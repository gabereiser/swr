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
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
)

func do_nothing(entity Entity, args ...string) {
	entity.Send("\r\n&rInput not recognized.&d\r\n")
}

func do_password(entity Entity, args ...string) {
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		if len(args) != 3 {
			player.Send("\r\nSyntax: password <oldpassword> <newpassword> <repeat newpassword>\r\n")
		} else {
			oldp := args[0]
			if encrypt_string(oldp) == player.Password {
				if strings.EqualFold(args[1], args[2]) {
					player.Password = encrypt_string(args[1])
					player.Send("\r\n&YPassword. Ok.&d\r\n")
				} else {
					player.Send("\r\n&RPassword Mis-match!&d\r\n")
				}
			} else {
				player.Send("\r\n&RPassword incorrect!&d\r\n")
			}
		}
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
	if entity_unspeakable_state(entity) {
		entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
		return
	}
	if entity.IsPlayer() {
		ch := entity.GetCharData()
		player := entity.(*PlayerProfile)
		if ch.State == ENTITY_STATE_DEAD {
			return
		}
		if ch.State == ENTITY_STATE_SLEEPING || ch.State == ENTITY_STATE_UNCONSCIOUS {
			entity.Send("\r\n&cIn your dreams?...&d\r\n")
			return
		}
		if len(args) == 0 { // l or look with no args
			roomId := entity.RoomId()
			shipId := entity.ShipId()
			ship := entity.GetShip()
			room := DB().GetRoom(roomId, shipId)
			if room != nil {
				if player.Priv >= 100 {
					entity.Send(fmt.Sprintf("\r\n%s\r\n",
						MakeTitle(sprintf("%s [%d]", room.Name, room.Id),
							ANSI_TITLE_STYLE_NORMAL,
							ANSI_TITLE_ALIGNMENT_CENTER)))
				} else {
					entity.Send(fmt.Sprintf("\r\n%s\r\n",
						MakeTitle(room.Name,
							ANSI_TITLE_STYLE_NORMAL,
							ANSI_TITLE_ALIGNMENT_CENTER)))
				}
				entity.Send(sprintf("&W%s&d\r\n\r\n", StitchParagraphs(telnet_encode(room.Desc), build_map(room))))
				entity.Send("Exits: \r\n")
				for dir, to_room_id := range room.Exits {
					to_room := DB().GetRoom(to_room_id, shipId)
					if k, ok := room.ExitFlags[dir]; ok {
						exit_flags := k
						ext := room_get_exit_status(exit_flags)
						entity.Send(sprintf("&G%s&W - &Y[&W%s&Y] &C%s&d\r\n", capitalize(dir), to_room.Name, ext))
					} else {
						entity.Send(sprintf("&G%s&W - &Y[&W%s&Y]&d\r\n", capitalize(dir), to_room.Name))
					}
				}
				entity.Send("\r\n")
				for i := range room.Items {
					item := room.Items[i]
					if item == nil {
						continue
					}
					if item.IsCorpse() {
						entity.Send("%s   &w%s %s&d\r\n", EMOJI_TOMBSTONE, item.GetData().Name)
					} else {
						entity.Send("&w%s %s&d\r\n", item.GetData().Name)
					}

				}

				for _, e := range room.GetEntities() {
					if e != entity {
						entity.Send("&P%s&d\r\n", e.GetCharData().Name)
					}
				}
				if shipId > 0 {
					if ship.GetData().Cockpit == roomId || (ship.GetData().Ramp == roomId && !ship.GetData().InSpace) {
						room := DB().GetRoom(ship.GetData().LocationId, 0)
						if room != nil {
							entity.Send(fmt.Sprintf("\r\nThrough your ships viewscreen you see...\r\n\r\n%s\r\n",
								MakeTitle(room.Name,
									ANSI_TITLE_STYLE_NORMAL,
									ANSI_TITLE_ALIGNMENT_CENTER)))
							entity.Send(sprintf("&W%s&d\r\n", StitchParagraphs(telnet_encode(room.Desc), build_map(room))))
							for dir, to_room_id := range room.Exits {
								to_room := DB().GetRoom(to_room_id, shipId)
								if k, ok := room.ExitFlags[dir]; ok {
									exit_flags := k
									ext := room_get_exit_status(exit_flags)
									entity.Send(sprintf("&G%s&W - &Y(&W%s&Y) &C%s&d\r\n", capitalize(dir), to_room.Name, ext))
								} else {
									entity.Send(sprintf("&G%s&W - &Y(&W%s&Y)&d\r\n", capitalize(dir), to_room.Name))
								}
							}
							entity.Send("\r\n")
							for i := range room.Items {
								item := room.Items[i]
								if item == nil {
									continue
								}
								if item.IsCorpse() {
									entity.Send("%s   &w%s&d\r\n", EMOJI_TOMBSTONE, item.GetData().Name)
								} else {
									entity.Send("&w%s&d\r\n", item.GetData().Name)
								}

							}

							for _, e := range room.GetEntities() {
								if e != entity {
									entity.Send("&P%s&d\r\n", e.GetCharData().Name)
								}
							}
						}
					}
				}
				return
			} else {
				log.Fatalf("Entity %s is in room %d, only it doesn't exist and crashed the server.", entity.GetCharData().Name, entity.RoomId())
			}

		} else {
			room := entity.GetRoom()
			for _, e := range room.GetEntities() {
				if e != entity {
					ch := e.GetCharData()
					for _, keyword := range ch.Keywords {
						if strings.HasPrefix(strings.ToLower(keyword), strings.ToLower(args[0])) {
							entity.Send("You look at %s and see...\r\n%s\r\n", ch.Name, ch.Desc)
							return
						}
					}
				}
			}
			item := room.FindItem(args[0])
			if item != nil {
				entity.Send("You look at %s and see...\r\n%s\r\n", item.GetData().Name, item.GetData().Desc)
				return
			}
			item = entity.FindItem(args[0])
			if item != nil {
				entity.Send("You look at %s and see...\r\n%s\r\n", item.GetData().Name, item.GetData().Desc)
				return
			}
		}
	}
	entity.Send("\r\n&dCan't find that here.\r\n")
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
	if entity.IsFighting() {
		entity.Send("\r\nYou are fighting!\r\n")
		return
	}
	db := DB()
	room := db.GetRoom(entity.RoomId(), entity.ShipId())
	if !room.HasExit(direction) {
		entity.Send("\r\nYou can't go that way.\r\n")
		return
	} else {
		to_room := db.GetRoom(room.Exits[direction], entity.ShipId())
		if to_room == nil {
			entity.Send("\r\n&RThat room doesn't exist!\r\n")
			return
		} else {
			locked := false
			closed := false
			flags := room.GetExitFlags(direction)
			if flags != nil {
				locked, closed = room_get_blocked_exit_flags(flags)
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
				for _, e := range room.GetEntities() {
					if entity_unspeakable_state(e) {
						continue
					}
					if e != entity {
						e.Send("\r\n%s has left going %s.\r\n", entity.GetCharData().Name, direction)
						if e.GetCharData().AI != nil {
							e.GetCharData().AI.OnMove(entity)
						}
					}
				}
				go room_prog_exec(entity, "leave", direction)
				entity.GetCharData().Room = to_room.Id
				do_look(entity)
				go room_prog_exec(entity, "enter", direction_reverse(direction))
				for _, e := range to_room.GetEntities() {
					if entity_unspeakable_state(e) {
						continue
					}
					if e != entity {
						e.Send("\r\n%s has arrived from the %s.\r\n", entity.GetCharData().Name, direction_reverse(direction))
						if e.GetCharData().AI != nil {
							e.GetCharData().AI.OnGreet(entity)
						}
					}
				}
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
		entity.GetRoom().SendToOthers(entity, sprintf("\r\n&d%s stands up.\r\n", ch.Name))
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
		entity.GetRoom().SendToOthers(entity, sprintf("\r\n&d%s sits down.\r\n", ch.Name))
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
	entity.GetRoom().SendToOthers(entity, sprintf("\r\n&d%s lays down and falls asleep.\r\n", ch.Name))
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
	room := db.GetRoom(entity.RoomId(), entity.ShipId())
	if len(args) == 0 {
		entity.Send("\r\n&ROpen what?&d\r\n")
		return
	}
	direction := get_direction_string(strings.ToLower(args[0]))
	if !room.HasExit(direction) {
		entity.Send("\r\n&ROpen what?&d.\r\n")
		return
	}
	room.OpenDoor(entity, direction, false)
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
	room := db.GetRoom(entity.RoomId(), entity.ShipId())
	if len(args) == 0 {
		entity.Send("\r\n&RClose what?&d\r\n")
		return
	}
	// TODO if args[0] is 'hatch' close the spaceship hatch/ramp.
	// For now we'll assume it's a direction door.
	direction := get_direction_string(strings.ToLower(args[0]))
	if !room.HasExit(direction) {
		entity.Send("\r\n&RClose what?&d.\r\n")
		return
	}
	room.CloseDoor(entity, direction, false)
}

func do_get(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\n&RGet what?&d\r\n")
		return
	}
	db := DB()
	ch := entity.GetCharData()
	if len(args) == 1 {
		// fetch an item from the room
		room := db.GetRoom(ch.Room, ch.Ship)
		item := room.FindItem(args[0])
		if item == nil {
			entity.Send("\r\n&dCan't seem to find that.\r\n")
		} else {
			if !entity_pickup_item(entity, item) {
				return
			}
			room.RemoveItem(item)
			room.SendToOthers(entity, sprintf("\r\n&P%s&d picks up &Y%s&d.\r\n", ch.Name, item.GetData().Name))
			entity.Send("\r\n&dYou pick up &Y%s&d.\r\n", item.GetData().Name)
			go room_prog_exec(entity, "get", item) // indiana jones...
			return
		}
	}
	// get <item> from
	if len(args) == 2 {
		entity.Send("\r\n&RGet &Y%s&W from &Rwhere?&d\r\n", args[0])
	}
	// get <item> from <item>
	if len(args) == 3 {

		from_container := args[2]
		item_name := args[0]
		ch := entity.GetCharData()
		room := db.GetRoom(ch.Room, ch.Ship)
		//on your person (backpack, bag)
		item := entity.FindItem(from_container)
		if item == nil {
			// in the room (corpse, continer)
			item = room.FindItem(from_container)
		}
		if item == nil {
			entity.Send("\r\nCan't seem to find that.\r\n")
			return
		} else {
			if item.IsContainer() || item.IsCorpse() {
				i := item.GetData().FindItemInContainer(item_name)
				if i == nil {
					entity.Send("\r\nCan't seem to find that in %s.\r\n", item.GetData().Name)
					return
				} else {
					if !entity_pickup_item(entity, i) {
						return
					}
					item.GetData().RemoveItem(i)
					entity.Send("\r\n&dYou pick up &Y%s&d from &Y%s&d.\r\n", i.GetData().Name, item.GetData().Name)
					return
				}
			}
		}
	}
	if len(args) > 3 {
		entity.Send("\r\n&CSyntax: &dget <item> | from <container>\r\n")
	}

}

func do_give(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\nSyntax: give <entity> <item>\r\n")
		entity.Send("-----------------------------------------------------\r\n")
		entity.Send("To give credits, use: give <entity> <quantity> credits\r\n")
		return
	}
	if len(args) == 1 {
		entity.Send("\r\n&RGive to who?&d\r\n")
		return
	}
	entity_name := args[0]
	item_name := "credits"
	quantity := 1
	if len(args) == 2 {
		item_name = args[1]
	}
	if len(args) == 3 {
		q, e := strconv.Atoi(args[1])
		if e != nil {
			entity.Send("\r\n&RUnable to determine quantity of credits!&d\r\n")
			return
		}
		quantity = q
	}
	var target Entity
	for _, e := range entity.GetRoom().GetEntities() {
		if e != entity {
			for _, k := range e.GetCharData().Keywords {
				if strings.HasPrefix(k, entity_name) {
					target = e
					break
				}
			}
			if target != nil {
				break
			}
		}
	}
	if target == nil {
		entity.Send("\r\n&RUnable to find &W%s&R!&d\r\n", entity_name)
		return
	}
	if item_name == "credits" {
		uq := uint(quantity)
		if entity.GetCharData().Gold < uq {
			entity.Send("\r\n&RNot enough credits!&d\r\n", entity_name)
			return
		}
		target.GetCharData().Gold += uq
		entity.GetCharData().Gold -= uq

		target.Send("\r\n&P%s&Y has given you &w%d&Y credits.&d\r\n", entity.GetCharData().Name, uq)
		if !target.IsPlayer() {
			if target.GetCharData().AI != nil {
				target.GetCharData().AI.OnGive(entity, quantity, nil)
			}
		}
		entity.Send("\r\n&YYou give &P%s&Y &w%d&Y credits.&d\r\n", target.GetCharData().Name, uq)
		return
	}
	item := entity.FindItem(item_name)
	if item == nil {
		entity.Send("\r\n&RUnable to find &W%s&R!&d\r\n", item_name)
		return
	}
	if item.GetData().Type == ITEM_TYPE_KEY {
		entity.Send("\r\n&RUnable to find &W%s&R!&d\r\n", item_name) // ;)
		return
	}
	if !entity_pickup_item(target, item) {
		entity.Send("\r\n&RThey are unable to carry &W%s&R!&d\r\n", item_name)
		return
	}
	entity.GetCharData().RemoveItem(item)
	target.Send("\r\n&P%s&Y has given you &W%s&Y.&d\r\n", entity.GetCharData().Name, item.GetData().Name)
	entity.Send("\r\n&YYou give &P%s&Y &W%s&Y.&d\r\n", target.GetCharData().Name, item.GetData().Name)
	if !target.IsPlayer() {
		if target.GetCharData().AI != nil {
			target.GetCharData().AI.OnGive(entity, 1, item)
		}
	}

}
func do_put(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\n&RPut what in what?&d\r\n")
		return
	}
	if len(args) == 3 {
		item_name := args[0]
		container_name := args[2]
		db := DB()
		room := db.GetRoom(entity.RoomId(), entity.ShipId())
		item := entity.FindItem(item_name)
		if item == nil {
			entity.Send("\r\n&RCan't seem to find that.&d\r\n")
			return
		}
		container := entity.FindItem(container_name)
		if container == nil {
			container = room.FindItem(container_name)
		}
		if container == nil {
			entity.Send("\r\nCan't seem to find that container.\r\n")
			return
		}
		container.GetData().AddItem(item)

	} else {
		entity.Send("\r\n&CSyntax: put <item> in <container>.&d\r\n")
	}
}
func do_drop(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\n&RDrop what?&d\r\n")
		return
	}
	db := DB()
	item_name := args[0]
	item := entity.FindItem(item_name)
	if item == nil {
		entity.Send("\r\nCan't find that in your inventory.\r\n")
		return
	}
	room := db.GetRoom(entity.RoomId(), entity.ShipId())
	room.AddItem(item)
	entity.GetCharData().RemoveItem(item)
	entity.Send("\r\n&YYou drop &W%s&Y.&d\r\n", item.GetData().Name)
	ch := entity.GetCharData()
	for _, e := range room.GetEntities() {
		if e == nil {
			continue
		}
		if e != entity {
			e.Send("\r\n&P%s&d drops &Y%s&d.\r\n", ch.Name, item.GetData().Name)
			if e.GetCharData().AI != nil {
				e.GetCharData().AI.OnDrop(entity, item)
			}
		}
	}
	go room_prog_exec(entity, "drop", item)
}

