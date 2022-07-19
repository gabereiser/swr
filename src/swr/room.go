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
}

func RoomFromMap(data map[string]interface{}) *RoomData {
	room := new(RoomData)
	room.Id = uint(data["id"].(int))
	room.Name = data["name"].(string)
	room.Desc = data["desc"].(string)
	room.Exits = make(map[string]uint)
	room.ExitFlags = make(map[string]interface{})
	room.Flags = make([]string, 0)
	room.RoomProgs = make(map[string]string)
	if d, ok := data["exits"]; ok {
		for dir, e := range d.(map[string]interface{}) {
			room.Exits[dir] = uint(e.(int))
		}
	}
	if d, ok := data["exflags"]; ok {
		for dir, f := range d.(map[string]interface{}) {
			room.ExitFlags[dir] = f
		}
	}
	if d, ok := data["flags"]; ok {
		for f := range d.([]interface{}) {
			flag := data["flags"].([]interface{})[f]
			room.Flags = append(room.Flags, flag.(string))
		}
	}
	if d, ok := data["roomProgs"]; ok {
		for evt, prog := range d.(map[string]string) {
			room.RoomProgs[evt] = prog
		}
	}
	return room
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

func room_get_exit_status(exitFlags map[string]interface{}) string {
	ret := " "
	locked := false
	if l, ok := exitFlags["locked"]; ok {
		locked = l.(bool)
	}
	if locked {
		ret += "(locked) "
	}
	return ret
}
