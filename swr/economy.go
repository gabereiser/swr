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
	"crypto/sha256"
	"strings"
)

var MINER_DIFFICULTY int = 4

func mineBlockForEntity(entity Entity) {
	if entity.IsPlayer() {
		// because crypto is hilarious... #hodl
		result := sprintf("%x", sha256.Sum256([]byte(sprintf("%x", random_float()))))
		if strings.EqualFold(result[:MINER_DIFFICULTY], strings.Repeat("0", MINER_DIFFICULTY)) {
			entity.GetCharData().Bank += 10000
			DB().SavePlayerData(entity.(*PlayerProfile))
			entity.Send("\r\n}YYou have won the lottery. You have been awarded &W%d&Y credits!&d\r\n", 10000)
			entity.Send("%s", result)
			MINER_DIFFICULTY += 2
		}
	}
}

func updateMinerDifficulty() {
	if roll_dice("1d100") == 100 {
		MINER_DIFFICULTY -= 1
		if MINER_DIFFICULTY < 4 {
			MINER_DIFFICULTY = 4
		}
	}
}
