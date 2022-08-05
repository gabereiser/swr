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
	"log"
	"strconv"
	"strings"
)

func do_area_create(entity Entity, args ...string) {
	if !entity.IsPlayer() {
		return
	}
	if len(args) < 3 {
		entity.Send("\r\nSyntax: acreate <areaname> <min vnum> <max vnum>\r\n")
		entity.Send("-----------------------------------------------------------\r\n")
		entity.Send("*NOTE* area create will create up to max vnum, but not including.\r\n")
		entity.Send("Using 100 200 will create 99 rooms starting at 100.\r\n")
		return
	}
	min_vnum, err := strconv.ParseInt(args[1], 10, 32)
	ErrorCheck(err)
	max_vnum, err := strconv.ParseInt(args[2], 10, 32)
	ErrorCheck(err)
	db := DB()
	for _, area := range db.areas {
		for _, r := range area.Rooms {
			if uint(min_vnum) < r.Id && r.Id < uint(max_vnum) {
				entity.Send("\r\n&RError! Vnum range already exists!&d\r\n")
				return
			}
		}
	}
	area := new(AreaData)
	area.Name = args[0]
	area.Author = entity.GetCharData().Name
	area.Rooms = make([]RoomData, 0)
	area.Items = make([]ItemSpawn, 0)
	area.Mobs = make([]MobSpawn, 0)
	area.Levels = []uint16{1, 100}
	area.Reset = 300
	area.ResetMsg = "The world seems to shift around you."
	for i := min_vnum; i < max_vnum; i++ {
		room := RoomData{
			Id:        uint(i),
			Name:      "A void",
			Desc:      "Somewhere in the void of space.",
			Flags:     make([]string, 0),
			Items:     make([]Item, 0),
			Exits:     make(map[string]uint),
			ExitFlags: make(map[string]*RoomExitFlag),
			Area:      area,
			RoomProgs: make(map[string]string),
		}
		area.Rooms = append(area.Rooms, room)
	}
	db.SaveArea(area)
	entity.Send("\r\n&YArea Create. Ok.&d\r\n")
}

func do_area_set(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		room := DB().GetRoom(player.RoomId(), player.ShipId())
		if room.Area == nil {
			entity.Send("\r\n&RNot in an area!&d\r\n")
			return
		}
		if len(args) < 2 {
			entity.Send("\r\nSyntax aset <field> <value>\r\n")
			entity.Send("-------------------------------------\r\n")
			entity.Send("Available Fields:\r\n")
			entity.Send("name, levels, author, reset, resetMsg")
			return
		}
		switch strings.ToLower(args[0]) {
		case "name":
			room.Area.Name = strings.TrimSpace(strings.Join(args[1:], " "))
		case "levels":
			if len(args) != 3 {
				entity.Send("\r\nSyntax: aset levels <min> <max>\r\n")
				return
			} else {
				min, _ := strconv.Atoi(args[1])
				max, _ := strconv.Atoi(args[2])
				room.Area.Levels[0] = uint16(min)
				room.Area.Levels[1] = uint16(max)
			}
		case "author":
			room.Area.Author = strings.TrimSpace(strings.Join(args[1:], " "))
		case "reset":
			r, _ := strconv.Atoi(args[1])
			room.Area.Reset = uint(r)
		case "resetmsg":
			room.Area.ResetMsg = strings.TrimSpace(strings.Join(args[1:], " "))
		default:
			entity.Send("\r\n&RInvalid field.&d\r\n")
		}
	}
	entity.Send("\r\n&YArea Set. Ok&d\r\n")
}

func do_area_remove(entity Entity, args ...string) {
	for i, area := range DB().areas {
		if strings.EqualFold(area.Name, args[0]) {
			DB().RemoveArea(DB().areas[i])
			entity.Send("\r\n&YArea Remove. Ok.&d\r\n")
			return
		}
	}
	entity.Send("\r\n&RArea not found.&d\r\n")
}

func do_area_reset(entity Entity, args ...string) {
	if len(args) == 0 {
		DB().ResetAll()
	} else {
		for i, area := range DB().areas {
			if strings.EqualFold(area.Name, args[0]) {
				area_reset(DB().areas[i])
				entity.Send("\r\n&YArea Reset. Ok.&d\r\n")
				return
			}
		}
		entity.Send("\r\n&RArea not found.&d\r\n")
	}
}
func do_area_save(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if entity.IsPlayer() {
		room := DB().GetRoom(entity.RoomId(), entity.ShipId())
		if room.Area != nil {
			DB().SaveArea(room.Area)
			entity.Send("\r\n&YArea Save. Ok.&d\r\n")
		} else {
			entity.Send("\r\n&RNot in an area file!&d\r\n")
		}
	}
}

func do_room_create(entity Entity, args ...string) {

}

