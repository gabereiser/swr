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
	Mob    uint   `yaml:"mob"`
	Room   uint   `yaml:"room"`
	entity Entity `yaml:"-"`
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
	ship      uint                     `yaml:"shipId,omitempty"`
	Name      string                   `yaml:"name"`
	Desc      string                   `yaml:"desc"`
	Exits     map[string]uint          `yaml:"exits"`
	ExitFlags map[string]*RoomExitFlag `yaml:"exflags,omitempty"`
	Flags     []string                 `yaml:"flags,flow,omitempty"`
	RoomProgs map[string]string        `yaml:"roomProgs,omitempty"`
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
func (r *RoomData) HasExit(direction string) bool {
	if _, ok := r.Exits[direction]; ok {
		return true
	}
	return false
}

func (r *RoomData) AddItem(item Item) {
	r.Items = append(r.Items, item)
}

func (r *RoomData) ShipId() uint {
	return r.ship
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

func (r *RoomData) HasFlag(flag string) bool {
	for _, f := range r.Flags {
		if strings.EqualFold(f, flag) {
			return true
		}
	}
	return false
}
func (r *RoomData) RemoveFlag(flag string) {
	index := -1
	for i, f := range r.Flags {
		if strings.EqualFold(f, flag) {
			index = i
		}
	}
	if index > -1 {
		ret := make([]string, 0)
		ret = append(ret, r.Flags[:index]...)
		ret = append(ret, r.Flags[index+1:]...)
		r.Flags = ret
	}
}
func (r *RoomData) SetFlag(flag string) {
	found := false
	for _, f := range r.Flags {
		if strings.EqualFold(f, flag) {
			found = true
		}
	}
	if !found {
		r.Flags = append(r.Flags, flag)
	}
}
func (r *RoomData) GetExitFlags(direction string) *RoomExitFlag {
	if _, ok := r.ExitFlags[direction]; ok {
		return r.ExitFlags[direction]
	}
	return nil
}

func (r *RoomData) GetExitRoom(direction string) *RoomData {
	if rid, ok := r.Exits[direction]; ok {
		return DB().GetRoom(rid, r.ShipId())
	}
	return nil
}

func (r *RoomData) OpenDoor(entity Entity, direction string, silent bool) {
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
			if !r.UnlockDoor(entity, direction, silent) {
				if !silent {
					entity.Send("\r\n&YIt's locked.&d\r\n")
				}
			}
			return
		}
		flags.Closed = false
		flags.Locked = false // just in-case this is called with a nil entity, the system wants to open the door.
		if entity != nil && !silent {
			entity.Send("\r\nThe door slides open.\r\n")
		}
		// schedule a closing of the door 15 seconds from now
		ScheduleFunc(func() {
			flags.Closed = true
			r.SendToRoom(sprintf("\r\nThe door to the %s slides closed.\r\n", direction))
		}, false, 15)

		r.SendToOthers(entity, sprintf("\r\nThe door to the %s slides open.\r\n", direction))
		tflags := to_room.GetExitFlags(direction_reverse(direction))
		if tflags != nil {
			if tflags.Closed {
				was_locked := (tflags.Key != 0) // keys set on doors are always lockable.
				if tflags.Locked {
					tflags.Locked = false
				}
				tflags.Closed = false
				to_room.SendToOthers(entity, sprintf("\r\nThe door to the %s slides open.\r\n", direction_reverse(direction)))
				ScheduleFunc(func() {
					tflags.Closed = true
					tflags.Locked = was_locked
					to_room.SendToRoom(sprintf("\r\nThe door to the %s slides closed.\r\n", direction_reverse(direction)))
				}, false, 15)
			}
		}

	} else {
		if entity != nil && !silent {
			entity.Send("\r\nIt's already open.\r\n")
		}
	}
}

