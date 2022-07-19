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
		roll += rand.Intn(sides-1) + 1
	}
	return roll
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

func get_skill_value(ch *CharData, skill string) int {
	if v, ok := ch.Skills[strings.ToLower(skill)]; ok {
		return v
	}
	return 0
}
