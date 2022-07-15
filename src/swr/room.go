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

type AreaData struct {
	Name     string               `yaml:"name"`
	Author   string               `yaml:"author,omitempty"`
	Levels   []uint16             `yaml:"levels,flow"`
	Reset    uint                 `yaml:"reset"`
	ResetMsg string               `yaml:"reset_msg`
	Rooms    map[uint]RoomData    `yaml:"rooms"`
	Mobs     map[uint]interface{} `yaml:"mobs,omitempty"`
	Items    map[uint]interface{} `yaml:"items,omitempty"`
}

type RoomData struct {
	Id        uint                   `yaml:"-"`
	Name      string                 `yaml:"name"`
	Desc      string                 `yaml:"desc,flow"`
	Exits     map[string]uint        `yaml:"exits,flow"`
	ExitFlags map[string]interface{} `yaml:"exflags,flow,omitempty"`
	Flags     []string               `yaml:"flags,flow,omitempty"`
	RoomProgs []string               `yaml:"room_progs,flow,omitempty"`
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