func (r *RoomData) CloseDoor(entity Entity, direction string, silent bool) {
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
		if entity != nil && !silent {
			entity.Send("\r\nThe door slides closed.\r\n")
		}
		r.SendToRoom(sprintf("\r\nThe door to the %s slides closed.\r\n", direction))
		tflags := to_room.GetExitFlags(direction_reverse(direction))
		if tflags != nil {
			if !tflags.Closed {
				tflags.Closed = true
				tflags.Locked = (tflags.Key != 0) // keys set on doors are always lockable.
				to_room.SendToRoom(sprintf("\r\nThe door to the %s slides closed.\r\n", direction_reverse(direction)))
			}
		}
	} else {
		if entity != nil && !silent {
			entity.Send("\r\nIt's already closed.\r\n")
		}
	}
}

func (r *RoomData) UnlockDoor(entity Entity, direction string, silent bool) bool {
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
			if !silent {
				entity.Send("\r\n&YYou hear a clunk as you unlock the door.&d\r\n")
			}
			ScheduleFunc(func() {
				flags.Locked = true
			}, false, 15)
		}
		flags.Locked = false
		tflags := to_room.GetExitFlags(direction_reverse(direction))
		was_locked := (tflags.Key != 0) // keys set on doors are always lockable...
		if tflags.Locked {
			tflags.Locked = false
			ScheduleFunc(func() {
				tflags.Locked = was_locked
			}, false, 15)
		}
		return true
	}
	if entity != nil && !silent {
		entity.Send("\r\n&RIt's unlocked.&d\r\n")
	}
	return false
}

func (r *RoomData) LockDoor(entity Entity, direction string, key Item) {
	flags := r.GetExitFlags(direction)
	if flags == nil {
		return
	}
	if !flags.Closed {
		r.CloseDoor(entity, direction, false)
	}
	flags.Locked = true
	flags.Key = key.GetTypeId()
	if entity != nil {
		entity.Send("\r\n&YWith a clunk you lock the door with %s %s.&d\r\n", key.GetData().Name, key.GetData().Name)
	}
}

func (r *RoomData) SendToRoom(message string) {
	for _, e := range DB().GetEntitiesInRoom(r.Id, r.ShipId()) {
		if e == nil {
			continue
		}
		if entity_unspeakable_state(e) {
			continue
		}
		e.Send(message)
	}
}
func (r *RoomData) SendToOthers(entity Entity, message string) {
	for _, e := range DB().GetEntitiesInRoom(r.Id, r.ShipId()) {
		if e == nil {
			continue
		}
		if e == entity {
			continue
		}
		if entity_unspeakable_state(e) {
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
		room := db.GetRoom(room_id, 0)
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

		room.SendToRoom(sprintf("\r\n&d%s&d\r\n", area.ResetMsg))
	}
	for _, spawn := range area.Items {
		room := db.GetRoom(spawn.Room, 0)
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
	for i := range area.Mobs {
		spawn := area.Mobs[i]
		mob := db.GetMob(spawn.Mob) // grabs the mob template
		if spawn.entity != nil {    // checks to see if we have a managed entity
			if spawn.entity.GetCharData().State == ENTITY_STATE_DEAD { // is it dead?
				//log.Printf("Removing dead entity %s\n", spawn.entity.GetCharData().Name)
				//db.RemoveEntity(spawn.entity)
				spawn.entity = nil // nil it out so we create a new one...
			} else {
				continue
			}
		}
		if spawn.entity == nil {
			//log.Printf("Nil entity for spawn, spawning %s\n", mob.GetCharData().Name)
			spawn.entity = db.SpawnEntity(mob)
			spawn.entity.GetCharData().Room = spawn.Room
			for _, e := range db.GetEntitiesInRoom(spawn.entity.GetCharData().Room, spawn.entity.GetCharData().Ship) {
				if e == nil {
					continue
				}
				if e.GetCharData().Id != spawn.entity.GetCharData().Id {
					spawn.entity.GetCharData().AI.OnGreet(e)
				}
			}
		}
		area.Mobs[i] = spawn
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
	room := DB().GetRoom(entity.RoomId(), entity.ShipId())
	if pg, ok := room.RoomProgs[evt]; ok {
		vm := mud_prog_init(entity)
		mud_prog_bind(vm, any...)
		_, err := vm.Run(pg)
		ErrorCheck(err)
	}
}
