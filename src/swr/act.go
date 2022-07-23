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

func do_save(entity Entity, args ...string) {
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		DB().SavePlayerData(player)
		entity.Send("\r\n&YSaved. Ok.&d\r\n")
	}
}
func do_look(entity Entity, args ...string) {
	if entity_unspeakable_state(entity) {
		entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
		return
	}
	if entity.IsPlayer() {
		ch := entity.GetCharData()
		if ch.State == ENTITY_STATE_DEAD {
			return
		}
		if ch.State == ENTITY_STATE_SLEEPING || ch.State == ENTITY_STATE_UNCONSCIOUS {
			entity.Send("\r\n&cIn your dreams?...&d\r\n")
			return
		}
		if len(args) == 0 { // l or look with no args
			roomId := entity.RoomId()
			room := DB().GetRoom(roomId)
			if room != nil {
				entity.Send(fmt.Sprintf("\r\n%s\r\n",
					MakeTitle(room.Name,
						ANSI_TITLE_STYLE_NORMAL,
						ANSI_TITLE_ALIGNMENT_CENTER)))
				entity.Send(room.Desc)
				entity.Send("\r\nExits:\r\n")
				for dir, to_room_id := range room.Exits {
					to_room := DB().GetRoom(to_room_id)
					if k, ok := room.ExitFlags[dir]; ok {
						exit_flags := k.(map[string]interface{})
						ext := room_get_exit_status(exit_flags)
						entity.Send(fmt.Sprintf(" - %s %s %s\r\n", dir, to_room.Name, ext))
					} else {
						entity.Send(fmt.Sprintf(" - %s %s\r\n", dir, to_room.Name))
					}
				}
				entity.Send("\r\n")
				for i := range room.Items {
					item := room.Items[i]
					entity.Send("&Y%s&d\r\n", item.GetData().Name)
				}

				for _, e := range DB().GetEntitiesInRoom(room.Id) {
					if e != entity {
						entity.Send("&P%s&d\r\n", e.GetCharData().Name)
					}
				}
			} else {
				log.Fatalf("Entity %s is in room %d, only it doesn't exist and crashed the server.", entity.GetCharData().Name, entity.RoomId())
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

func do_north(entity Entity, args ...string) {
	do_direction(entity, "north")
}
func do_northwest(entity Entity, args ...string) {
	do_direction(entity, "northwest")
}
func do_northeast(entity Entity, args ...string) {
	do_direction(entity, "northeast")
}
func do_east(entity Entity, args ...string) {
	do_direction(entity, "east")
}
func do_west(entity Entity, args ...string) {
	do_direction(entity, "west")
}
func do_southeast(entity Entity, args ...string) {
	do_direction(entity, "southeast")
}
func do_southwest(entity Entity, args ...string) {
	do_direction(entity, "southwest")
}
func do_south(entity Entity, args ...string) {
	do_direction(entity, "south")
}
func do_up(entity Entity, args ...string) {
	do_direction(entity, "up")
}
func do_down(entity Entity, args ...string) {
	do_direction(entity, "down")
}

func do_direction(entity Entity, direction string) {
	if entity_unspeakable_state(entity) {
		entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
		return
	}
	if entity.GetCharData().State == ENTITY_STATE_SITTING {
		entity.Send("\r\nYou are unable to move while sitting.\r\n")
		return
	}
	db := DB()
	room := db.GetRoom(entity.RoomId())
	if !room.HasExit(direction) {
		entity.Send("\r\nYou can't go that way.\r\n")
		return
	} else {
		to_room := db.GetRoom(room.Exits[direction])
		if to_room == nil {
			entity.Send("\r\n&RThat room doesn't exist!\r\n")
			return
		} else {
			locked := false
			closed := false
			if flags, ok := room.ExitFlags[direction]; ok {
				locked, closed = room_get_blocked_exit_flags(flags.(map[string]interface{}))
			}
			if locked {
				entity.Send("\r\nIt's locked.\r\n")
				return
			}
			if closed {
				entity.Send("\r\nIt's closed.\r\n")
				return
			}
			if entity.CurrentMv() > 0 {
				entity.GetCharData().Mv[0]--
				entity.GetCharData().Room = to_room.Id
				for _, e := range room.GetEntities() {
					if entity_unspeakable_state(e) {
						continue
					}
					if e != entity {
						e.Send("\r\n%s has left going %s.\r\n", entity.GetCharData().Name, direction)
						if e.GetCharData().AI != nil {
							e.GetCharData().AI.OnEnter(entity)
						}
					}
				}
				for _, e := range to_room.GetEntities() {
					if entity_unspeakable_state(e) {
						continue
					}
					if e != entity {
						e.Send("\r\n%s has arrived from the %s.\r\n", entity.GetCharData().Name, direction_reverse(direction))
						if e.GetCharData().AI != nil {
							e.GetCharData().AI.OnEnter(entity)
						}
					}
				}
				do_look(entity)
				return
			} else {
				entity.Send("\r\n&You are too exhausted.\r\n")
				return
			}
		}

	}
}

func do_stand(entity Entity, args ...string) {
	ch := entity.GetCharData()
	if ch.State == ENTITY_STATE_DEAD {
		entity.Send("\r\n&RYou can't move when you're dead.&d\r\n")
		return
	}
	if ch.State == ENTITY_STATE_UNCONSCIOUS {
		entity.Send("\r\n&YYou are unconscious...&d\r\n")
		return
	}
	if ch.State == ENTITY_STATE_SITTING || ch.State == ENTITY_STATE_SLEEPING {
		ch.State = ENTITY_STATE_NORMAL
		for _, e := range DB().GetEntitiesInRoom(entity.RoomId()) {
			if entity_unspeakable_state(e) {
				continue
			}
			if e != entity {
				e.Send("\r\n&d%s stands to their feet.\r\n", ch.Name)
			}
		}
		entity.Send("\r\n&dYou spring to your feet.\r\n")
		return
	} else {
		entity.Send("\r\n&dYou are already standing.\r\n")
		return
	}
}

func do_sit(entity Entity, args ...string) {
	ch := entity.GetCharData()
	if ch.State == ENTITY_STATE_DEAD {
		entity.Send("\r\n&RYou can't move when you're dead.&d\r\n")
		return
	} else if ch.State == ENTITY_STATE_UNCONSCIOUS {
		entity.Send("\r\n&YYou are unconscious...&d\r\n")
		return
	} else if ch.State == ENTITY_STATE_NORMAL {
		ch.State = ENTITY_STATE_SITTING
		for _, e := range DB().GetEntitiesInRoom(entity.RoomId()) {
			if entity_unspeakable_state(e) {
				continue
			}
			if e != entity {
				e.Send("\r\n&d%s sits down.\r\n", ch.Name)
			}
		}
		entity.Send("\r\n&dYou sit down.\r\n")
	} else if ch.State == ENTITY_STATE_SITTING {
		entity.Send("\r\n&dYou are already sitting.\r\n")
	} else {
		entity.Send("\r\n&dYou can't do that right now.\r\n")
	}
}

func do_sleep(entity Entity, args ...string) {
	ch := entity.GetCharData()
	if ch.State == ENTITY_STATE_DEAD {
		entity.Send("\r\n&RYou're already permanently asleep (*DEAD*).&d\r\n")
		return
	}
	if ch.State == ENTITY_STATE_SLEEPING {
		entity.Send("\r\n&dYou're already asleep.\r\n")
		return
	}
	if ch.State == ENTITY_STATE_UNCONSCIOUS {
		entity.Send("\r\n&dYou're unconscious.\r\n")
		return
	}
	if ch.State == ENTITY_STATE_FIGHTING {
		entity.Send("\r\n&dYou can't sleep when you're fighting.\r\n")
		return
	}
	if ch.State == ENTITY_STATE_GUNNING || ch.State == ENTITY_STATE_PILOTING {
		entity.Send("\r\n&dYou can't sleep.\r\n")
		return
	}
	ch.State = ENTITY_STATE_SLEEPING
	entity.Send("\r\n&dYou lay down and fall asleep.\r\n")
	for _, e := range DB().GetEntitiesInRoom(entity.RoomId()) {
		if entity_unspeakable_state(e) {
			continue
		}
		if e != entity {
			e.Send("\r\n&d%s lays down and falls asleep.\r\n", ch.Name)
		}
	}
}

func do_open(entity Entity, args ...string) {
	if entity_unspeakable_state(entity) {
		entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
		return
	}
	if entity.GetCharData().State == ENTITY_STATE_SITTING {
		entity.Send("\r\nYou are unable to move while sitting.\r\n")
		return
	}
	db := DB()
	room := db.GetRoom(entity.RoomId())
	if len(args) == 0 {
		entity.Send("\r\n&ROpen what?&d\r\n")
		return
	}
	direction := strings.ToLower(args[0])
	if !room.HasExit(direction) {
		entity.Send("\r\n&ROpen what?&d.\r\n")
	}
	if flags, ok := room.ExitFlags[direction]; ok {
		locked, closed := room_get_blocked_exit_flags(flags.(map[string]interface{}))
		if locked {
			// TODO: Search users inventory for the key.
			entity.Send("\r\n&RYou don't have the key.&d\r\n")
			return
		}
		if closed {
			entity.Send("\r\n&GYou open the door.&d\r\n")
			room.ExitFlags[direction].(map[string]interface{})["closed"] = false
			to_room := db.GetRoom(room.Exits[direction])
			reverse_direction := direction_reverse(direction)
			if room.Id == to_room.Exits[reverse_direction] {
				if _, ok := to_room.ExitFlags[reverse_direction]; ok {
					to_room.ExitFlags[reverse_direction].(map[string]interface{})["closed"] = false
				}
			}
			for _, e := range db.GetEntitiesInRoom(entity.RoomId()) {
				if e != nil {
					if e != entity {
						e.Send("\r\n&P%s&d opens the door to the %s.\r\n", entity.GetCharData().Name, direction)
					}
				}
			}
			for _, e := range db.GetEntitiesInRoom(to_room.Id) {
				if e != nil {
					e.Send("\r\n&P%s&d opens the door to the %s.\r\n", entity.GetCharData().Name, reverse_direction)
				}
			}
		}
	} else {
		entity.Send("\r\nYou can't close a door that doesn't exist.\r\n")
	}
}

func do_close(entity Entity, args ...string) {
	if entity_unspeakable_state(entity) {
		entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
		return
	}
	if entity.GetCharData().State == ENTITY_STATE_SITTING {
		entity.Send("\r\nYou are unable to move while sitting.\r\n")
		return
	}
	db := DB()
	room := db.GetRoom(entity.RoomId())
	if len(args) == 0 {
		entity.Send("\r\n&ROpen what?&d\r\n")
		return
	}
	// TODO if args[0] is 'hatch' close the spaceship hatch/ramp.
	// For now we'll assume it's a direction door.
	direction := strings.ToLower(args[0])
	if !room.HasExit(direction) {
		entity.Send("\r\n&ROpen what?&d.\r\n")
	}
	if flags, ok := room.ExitFlags[direction]; ok {
		locked, closed := room_get_blocked_exit_flags(flags.(map[string]interface{}))
		if locked {
			// TODO: Search users inventory for the key.
			entity.Send("\r\n&RYou don't have the key.&d\r\n")
			return
		}
		if !closed && !locked {
			entity.Send("\r\n&GYou close the door.&d\r\n")
			room.ExitFlags[direction].(map[string]interface{})["closed"] = true
			to_room := db.GetRoom(room.Exits[direction])
			reverse_direction := direction_reverse(direction)
			if room.Id == to_room.Exits[reverse_direction] {
				if _, ok := to_room.ExitFlags[reverse_direction]; ok {
					to_room.ExitFlags[reverse_direction].(map[string]interface{})["closed"] = true
				}
			}
			for _, e := range db.GetEntitiesInRoom(entity.RoomId()) {
				if e != nil {
					if e != entity {
						e.Send("\r\n&P%s&d closes the door to the %s.\r\n", entity.GetCharData().Name, direction)
					}
				}
			}
			for _, e := range db.GetEntitiesInRoom(to_room.Id) {
				if e != nil {
					e.Send("\r\n&P%s&d closes the door to the %s.\r\n", entity.GetCharData().Name, reverse_direction)
				}
			}
		}
	}
}
