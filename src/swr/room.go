/*  Star Wars Role-Playing Mud
 *  Copyright (C) 2022 @{See Authors}
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */
package swr

import (
	"fmt"
	"log"
	"strings"
)

type MobSpawn struct {
	Mob  uint `yaml:"mob"`
	Room uint `yaml:"room"`
}

type ItemSpawn struct {
	Item uint `yaml:"item"`
	Room uint `yaml:"room"`
}

type AreaData struct {
	Name     string      `yaml:"name"`
	Author   string      `yaml:"author,omitempty"`
	Levels   []uint16    `yaml:"levels,flow"`
	Reset    uint        `yaml:"reset"`
	ResetMsg string      `yaml:"reset_msg"`
	Rooms    []RoomData  `yaml:"rooms"`
	Mobs     []MobSpawn  `yaml:"mobs,omitempty"`
	Items    []ItemSpawn `yaml:"items,omitempty"`
}

type RoomData struct {
	Id        uint                     `yaml:"id"`
	Name      string                   `yaml:"name"`
	Desc      string                   `yaml:"desc,flow"`
	Exits     map[string]uint          `yaml:"exits,flow"`
	ExitFlags map[string]*RoomExitFlag `yaml:"exflags,flow,omitempty"`
	Flags     []string                 `yaml:"flags,flow,omitempty"`
	RoomProgs map[string]string        `yaml:"roomProgs,flow,omitempty"`
	Area      *AreaData                `yaml:"-"`
	Items     []Item                   `yaml:"-"`
}

type RoomExitFlag struct {
	Locked bool `yaml:"locked,omitempty"`
	Closed bool `yaml:"closed,omitempty"`
	Key    uint `yaml:"key,omitempty"`
}

func (r *RoomData) String() string {
	return fmt.Sprintf("ROOM:[%d-%s]", r.Id, r.Name)
}
func (r *RoomData) GetEntities() []Entity {
	return DB().GetEntitiesInRoom(r.Id)
}
func (r *RoomData) HasExit(direction string) bool {
	if _, ok := r.Exits[direction]; ok {
		return true
	}
	return false
}

func (r *RoomData) AddItem(item Item) {
	r.Items = append(r.Items, item)
}
func (r *RoomData) RemoveItem(item Item) {
	idx := -1
	for id := range r.Items {
		i := r.Items[id]
		if i == nil {
			continue
		}
		if i.GetId() == item.GetId() {
			idx = id
		}
	}
	ret := make([]Item, len(r.Items)-1)
	ret = append(ret, r.Items[:idx]...)
	ret = append(ret, r.Items[idx+1:]...)
	r.Items = ret
}

func (r *RoomData) FindItem(keyword string) Item {
	for id := range r.Items {
		i := r.Items[id]
		if i == nil {
			continue
		}
		keys := i.GetKeywords()
		for k := range keys {
			key := keys[k]
			if strings.HasPrefix(key, keyword) {
				return i
			}
		}
	}
	return nil
}

func (r *RoomData) GetExitFlags(direction string) *RoomExitFlag {
	if _, ok := r.ExitFlags[direction]; ok {
		return r.ExitFlags[direction]
	}
	return nil
}

func (r *RoomData) GetExitRoom(direction string) *RoomData {
	if _, ok := r.Exits[direction]; ok {
		return DB().GetRoom(r.Exits[direction])
	}
	return nil
}

func (r *RoomData) OpenDoor(entity Entity, direction string) {
	flags := r.GetExitFlags(direction)
	to_room := r.GetExitRoom(direction)
	if to_room == nil {
		return
	}
	if flags == nil {
		return
	}
	if flags.Closed {
		if flags.Locked && entity != nil {
			if !r.UnlockDoor(entity, direction) {
				entity.Send("\r\n&YIt's locked.&d\r\n")
			}
			return // make them call open again after it's unlocked.
		}
		flags.Closed = false
		flags.Locked = false // just in-case this is called with a nil entity, the system wants to open the door.
		if entity != nil {
			entity.Send("\r\n&GYou open the door.&d\r\n")
		}
		// schedule a closing of the door 15 seconds from now
		ScheduleFunc(func() {
			// nil Entity because no one closes it, the system does.
			r.CloseDoor(nil, direction)
		}, false, 15)

		// tell others the door is open
		for _, e := range r.GetEntities() {
			if e != nil {
				if e != entity {
					e.Send("\r\nThe door to the %s opens.\r\n", direction)
				}
			}
		}
		to_room.OpenDoor(nil, direction_reverse(direction))
	}
}

func (r *RoomData) CloseDoor(entity Entity, direction string) {
	flags := r.GetExitFlags(direction)
	to_room := r.GetExitRoom(direction)
	if to_room == nil {
		return
	}
	if flags == nil {
		return
	}
	if !flags.Closed {
		flags.Closed = true
		if entity != nil {
			entity.Send("\r\n&GYou close the door.&d\r\n")
		}
		for _, e := range r.GetEntities() {
			if e != nil {
				if e != entity {
					e.Send("\r\nThe door to the %s closes.\r\n", direction)
				}
			}
		}
		to_room.CloseDoor(nil, direction_reverse(direction))
	}
}

