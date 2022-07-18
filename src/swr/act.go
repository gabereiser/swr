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
	"strings"
)

func do_nothing(entity Entity, args ...string) {
	entity.Send("\r\n&rInput not recognized.&d\r\n")
}

func do_quit(entity Entity, args ...string) {
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		player.Client.Close()
	}
}
func do_save(entity Entity, args ...string) {
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		DB().SavePlayerData(player)
		entity.Send("\r\n&YSaved. Ok.&d\r\n")
	}
}
func do_look(entity Entity, args ...string) {
	if entity.IsPlayer() {
		if len(args) == 0 { // l or look with no args
			roomId := entity.RoomId()
			log.Printf("Entity RoomId %d", roomId)
			room := DB().GetRoom(roomId)
			if room != nil {
				entity.Send(fmt.Sprintf("\r\n-=-=-=-=-=-=-=-=-=( %s %d )=-=-=-=-=-=-=-=-=-\r\n", room.Name, room.Id))
				entity.Send(room.Desc)
				entity.Send("\r\nExits:\r\n")
				for dir, toRoom := range room.Exits {
					if k, ok := room.ExitFlags[dir]; ok {
						exit_flags := k.(map[string]interface{})
						ext := room_get_exit_status(exit_flags)
						entity.Send(fmt.Sprintf(" - %s%s[%d]\r\n", dir, ext, toRoom))
					} else {
						entity.Send(fmt.Sprintf(" - %s [%d]\r\n", dir, toRoom))
					}
				}
			} else {
				log.Fatalf("Entity %s is in room %d, only it doesn't exist and crashed the server.", entity.Name(), entity.RoomId())
			}

		} else {
			for _, e := range DB().GetEntitiesInRoom(entity.RoomId()) {
				if e != entity {
					ch := e.GetCharData()
					for _, keyword := range ch.Keywords {
						if strings.HasPrefix(strings.ToLower(keyword), strings.ToLower(args[0])) {
							entity.Send("You look at %s and see...\r\n%s\r\n", ch.Title, ch.Desc)
							return
						}
					}
				}
			}
		}
	}

}

func do_say(entity Entity, args ...string) {
	words := strings.Join(args, " ")
	var speaker *CharData
	if entity.IsPlayer() {
		speaker = &(entity.(*PlayerProfile).Char)
	} else {
		speaker = entity.(*CharData)
	}
	if entity.IsPlayer() {
		entity.Send(fmt.Sprintf("You say \"%s\"\n", words))
	}
	entities := DB().GetEntitiesInRoom(speaker.RoomId())
	for _, ex := range entities {
		if ex != entity {
			if ex.IsPlayer() {
				listener := &(ex.(*PlayerProfile).Char)
				ex.Send(fmt.Sprintf("%s says \"%s\"\n", speaker.Name(), language_spoken(speaker, listener, words)))
			} else {
				ex.Send(fmt.Sprintf("%s says \"%s\"\n", speaker.Name(), words))
			}
		}
	}
}

func do_who(entity Entity, args ...string) {
	db := DB()
	total := 0
	entity.Send("\r\n&G-=-=-=-=-=-=-=-=-=-=-=-=-=-=( &WWho&G )=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-&d\r\n")
	for _, e := range db.entities {
		if e.IsPlayer() {
			player := e.(*PlayerProfile)
			entity.Send(fmt.Sprintf("&W%-54s&G [ &WLevel %2d&G ]\r\n", player.Char.Title, player.Char.Level))
			total++
		}
	}
	entity.Send(fmt.Sprintf("\r\n&G-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=( &W%3d&Y Online&G )=-=-=&d\r\n", total))
}

func do_score(entity Entity, args ...string) {
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		char := player.Char
		player.Send("\r\n&c╒════════════════( &W%16s&c )══════╕&d\r\n", char.CharName)
		player.Send("&c│ Title: &G%-25s&c         │&d▒\r\n", char.Title)
		player.Send("&c│  Race: &G%-25s&c         │&d▒\r\n", char.Race)
		player.Send("&c│ Level: &G%-25d&c         │&d▒\r\n", char.Level)
		player.Send("&c├─( Stats )────────────────────────────────┤&d▒\r\n")
		player.Send("&c│ STR: &G%-2d&c               XP: &G%-14d&c │&d▒\r\n", char.Stats[0], char.XP)
		player.Send("&c│ INT: &G%-2d&c            MONEY: &G%-14d&c │&d▒\r\n", char.Stats[1], char.Gold)
		player.Send("&c│ DEX: &G%-2d&c             BANK: &G%-14d&c │&d▒\r\n", char.Stats[2], char.Bank)
		player.Send("&c│ WIS: &G%-2d&c                                  │&d▒\r\n", char.Stats[3])
		player.Send("&c│ CON: &G%-2d&c                                  │&d▒\r\n", char.Stats[4])
		player.Send("&c│ CHA: &G%-2d&c                                  │&d▒\r\n", char.Stats[5])
		player.Send("&c╞══════════════════════════════════════════╡&d▒\r\n")
		player.Send("&c│ Weight: &G%3d kg&c                           │&d▒\r\n", char.CurrentWeight())
		player.Send("&c│ Inventory: &G%3d&p(%3d)&c                      │&d▒\r\n", char.CurrentInventoryCount(), (char.Level*3)+char.Stats[0])
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
		player.Send("&c└──────────────────────────────────────────┘▒&d\r\n")
		player.Send(" ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒\r\n")
	}
}