func do_statsys(entity Entity, args ...string) {
	db := DB()
	mobCount := 0
	playerCount := 0
	for _, e := range db.entities {
		if e == nil {
			continue
		}
		if !e.IsPlayer() {
			mobCount++
		} else {
			playerCount++
		}
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	entity.Send("\r\n%s\r\n", MakeTitle("System Stats", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
	entity.Send("&G      System Name:&W %s\r\n", Config().Name)
	entity.Send("&G    Total Systems: &W%-3d       &GTotal Areas: &W%-3d&d\r\n", len(db.starsystems), len(db.areas))
	entity.Send("&G       Total Mobs: &W%-12d &GTotal Rooms: &W%-12d&d\r\n", mobCount, len(db.rooms))
	entity.Send("&G      Total Ships: &W%-4d&d\r\n", len(db.ships))
	entity.Send("\r\n%s\r\n", MakeTitle("OS", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
	entity.Send("&Y           Name&d: %s\r\n", runtime.GOOS)
	entity.Send("&Y           Arch&d: %s\r\n", runtime.GOARCH)
	entity.Send("\r\n%s\r\n", MakeTitle("CPU", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
	entity.Send("&Y          Cores&d: %d\r\n", runtime.NumCPU())
	entity.Send("&Y  Total Threads&d: %d\r\n", runtime.NumGoroutine())
	entity.Send("\r\n%s\r\n", MakeTitle("Memory", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
	entity.Send("&Y Current Memory&d: %.4f mb\r\n", bytes_to_mb(m.Alloc))
	entity.Send("&YReserved Memory&d: %.4f mb\r\n", bytes_to_mb(m.Sys))
	entity.Send("&Y    Misc Memory&d: %.4f mb\r\n", bytes_to_mb(m.OtherSys))
	entity.Send("&Y       GC Count&d: %d\r\n", m.NumGC)
	entity.Send("\r\n")
}