func do_room_set(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if !entity.IsPlayer() {
		return
	}
	player := entity.(*PlayerProfile)
	room := DB().GetRoom(player.RoomId(), player.ShipId())
	if len(args) < 2 {
		entity.Send("\r\nSyntax rset <field> <value>\r\n")
		entity.Send("-------------------------------------\r\n")
		entity.Send("Available Fields:\r\n")
		entity.Send("name, desc, flags")
		return
	}
	switch args[0] {
	case "name":
		room.Name = strings.TrimSpace(strings.Join(args[1:], " "))
	case "desc":
		room.Desc = consolify(strings.TrimSpace(strings.Join(args[1:], " ")))
	case "flags":
		if room.HasFlag(args[1]) {
			room.RemoveFlag(args[1])
		} else {
			room.SetFlag(args[1])
		}
	default:
		entity.Send("\r\n&RField invalid.&d\r\n")
		return
	}
	// Set the data on the AreaRoom []Rooms slice so that when the area is saved, the changes to the room are too.
	for i, r := range room.Area.Rooms {
		if r.Id == room.Id {
			room.Area.Rooms[i] = *room
		}
	}
	entity.Send("\r\n&YSet. Ok.&d\r\n")
}

func do_room_make_exit(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if len(args) != 2 {
		entity.Send("\r\nSyntax: rexit <dir> <roomId>\r\n")
		entity.Send("--------------------------------------\r\n")
		entity.Send("*NOTE* Rooms must be on the same ship/planet.\r\n")
		entity.Send("Rooms cannot be joined across the galaxy.\r\n")
		entity.Send("To delete an exit, supply roomId \"0\".\r\n")
		return
	}
	dir := get_direction_string(args[0])
	vnum, _ := strconv.Atoi(args[1])
	room := DB().GetRoom(entity.RoomId(), entity.ShipId())
	if room == nil {
		entity.Send("\r\n&RFATAL! Unable to determine your room!!!&d\r\n")
		log.Printf("FATAL: Unable to determine room for roomId %d and locationId %d", entity.RoomId(), entity.ShipId())
		return
	}
	if vnum == 0 {
		delete(room.Exits, dir)
		for i, r := range room.Area.Rooms {
			if r.Id == room.Id {
				room.Area.Rooms[i] = *room
			}
		}
		entity.Send("\r\n&YExit. Ok&d\r\n")
		return
	}
	to_room := DB().GetRoom(uint(vnum), entity.ShipId())
	if to_room == nil {
		entity.Send("\r\n&RFATAL! Unable to find exit room!!!&d\r\n")
		log.Printf("FATAL: Unable to determine room for roomId %d and locationId %d", vnum, entity.ShipId())
		return
	}
	if room.ship != to_room.ship {
		entity.Send("\r\n&RFATAL! Rooms are not in the same area!!!&d\r\n")
		return
	}
	room.Exits[dir] = to_room.Id
	entity.Send("\r\n&YExit. Ok.&d\r\n")
}

func do_room_remove(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if len(args) != 1 {
		entity.Send("\r\nSyntax: rremove <vnum>\r\n")
		entity.Send("-----------------------------------------------------------------\r\n")
		entity.Send("*NOTE* rremove will reset a room back to it's prototype state.\r\n")
		entity.Send("Use with caution.\r\n")
		entity.Send("Will reset a room in your current area/ship.\r\n")
		return
	}
	vnum, _ := strconv.Atoi(args[0])
	room := DB().GetRoom(uint(vnum), entity.ShipId())
	if room == nil {
		entity.Send("\r\n&RRoom not found in your area.&d\r\n")
		return
	}
	for dir, rId := range room.Exits {
		eroom := DB().GetRoom(rId, entity.ShipId())
		delete(eroom.Exits, direction_reverse(dir))
	}
	room.Name = "A void"
	room.Desc = "Somewhere in the void of space."
	room.Exits = make(map[string]uint)
	room.ExitFlags = make(map[string]*RoomExitFlag)
	room.Flags = make([]string, 0)
	room.Items = make([]Item, 0)
	room.RoomProgs = make(map[string]string)
	if room.ship > 0 {
		ship := DB().GetShip(room.ship)
		ship.GetData().Rooms[room.Id] = room
	} else {
		for i, r := range room.Area.Rooms {
			if r.Id == room.Id {
				room.Area.Rooms[i] = *room
			}
		}
	}
	entity.Send("\r\n&YRemove. Ok.&d\r\n")

}

func do_item_create(entity Entity, args ...string) {

}

func do_item_set(entity Entity, args ...string) {

}

func do_item_remove(entity Entity, args ...string) {

}

func do_item_stat(entity Entity, args ...string) {

}

func do_item_find(entity Entity, args ...string) {

}

func do_room_find(entity Entity, args ...string) {

}

func do_mob_create(entity Entity, args ...string) {

}

func do_mob_set(entity Entity, args ...string) {

}

func do_mob_remove(entity Entity, args ...string) {

}

