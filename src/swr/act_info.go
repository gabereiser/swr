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

func do_quit(entity Entity, args ...string) {
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		player.Client.Close()
	}
}

func do_qui(entity Entity, args ...string) {
	entity.Send("\r\n}RYou'll have to be more specific when quitting!&d\r\n&RType &Wquit&R to quit!&d\r\n")
}

func do_who(entity Entity, args ...string) {
	db := DB()
	total := 0
	entity.Send("\r\n")
	entity.Send(MakeTitle("Who", ANSI_TITLE_STYLE_NORMAL, ANSI_TITLE_ALIGNMENT_CENTER))
	for _, e := range db.entities {
		if e.IsPlayer() {
			player := e.(*PlayerProfile)
			entity.Send(fmt.Sprintf("&W%-54s&G [ &WLevel %2d&G ]\r\n", player.Char.Title, player.Char.Level))
			total++
		}
	}
	entity.Send("\r\n")
	entity.Send(MakeTitle(fmt.Sprintf("%d Online", total), ANSI_TITLE_STYLE_NORMAL, ANSI_TITLE_ALIGNMENT_RIGHT))
	entity.Send("\r\n")
}
