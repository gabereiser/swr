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
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func Init() {
	// Ensure that the player directories exists
	for _, p := range "abcdefghijklmnopqrstuvwxyz" {
		_ = os.MkdirAll(fmt.Sprintf("data/accounts/%s", string(p)), 0755)
	}
	_ = os.MkdirAll("backup", 0755)
	// Start the scheduler
	Scheduler()
}

func Main() {

	log.Printf("Starting version %s\n", version)

	DB().Load()
	CommandsLoad()
	LanguageLoad()
	StartBackup()
	area := new(AreaData)
	area.Name = "template"
	area.Author = "Admin"
	area.Levels = []uint16{1, 10}
	area.Reset = 360
	area.ResetMsg = "You hear the sound of sqeaking in the distance."
	area.Rooms = make(map[uint]RoomData)
	area.Items = make(map[uint]interface{})
	area.Mobs = make(map[uint]interface{})
	for i := 1; i <= 10; i++ {
		area.Rooms[uint(i)] = RoomData{
			Id:        uint(i),
			Name:      "The void",
			Desc:      "All you see is black as you are in a void. Nothing exists here.",
			Exits:     make(map[string]uint),
			ExitFlags: make(map[string]interface{}),
			RoomProgs: make([]string, 0),
			Flags:     make([]string, 0),
		}
	}
	for i := 1; i <= 10; i++ {
		area.Mobs[uint(i)] = CharData{
			Room:      uint(i),
			CharName:  "A generic mob",
			Keywords:  []string{"man", "male", "mob"},
			Title:     "",
			Desc:      "A generic mob stands here. A blank stare in it's eyes.",
			Race:      race_list[i-1],
			Gender:    "Male",
			Level:     1,
			XP:        0,
			Gold:      0,
			Hp:        []uint16{10, 10},
			Mp:        []uint16{0, 0},
			Mv:        []uint16{10, 10},
			Stats:     []uint16{10, 10, 10, 10, 10, 10},
			Skills:    map[string]int{"martial arts": 1, "kick": 1, "piloting": 1},
			Languages: map[string]int{"basic": 100, strings.ToLower(race_list[i-1]): 100},
			Speaking:  "basic",
			Equipment: map[string]Item{"head": ItemData{}, "torso": ItemData{}, "waist": ItemData{}, "legs": ItemData{}, "feet": ItemData{}, "hands": ItemData{}, "weapon": ItemData{}},
			Inventory: []Item{},
			State:     "normal",
			Brain:     "generic",
		}
	}
	DB().SaveArea(area)
	ServerStart(Config().Addr)

	DB().Save()
}

func GetVersion() string {
	return version
}

func roll_dice(d20 string) uint {
	p := strings.Split(strings.ToLower(d20), "d")
	num_dice, _ := strconv.Atoi(p[0])
	sides, _ := strconv.Atoi(p[1])
	roll := uint(0)
	for i := 0; i < num_dice; i++ {
		roll += (rand.Intn(sides-1) + 1)
	}
	return roll
}
