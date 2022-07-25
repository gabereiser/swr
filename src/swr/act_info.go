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
		if entity.IsFighting() {
			entity.Send("\r\n&RYou can't quit while fighting!&d\r\n")
			return
		}
		entity.GetCharData().State = ENTITY_STATE_SLEEPING
		player := entity.(*PlayerProfile)
		player.Client.Close()
		entity.Send("\r\n&CThe world slowly fades away as you close your eyes and leave the game...&d\r\n\r\n")
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
		if e == nil {
			continue
		}
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

func do_score(entity Entity, args ...string) {
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		char := player.Char
		player.Send("\r\n&c╒═══( &W%-16s&c )═══════════════════╕&d\r\n", char.Name)
		player.Send("&c│ Title: &G%-25s&c         │&d▒\r\n", char.Title)
		player.Send("&c│  Race: &G%-25s&c         │&d▒\r\n", char.Race)
		player.Send("&c│ Level: &G%-25d&c         │&d▒\r\n", char.Level)
		player.Send("&c├─( Stats )────────────────────────────────┤&d▒\r\n")
		player.Send("&c│ STR: &G%-2d&c               XP: &G%-14d&c │&d▒\r\n", char.Stats[0], char.XP)
		player.Send("&c│ INT: &G%-2d&c         NEXT LVL: &G%-14d&c │&d▒\r\n", char.Stats[1], get_xp_for_level(char.Level))
		player.Send("&c│ DEX: &G%-2d&c            MONEY: &G%-14d&c │&d▒\r\n", char.Stats[2], char.Gold)
		player.Send("&c│ WIS: &G%-2d&c             BANK: &G%-14d&c │&d▒\r\n", char.Stats[3], char.Bank)
		player.Send("&c│ CON: &G%-2d&c                                  │&d▒\r\n", char.Stats[4])
		player.Send("&c│ CHA: &G%-2d&c                                  │&d▒\r\n", char.Stats[5])
		player.Send("&c╞══════════════════════════════════════════╡&d▒\r\n")
		player.Send("&c│ Weight: &G%3d kg&p(%3d kg)&c                   │&d▒\r\n", char.CurrentWeight(), char.MaxWeight())
		player.Send("&c│ Inventory: &G%3d&p(%3d)&c                      │&d▒\r\n", char.CurrentInventoryCount(), char.MaxInventoryCount())
		player.Send("&c├─( Equipment )────────────────────────────┤▒&d\r\n")
		player.Send("&c│       Head: &d%-20s&c         │&d▒\r\n", "None")
		player.Send("&c│      Torso: &d%-20s&c         │&d▒\r\n", "None")
		player.Send("&c│      Waist: &d%-20s&c         │&d▒\r\n", "None")
		player.Send("&c│       Legs: &d%-20s&c         │&d▒\r\n", "None")
		player.Send("&c│       Feet: &d%-20s&c         │&d▒\r\n", "None")
		player.Send("&c│      Hands: &d%-20s&c         │&d▒\r\n", "None")
		player.Send("&c│                                          │&d▒\r\n")
		player.Send("&c│     &RWeapon: &d%-20s&c         │&d▒\r\n", "None")
		player.Send("&c│                                          │&d▒\r\n")
		player.Send("&c├──( Skills )──────────────────────────────┤&d▒\r\n")
		for s, v := range char.Skills {
			player.Send("&c│ &w%-25s&c          &w%3d&c   │&d▒\r\n", s, v)
		}
		player.Send("&c├──( Languages )───────────────────────────┤&d▒\r\n")
		for s, v := range char.Languages {
			player.Send("&c│ &w%-25s&c          &w%3d&c   │&d▒\r\n", s, v)
		}
		player.Send("&c│   &cSpeaking: &w%-20s&c         │&d▒\r\n", char.Speaking)
		player.Send("&c└──────────────────────────────────────────┘&d▒\r\n")
		player.Send(" ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒\r\n")
	}
}

func do_inventory(entity Entity, args ...string) {
	player := entity.(*PlayerProfile)
	ch := entity.GetCharData()
	player.Send("\r\n&c╒═══( Inventory )═══════════════════╕\r\n")
	player.Send("&c├───────────────────────────────────┤&d▒\r\n")
	for _, item := range ch.Inventory {
		player.Send("&c│ %-34s│&d▒\r\n", item.GetData().Name)
	}
	player.Send("&c└───────────────────────────────────┘&d▒\r\n")
	player.Send(" ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒&d\r\n")
}
