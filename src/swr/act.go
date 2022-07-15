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

func do_nothing(entity Entity, args ...interface{}) {
	entity.Send("\r\n&WHuh?\r\n")
}

func do_look(entity Entity, args ...interface{}) {
	if len(args) <= 1 { // l or look with no args
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
		fmt.Printf("Args %v\n", args)
	}
}

func do_say(entity Entity, args ...interface{}) {
	parts := args[0].([]string)
	if strings.ToLower(parts[0]) == "say" {
		parts = parts[1:]
	}
	words := strings.Join(parts, " ")
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

func do_who(entity Entity, args ...interface{}) {
	db := DB()
	total := 0
	entity.Send("\r\n&G-=-=-=-=-=-=-=-=-=-=-=-=-=-=( &WWho&G )=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-&d\r\n")
	for _, e := range db.entities {
		if e.IsPlayer() {
			player := e.(*PlayerProfile)
			entity.Send(fmt.Sprintf("&W%-49s&G\t[ &WLevel %2d&G ]\r\n", player.Char.Title, player.Char.Level))
			total++
		}
	}
	entity.Send(fmt.Sprintf("\r\n&G-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=( &W%3d&Y Online&G )=-=-=&d\r\n", total))
}
