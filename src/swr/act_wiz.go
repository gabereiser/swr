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
func do_room_stat(entity Entity, args ...string) {
	if len(args) > 2 {
		entity.Send("\r\nSyntax: rstat <vnum?> <shipId?>     | vnum and shipId are optional, but vnum must be supplied with shipId.\r\n")
		return
	}
	room := entity.GetRoom()
	vnum := room.Id
	shipId := room.ship
	if len(args) == 2 {
		value, _ := strconv.Atoi(args[0])
		vnum = uint(value)
		value, _ = strconv.Atoi(args[1])
		shipId = uint(value)
	}
	if len(args) == 1 {
		value, _ := strconv.Atoi(args[0])
		vnum = uint(value)
	}
	room = DB().GetRoom(vnum, shipId)
	ship := "None"

	if shipId > 0 {
		s := DB().GetShip(shipId)
		ship = s.GetData().Name + " (" + s.GetData().Type + ")"
	}
	entity.Send("\r\n%s\r\n", MakeTitle("Room Stat", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
	entity.Send("     &GName: &W\"%s\"&d\r\n", room.Name)
	entity.Send("     &GVNum: &W%-7d &GShip: &W%s&d\r\n", room.Id, ship)
	entity.Send("     &GArea: &W%s&d\r\n", room.Area.Name)
	entity.Send("    &GFlags: &W%v&d\r\n", room.Flags)
	entity.Send("     &GDesc: &W\"%s\"&d\r\n", room.Desc)
	entity.Send("    &GExits: &W%+v&d\r\n", room.Exits)
	entity.Send("&GExitFlags: &W%+v&d\r\n", room.ExitFlags)
	entity.Send("&GRoomProgs: &d\r\n")
	for name, value := range room.RoomProgs {
		entity.Send("&y%s&w:%s&d\r\n", name, value)
	}

}
func do_room_set(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if !entity.IsPlayer() {
		return
	}
	player := entity.(*PlayerProfile)
	room := player.GetRoom()
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
	if entity == nil {
		return
	}
	if len(args) != 3 {
		entity.Send("\r\nSyntax: ocreate <filename> <itemtype> <item name>\r\n")
		return
	}
	room := DB().GetRoom(entity.RoomId(), entity.ShipId())
	if room.ship > 0 {
		entity.Send("\r\n&RItem's must be created on planets. So we can use the area name in the filepath.&d\r\n&xsorry...&d\r\n")
		return
	}
	filename := args[0]
	// no need to use .yml in the filename, the path of the item will be /data/items/<area>/<file_name>.yml
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		filename = strings.TrimSuffix(filename, ".yaml")
		filename = strings.TrimSuffix(filename, ".yml")
	}
	itemtype := args[1]
	if !item_is_item_type(itemtype) {
		entity.Send("\r\n&RInvalid item type.&d\r\n")
		return
	}
	itemname := strings.TrimSpace(strings.Join(args[2:], " "))
	words := strings.Split(strings.ToLower(itemname), " ")
	for i, s := range words {
		if len(s) < 4 {
			words[i] = ""
			// simple words like a, or an, or the aren't good keywords.
			// You can use the itemtype as a keyword to add a proper one after creation.
		}
	}
	keywords := strings.Split(strings.Join(words, " "), " ")
	item := new(ItemData)
	item.Id = DB().GetNextItemVnum()
	item.OId = DB().GetNextItemVnum()
	item.Name = itemname
	item.Keywords = make([]string, 0)
	item.Keywords = append(item.Keywords, itemtype)
	item.Keywords = append(item.Keywords, keywords...)

	item.Desc = "A lob of goo."
	item.Type = itemtype
	item.Filename = sprintf("data/items/%s/%s.yml", strings.ToLower(strings.ReplaceAll(room.Area.Name, " ", "")), strings.ToLower(filename))
	DB().SaveItem(item)
	DB().LoadItem(item.Filename)
	entity.Send("\r\n&YObject Create. Ok.&d\r\n")
	room.AddItem(item_clone(item))
}

func do_item_set(entity Entity, args ...string) {
	room := DB().GetRoom(entity.RoomId(), entity.ShipId())
	if len(args) < 3 {
		entity.Send("\r\nSyntax: oset <item> <field> <value>\r\n")
		return
	}
	item := room.FindItem(args[0])
	if item == nil {
		entity.Send("\r\n&RUnable to find item &W%s&R!!&d\r\n", args[0])
		return
	}
	i := item.GetData()
	switch strings.ToLower(args[1]) {
	case "name":
		i.Name = strings.TrimSpace(strings.Join(args[2:], " "))
	case "desc":
		i.Desc = consolify(strings.TrimSpace(strings.Join(args[2:], " ")))
	case "type":
		if !item_is_item_type(args[2]) {
			entity.Send("\r\n&RInvalid type.&d\r\n")
			return
		}
		i.Type = args[2]
	case "keywords":
		switch args[2] {
		case "add":
			i.Keywords = append(i.Keywords, args[3])
		case "rm":
			ki := -1
			for i, k := range i.Keywords {
				if k == args[3] {
					ki = i
				}
			}
			if ki > -1 {
				ret := make([]string, 0)
				ret = append(ret, i.Keywords[:ki]...)
				ret = append(ret, i.Keywords[ki+1:]...)
				i.Keywords = ret
			}
		default:
			entity.Send("\r\nSyntax oset <item> keywords <add/rm> <keyword>\r\n")
			return
		}
	case "value":
		value, _ := strconv.Atoi(args[2])
		i.Value = value
	case "wearloc":
		if !item_is_wearable_slot(args[2]) {
			entity.Send("\r\n&RInvalid wearLoc.&d\r\n")
			return
		}
		i.WearLoc = &args[2]
	case "weapontype":
		if !item_is_weapon_type(args[2]) {
			entity.Send("\r\n&RInvalid weapon type.&d\r\n")
			return
		}
		i.WeaponType = &args[2]
	case "weight":
		value, _ := strconv.Atoi(args[2])
		i.Weight = value
	case "ac":
		value, _ := strconv.Atoi(args[2])
		i.AC = value
	default:
		entity.Send("\r\nSyntax: oset <item> <field> <value>\r\n")
		entity.Send("--------------------------------------------\r\n")
		entity.Send("Fields are:\r\n")
		entity.Send("name, desc, type, keywords, value, wearLoc, weaponType, weight, ac\r\n")
		return
	}
	i.Id = i.OId
	DB().SaveItem(i)
	DB().items[i.Id] = i
	entity.Send("\r\nObject Set. Ok.&d\r\n")
}

func do_item_remove(entity Entity, args ...string) {

}

func do_item_stat(entity Entity, args ...string) {
	if len(args) != 1 {
		entity.Send("\r\nSyntax: ostat <item>\r\n")
		return
	}
	room := DB().GetRoom(entity.RoomId(), entity.ShipId())
	var item Item
	for _, i := range room.Items {
		for _, k := range i.GetData().Keywords {
			if strings.HasPrefix(k, args[0]) {
				item = i
			}
		}
	}
	if item != nil {
		i := item.GetData()
		entity.Send("\r\n%s\r\n", MakeTitle("Object Stats", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
		entity.Send("&GObject Name: &W%s&d", i.Name)
		entity.Send("&GObject   ID: &W%-9d &GOID: &W%-9d&d\r\n", i.Id, i.OId)
		entity.Send("&GObject Desc: &W%s&d\r\n", i.Desc)
		entity.Send("--------------------------------------------")
		entity.Send("&GType: &W%s &GValue: &W%-6d&d\r\n", i.Type, i.Value)
		weaponType := ""
		isWeapon := " "
		if i.WeaponType != nil {
			weaponType = *i.WeaponType
			isWeapon = "x"
		}
		entity.Send("&G  IsWeapon: [%s]    Weapon Type: &W%s&d\r\n", isWeapon, weaponType)
		wearLocation := ""
		isWearable := " "
		if i.WearLoc != nil {
			wearLocation = *i.WearLoc
			isWearable = "x"
		}
		entity.Send("&GIsWearable: [%s]  Wear Location: &W%s&d\r\n", isWearable, wearLocation)
		entity.Send("&GWeight: &W%5d&G kg&d\r\n", i.Weight)

	}
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
	room := target.GetRoom()
	room.SendToOthers(target, sprintf("\r\n%s has left.\r\n", target.GetCharData().Name))
	target.GetCharData().Room = uint(room_id)
	room = DB().GetRoom(uint(room_id), 0)
	room.SendToOthers(target, sprintf("\r\n%s has appeared.\r\n", target.GetCharData().Name))
	target.Send("\r\nYou feel a rush of air as your surroundings quickly change.\r\n")
}

func do_advance(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if len(args) == 0 {
		// we are advancing ourselves...
		ch := entity.GetCharData()
		for i := ch.Level; i < 100; i++ {
			entity_advance_level(entity)
		}
		log.Printf("ADMIN (ADVANCE): %s has been advanced a level!", ch.Name)

	} else {
		if len(args) > 2 {
			entity.Send("\r\nSyntax: advance <charactername> <level>")
			return
		} else {
			l, e := strconv.Atoi(args[1])
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
			if p.GetCharData().Level > uint(l) {
				entity.Send("\r\n&RCharacter is already level &W%d&R!&d\r\n", p.GetCharData().Level)
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
	room := player.GetRoom()
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
