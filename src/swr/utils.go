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
	"math/rand"
	"strconv"
	"strings"
)

func roll_dice(d20 string) int {
	p := strings.Split(strings.ToLower(d20), "d")
	num_dice, _ := strconv.Atoi(p[0])
	sides, _ := strconv.Atoi(p[1])
	roll := 0
	for i := 0; i < num_dice; i++ {
		roll += rand.Intn(sides) + 1
	}
	return roll
}

func rand_min_max(min int, max int) int {
	return min + rand.Intn((max-min)+1)
}

func umin(min uint, value uint) uint {
	if value < min {
		return min
	}
	return value
}

func gen_player_char_id() uint {
	return uint(rand.Intn(9000000000)) + 9000000000
}

func gen_npc_char_id() uint {
	return uint(rand.Intn(1000000000)) + 1000000000
}

func gen_item_id() uint {
	return uint(rand.Intn(2000000000)) + 2000000000
}

func tune_random_frequency() string {
	buf := ""
	buf += strconv.Itoa(rand.Intn(2) + 1) // 1,2,3
	buf += strconv.Itoa(rand.Intn(9))     // 0-9
	buf += strconv.Itoa(rand.Intn(9))     // 0-9
	buf += "."
	buf += strconv.Itoa(rand.Intn(9))
	switch rand.Intn(3) {
	case 0:
		buf += "00"
	case 1:
		buf += "25"
	case 2:
		buf += "50"
	case 3:
		buf += "75"
	}
	return buf
}

func direction_reverse(direction string) string {
	switch direction {
	case "north":
		return "south"
	case "south":
		return "north"
	case "east":
		return "west"
	case "west":
		return "east"
	case "northwest":
		return "southeast"
	case "northeast":
		return "southwest"
	case "southwest":
		return "northeast"
	case "southeast":
		return "northwest"
	case "up":
		return "down"
	case "down":
		return "up"
	default:
		return "somewhere"
	}
}

func get_gender_for_code(gender string) string {
	g := strings.ToLower(gender)
	if g[0:1] == "m" {
		return "Male"
	}
	if g[0:1] == "f" {
		return "Female"
	}
	if g[0:1] == "n" {
		return "Neuter"
	}
	return "Male"
}

func distance_between_points(origin []float32, dest []float32) float32 {
	return 0
}

var ZERO_DISTANCE float32 = distance_between_points([]float32{0.0, 0.0}, []float32{0.0, 0.0})
var MAX_DISTANCE float32 = distance_between_points([]float32{-10000000.0, -10000000.0}, []float32{10000000.0, 10000000.0})
