/*  Space Wars Rebellion Mud
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
	return r.ExitFlags[direction]
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
	ret := " "
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
			room.ExitFlags[dir].Closed = f.Closed
			room.ExitFlags[dir].Locked = f.Locked
			room.ExitFlags[dir].Key = f.Key
		}
		rem_items := make([]Item, 0)
		for _, i := range room.Items {
			if i != nil {
				if i.IsCorpse() {
					rem_items = append(rem_items, i)
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
				if i.GetId() == item.GetId() {
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
			if e.GetCharData().Name == mob.GetCharData().Name {
				exists = true
			}
		}
		if !exists {
			m := db.SpawnEntity(mob)
			m.GetCharData().Room = spawn.Room
			if mob.GetCharData().AI == nil {
				m.GetCharData().AI = MakeGenericBrain(m)
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
