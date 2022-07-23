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

import "fmt"

type AreaData struct {
	Name     string     `yaml:"name"`
	Author   string     `yaml:"author,omitempty"`
	Levels   []uint16   `yaml:"levels,flow"`
	Reset    uint       `yaml:"reset"`
	ResetMsg string     `yaml:"reset_msg`
	Rooms    []RoomData `yaml:"rooms"`
}

type RoomData struct {
	Id        uint                   `yaml:"id"`
	Name      string                 `yaml:"name"`
	Desc      string                 `yaml:"desc,flow"`
	Exits     map[string]uint        `yaml:"exits,flow"`
	ExitFlags map[string]interface{} `yaml:"exflags,flow,omitempty"`
	Flags     []string               `yaml:"flags,flow,omitempty"`
	RoomProgs map[string]string      `yaml:"roomProgs,flow,omitempty"`
	Area      *AreaData              `yaml:"-"`
	Items     []Item                 `yaml:"-"`
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
		if i.GetId() == item.GetId() {
			idx = id
		}
	}
	ret := make([]Item, len(r.Items)-1)
	ret = append(ret, r.Items[:idx]...)
	ret = append(ret, r.Items[idx+1:]...)
	r.Items = ret
}
func room_get_blocked_exit_flags(exitFlags map[string]interface{}) (locked bool, closed bool) {
	locked = false
	closed = false
	if c, ok := exitFlags["closed"]; ok {
		closed = c.(bool)
	}
	if l, ok := exitFlags["locked"]; ok {
		locked = l.(bool)
	}
	return locked, closed
}
func room_get_exit_status(exitFlags map[string]interface{}) string {
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