func (r *RoomData) UnlockDoor(entity Entity, direction string) bool {
	flags := r.GetExitFlags(direction)
	to_room := r.GetExitRoom(direction)
	if to_room == nil {
		return false
	}
	if flags == nil {
		return false
	}
	if flags.Locked {
		if entity != nil {
			key := entity.GetCharData().GetItem(flags.Key)
			if key == nil {
				entity.Send("\r\n&RYou don't have the key.&d\r\n")
				return false
			}
			entity.Send("\r\n&YYou unlock the door.&d\r\n")
			ScheduleFunc(func() {
				// nil Entity because no one closes it, the system does.
				r.CloseDoor(nil, direction)
				r.LockDoor(nil, direction, key)
			}, false, 15)
		}
		flags.Locked = false
		to_room.UnlockDoor(nil, direction_reverse(direction))
		return true
	}
	return false
}

func (r *RoomData) LockDoor(entity Entity, direction string, key Item) {
	flags := r.GetExitFlags(direction)
	if flags == nil {
		return
	}
	flags.Closed = true
	flags.Locked = true
	flags.Key = key.GetTypeId()
	if entity != nil {
		entity.Send("\r\n&YYou lock the door with %s %s.&d\r\n", get_preface_for_name(key.GetData().Name), key.GetData().Name)
	}
}

func (r *RoomData) SendToRoom(message string) {
	for _, e := range r.GetEntities() {
		if e == nil {
			continue
		}
		e.Send(message)
	}
}
func room_get_blocked_exit_flags(exitFlags *RoomExitFlag) (locked bool, closed bool) {
	locked = false
	closed = false
	locked = exitFlags.Locked
	closed = exitFlags.Closed
	return locked, closed
}
func room_get_exit_status(exitFlags *RoomExitFlag) string {
	ret := ""
	locked, closed := room_get_blocked_exit_flags(exitFlags)
	if closed {
		ret += "(closed) "
	}
	if locked {
		ret += "(locked) "
	}
	return ret
}

func area_reset(area *AreaData) {
	db := DB()
	if area == nil {
		return
	}
	for _, r := range area.Rooms {
		room_id := r.Id
		room := db.GetRoom(room_id)
		if room == nil {
			log.Printf("Error: roomId %d doesn't exist! area_reset(%s)", room_id, area.Name)
			continue
		}
		for dir, f := range r.ExitFlags {
			room.ExitFlags[dir] = f
		}
		rem_items := make([]Item, 0)
		for _, i := range room.Items {
			if i != nil {
				if i.IsCorpse() {
					rem_items = append(rem_items, i)
				}
				if i.IsContainer() {
					if i.GetData().Type == ITEM_TYPE_TRASH_BIN {
						i.GetData().Items = make([]Item, 0)
					}
				}
			}
		}
		for _, i := range rem_items {
			room.RemoveItem(i)
		}
		for _, e := range room.GetEntities() {
			if e != nil {
				if e.IsPlayer() {
					e.Send("\r\n&d%s&d\r\n", area.ResetMsg)
				}
			}
		}
	}
	for _, spawn := range area.Items {
		room := db.GetRoom(spawn.Room)
		item := db.GetItem(spawn.Item)
		exists := false
		for _, i := range room.Items {
			if i != nil {
				if i.GetData().Name == item.GetData().Name {
					exists = true
				}
			}
		}
		if !exists {
			room.AddItem(item_clone(item))
		}
	}
	for _, spawn := range area.Mobs {
		mob := db.GetMob(spawn.Mob)
		exists := false
		for _, e := range db.GetEntitiesInRoom(spawn.Room) {
			if e == nil {
				continue
			}
			if e.IsPlayer() {
				continue
			}
			if e.GetCharData().Id == mob.GetCharData().Id || e.GetCharData().OId == mob.GetCharData().Id {
				exists = true
			}
		}
		if !exists {
			m := db.SpawnEntity(mob)
			m.GetCharData().Room = spawn.Room
			if mob.GetCharData().AI == nil {
				m.GetCharData().AI = MakeGenericBrain(m)
				m.GetCharData().AI.OnSpawn()
				for _, e := range db.GetEntitiesInRoom(spawn.Room) {
					if e == nil {
						continue
					}
					if e != m {
						m.GetCharData().AI.OnGreet(e)
					}
				}
			}
		}
	}
	ScheduleFunc(func() {
		area_reset(area)
	}, false, area.Reset)
}

func get_direction_string(direction string) string {
	direction = strings.TrimSpace(strings.ToLower(direction))
	if strings.HasPrefix(direction, "ne") {
		return "northeast"
	}
	if strings.HasPrefix(direction, "nw") {
		return "northwest"
	}
	if strings.HasPrefix(direction, "se") {
		return "southeast"
	}
	if strings.HasPrefix(direction, "sw") {
		return "southwest"
	}
	if strings.HasPrefix(direction, "s") {
		return "south"
	}
	if strings.HasPrefix(direction, "e") {
		return "east"
	}
	if strings.HasPrefix(direction, "w") {
		return "west"
	}
	if strings.HasPrefix(direction, "n") {
		return "north"
	}
	if strings.HasPrefix(direction, "u") {
		return "up"
	}
	if strings.HasPrefix(direction, "d") {
		return "down"
	}
	return direction
}

func room_prog_exec(entity Entity, evt string, any ...interface{}) {
	room := DB().GetRoom(entity.RoomId())
	if pg, ok := room.RoomProgs[evt]; ok {
		vm := mud_prog_init(entity)
		mud_prog_bind(vm, any...)
		_, err := vm.Run(pg)
		ErrorCheck(err)
	}
}