func do_mob_reset(entity Entity, args ...string) {

}

func do_mob_stat(entity Entity, args ...string) {

}

func do_mob_find(entity Entity, args ...string) {

}

func do_transfer(entity Entity, args ...string) {
	if len(args) < 2 {
		entity.Send("\r\n&RTransfer who, where?&d\r\nSyntax: transfer <entity_name> <room_id>\r\n")
		return
	}
	entity_name := args[0]
	target := DB().GetPlayerEntityByName(entity_name)
	if target == nil {
		entity.Send("\r\n&RCouldn't find target entity to transfer!&d\r\n")
		return
	}
	room_id, err := strconv.Atoi(args[1])
	if err != nil {
		entity.Send("\r\n&RUnable to parse room_id!&d\r\n")
		return
	}
	room := DB().GetRoom(target.RoomId(), target.ShipId())
	for _, e := range DB().GetEntitiesInRoom(room.Id, target.ShipId()) {
		if e == nil {
			continue
		}
		if e != target {
			e.Send("\r\n&C%s&d has left.\r\n")
			if e.GetCharData().AI != nil {
				e.GetCharData().AI.OnMove(entity)
			}
		}
	}
	target.GetCharData().Room = uint(room_id)
	room = DB().GetRoom(uint(room_id), 0)
	for _, e := range DB().GetEntitiesInRoom(room.Id, 0) {
		if e == nil {
			continue
		}
		if e != target {
			e.Send("\r\n&C%s&d has appeared.\r\n")
			if e.GetCharData().AI != nil {
				e.GetCharData().AI.OnGreet(entity)
			}
		}
	}
}

func do_advance(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if len(args) == 0 {
		// we are advancing ourselves...
		ch := entity.GetCharData()
		for i := ch.Level; i <= 100; i++ {
			entity_advance_level(entity)
		}
		log.Printf("ADMIN (ADVANCE): %s has been advanced a level!", ch.Name)

	} else {
		if len(args) > 2 {
			entity.Send("\r\nSyntax: advance <charactername> <level>")
			return
		} else {
			l, e := strconv.ParseInt(args[1], 10, 32)
			ErrorCheck(e)
			if e != nil {
				entity.Send("\r\n&RUnable to parse <level>&d\r\n")
				return
			}
			p := DB().GetPlayerEntityByName(args[0])
			if p == nil {
				entity.Send("\r\n&RUnable to find player %s", args[0])
				return
			}
			for i := p.GetCharData().Level; i <= uint(l); i++ {
				entity_advance_level(p)
			}
		}
	}
}

func do_dig(entity Entity, args ...string) {
	if len(args) < 2 {
		entity.Send("\r\nSyntax: dig <dir> <room name>\r\n")
		return
	}
	if !entity.IsPlayer() {
		return
	}
	player := entity.(*PlayerProfile)
	if player.Priv != 100 {
		entity.Send("\r\n&ROnly Immortals can dig, dig?&d\r\n")
	}
	db := DB()
	room := db.GetRoom(player.RoomId(), player.ShipId())
	dir := get_direction_string(strings.ToLower(args[0]))
	if _, ok := room.Exits[dir]; ok {
		entity.Send("\r\n&RRoom already exists in that direction!&d\r\n")
	} else {
		lastVnum := uint(0)
		for _, r := range room.Area.Rooms {
			if r.Id > lastVnum {
				lastVnum = r.Id
			}
		}
		next_id := db.GetNextRoomVnum(room.Id, room.ship)
		if next_id == 0 {
			entity.Send("\r\n&RUnable to determine next room vnum.&d\r\n")
			return
		}

		log.Printf("Found next vnum of %d from room %d", next_id, room.Id)

		room.Exits[dir] = next_id
		next_room := db.GetRoom(next_id, room.ship)

		if next_room == nil {
			next_room = &RoomData{
				Id:        next_id,
				ship:      room.ship,
				Name:      strings.TrimSpace(strings.Join(args[1:], " ")),
				Desc:      sprintf("Room Dugged by %s", entity.GetCharData().Name),
				Exits:     make(map[string]uint),
				ExitFlags: make(map[string]*RoomExitFlag),
				Flags:     []string{"prototype"},
				RoomProgs: make(map[string]string),
				Area:      room.Area,
			}
		}

		next_room.Name = strings.TrimSpace(strings.Join(args[1:], " "))
		next_room.Exits[direction_reverse(dir)] = room.Id
		if next_room.ship > 0 {
			ship := db.GetShip(next_room.ship)
			s := ship.GetData()
			s.Rooms[next_id] = next_room
		} else {
			db.rooms[next_id] = next_room
			for i, r := range room.Area.Rooms {
				if r.Id == next_room.Id {
					room.Area.Rooms[i] = *next_room
				}
			}
		}
		entity.Send("\r\n&GDug a room to the %s&d\r\n", dir)
	}
}
