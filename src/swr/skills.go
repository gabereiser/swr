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

var skill_list []string = []string{
	"aerobics",
	"astrophysics",
	"astronomy",
	"bartering",
	"blasters",
	"bowcasters",
	"business",
	"chemical-analysis",
	"chemical-synthesis",
	"claymores",
	"cloning",
	"defusing",
	"engineering",
	"electronics",
	"first-aid",
	"force",
	"force-pikes",
	"grenades",
	"gunnery",
	"healing",
	"hunting",
	"hyperdrives",
	"lightsabers",
	"lore",
	"martial-arts",
	"mines",
	"missiles",
	"piloting",
	"production",
	"rifles",
	"repeaters",
	"targeting",
	"telekinesis",
	"tracking",
	"vibro-blades",
	"xenosciences",
}

func is_skill(s string) bool {
	for _, skill := range skill_list {
		if skill == s {
			return true
		}
	}
	return false
}
