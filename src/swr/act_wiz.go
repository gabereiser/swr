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
			Id:   uint(i),
			Name: "A void",
			Desc: "Somewhere in the void of space.",
		}
		area.Rooms = append(area.Rooms, room)
	}
	db.SaveArea(area)
	entity.Send("\r\n&GDone&d\r\n")
}

func do_area_set(entity Entity, args ...string) {

}

func do_area_remove(entity Entity, args ...string) {
	for i, area := range DB().areas {
		if strings.EqualFold(area.Name, args[0]) {
			DB().RemoveArea(DB().areas[i])
			entity.Send("\r\n&GArea Removed.&d\r\n")
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
				entity.Send("\r\n&YReset. Ok.&d\r\n")
				return
			}
		}
		entity.Send("\r\n&RArea not found.&d\r\n")
	}
}

func do_room_create(entity Entity, args ...string) {

}

func do_room_set(entity Entity, args ...string) {

}

func do_room_edit(entity Entity, args ...string) {

}

func do_room_reset(entity Entity, args ...string) {

}

func do_room_remove(entity Entity, args ...string) {

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
