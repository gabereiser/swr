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
	"os"
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
		entity.Send("\r\nSyntax: areset <areaname>\r\n")
		return
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
		DB().SaveMobs()
		DB().SaveItems()
		room := DB().GetRoom(entity.RoomId(), entity.ShipId())
		if room.Area != nil {
			DB().SaveArea(room.Area)
			entity.Send("\r\n&YArea Save. Ok.&d\r\n")
		} else {
			entity.Send("\r\n&RNot in an area file!&d\r\n")
		}
	}
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
		value, e := strconv.Atoi(args[0])
		ErrorCheck(e)
		if e != nil {
			entity.Send("\r\n&RUnable to parse room vnum.&d\r\n")
			return
		}
		vnum = uint(value)
		value, e = strconv.Atoi(args[1])
		ErrorCheck(e)
		if e != nil {
			entity.Send("\r\n&RUnable to parse ship vnum.&d\r\n")
			return
		}
		shipId = uint(value)
	}
	if len(args) == 1 {
		value, e := strconv.Atoi(args[0])
		ErrorCheck(e)
		if e != nil {
			entity.Send("\r\n&RUnable to parse room vnum.&d\r\n")
			return
		}
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
	entity.Send("   &GSpawns: &d\r\n")
	for _, ms := range room.Area.Mobs {
		if ms.Room != room.Id {
			continue
		}
		mob := DB().GetMob(ms.Mob)
		m := mob.GetCharData()
		entity.Send(sprintf("&Y[&W%d&Y]&d%-26s", m.OId, tstring(m.Name, 23)))
	}
	entity.Send("\r\n")

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
	if len(args) < 2 {
		entity.Send("\r\nSyntax: rexit <dir> <roomId> [closed?]\r\n")
		entity.Send("--------------------------------------\r\n")
		entity.Send("*NOTE* Rooms must be on the same ship/planet.\r\n")
		entity.Send("Rooms cannot be joined across the galaxy.\r\n")
		entity.Send("To delete an exit, supply roomId \"0\".\r\n")
		entity.Send("To close an exit, rexit <dir> <roomId> 1.\r\n")
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
		to_room := DB().GetRoom(room.Exits[dir], entity.ShipId())
		delete(room.Exits, dir)
		delete(to_room.Exits, direction_reverse(dir))
		if room.ExitFlags != nil {
			delete(room.ExitFlags, dir)
		}
		if to_room.ExitFlags != nil {
			delete(to_room.ExitFlags, direction_reverse(dir))
		}
		for i, r := range room.Area.Rooms {
			if r.Id == room.Id {
				room.Area.Rooms[i] = *room
			}
		}
		for i, r := range to_room.Area.Rooms {
			if r.Id == to_room.Id {
				to_room.Area.Rooms[i] = *to_room
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
	to_room.Exits[direction_reverse(dir)] = room.Id
	if len(args) == 3 {
		if args[2] == "1" {
			if room.ExitFlags == nil {
				room.ExitFlags = make(map[string]*RoomExitFlag)
			}
			if to_room.ExitFlags == nil {
				to_room.ExitFlags = make(map[string]*RoomExitFlag)
			}
			room.ExitFlags[dir] = &RoomExitFlag{
				Closed: true,
			}
			to_room.ExitFlags[direction_reverse(dir)] = &RoomExitFlag{
				Closed: true,
			}
		}
	}
	for i, r := range room.Area.Rooms {
		if r.Id == room.Id {
			// copy over the DB.[]rooms room back to the Area's []Rooms
			room.Area.Rooms[i] = *room
		}
	}
	for i, r := range to_room.Area.Rooms {
		if r.Id == to_room.Id {
			// copy over the DB.[]rooms room back to the Area's []Rooms
			to_room.Area.Rooms[i] = *to_room
		}
	}
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
	if len(args) < 3 {
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
		if len(s) < 3 {
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
	item.Keywords = append(item.Keywords, keywords...)

	item.Desc = "A lob of goo."
	item.Type = itemtype
	if itemtype == ITEM_TYPE_CONTAINER {
		item.Desc = "A box of goo."
		item.Items = make([]Item, 0)
	}
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
	if len(args) != 1 {
		entity.Send("\r\nSyntax: oremove <item>\r\n")
		return
	}
	room := entity.GetRoom()
	item := room.FindItem(args[0])
	if item == nil {
		entity.Send("\r\n&RUnable to find item.")
	}
	item_id := item.GetData().OId
	room.RemoveItem(item)
	delete(DB().items, item_id)
	for _, a := range DB().areas {
		for i, isp := range a.Items {
			if isp.Item == item_id {
				// remove the mobspawn that has this mob listed
				ret := make([]ItemSpawn, 0)
				ret = append(ret, a.Items[:i]...)
				ret = append(ret, a.Items[i+1:]...)
				a.Items = ret
			}
		}
	}
}

func do_item_spawn(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\nSyntax: ospawn <item_vnum>\r\n")
		return
	}
	vnum, err := strconv.Atoi(args[0])
	if err != nil {
		entity.Send("\r\n&RUnable to parse argument as number.&d\r\n")
		return
	}
	item := DB().items[uint(vnum)]
	if item == nil {
		entity.Send("\r\n&RUnable to find item with vnum: &W%d&d\r\n", vnum)
		return
	} else {
		room := entity.GetRoom()
		room.Area.Items = append(room.Area.Items, ItemSpawn{
			Item: item.OId,
			Room: room.Id,
		})
	}
	entity.Send("\r\n&YObject Spawn. Ok.&d\r\n")
}

func do_item_stat(entity Entity, args ...string) {
	if len(args) != 1 {
		entity.Send("\r\nSyntax: ostat <item>\r\n")
		return
	}
	room := entity.GetRoom()
	item := room.FindItem(args[0])
	if item != nil {
		i := item.GetData()
		entity.Send("\r\n%s\r\n", MakeTitle("Object Stats", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
		entity.Send("&GFilename: &W%s&d\r\n", i.Filename)
		entity.Send("&GID: &W%-9d &GOID: &W%-9d&d\r\n", i.Id, i.OId)
		entity.Send("&GName: &W%s&d\r\n", i.Name)
		entity.Send("&GDesc: &W%s&d\r\n", i.Desc)
		entity.Send("&c--------------------------------------------&d\r\n")
		entity.Send("&GType: &W%s &GValue: &W%-6d&d\r\n", i.Type, i.Value)
		entity.Send("&GWeight: &W%5d&G kg&d\r\n", i.Weight)
		weaponType := ""
		isWeapon := " "
		if i.WeaponType != nil {
			weaponType = *i.WeaponType
			isWeapon = "x"
		}
		entity.Send("&G   IsWeapon: [%s]    Weapon Type: &W%s&d\r\n", isWeapon, weaponType)
		wearLocation := ""
		isWearable := " "
		if i.WearLoc != nil {
			wearLocation = *i.WearLoc
			isWearable = "x"
		}
		isContainer := " "
		if i.Type == ITEM_TYPE_CONTAINER {
			isContainer = "x"
		}
		entity.Send("&G IsWearable: [%s]  Wear Location: &W%s&d\r\n", isWearable, wearLocation)
		entity.Send("&GIsContainer: [%s]&d\r\n", isContainer)
		if i.Type == ITEM_TYPE_CONTAINER {
			for _, i := range i.Items {
				entity.Send("&Y[&W%d&Y]&w%s&d\r\n", i.GetData().Name)
			}
		}

	}
}

func do_item_find(entity Entity, args ...string) {

}

func do_room_find(entity Entity, args ...string) {
	if len(args) == 0 {
		room := entity.GetRoom()
		entity.Send("\r\n%s\r\n", MakeTitle("Rooms", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
		rlist := make([]string, 0)
		for _, r := range room.Area.Rooms {
			if r.Name == "A void" {
				continue
			}
			n := sprintf("&Y[&W%d&Y]&d%-26s", r.Id, tstring(r.Name, 23))
			rlist = append(rlist, n)
		}
		p1 := (len(rlist) / 3) + 1
		p2 := p1 + p1
		rlist1 := rlist[:p1]
		rlist2 := rlist[p1:p2]
		rlist3 := rlist[p2:]
		pad := strings.Repeat(" ", 26)
		for i := 0; i <= p1; i++ {
			r1 := pad
			r2 := pad
			r3 := pad
			if i < len(rlist1) {
				r1 = sprintf("%-26s", rlist1[i])
			}
			if i < len(rlist2) {
				r2 = sprintf("%-26s", rlist2[i])
			}
			if i < len(rlist3) {
				r3 = sprintf("%-26s", rlist3[i])
			}
			entity.Send(sprintf("%-26s %-26s %-26s\r\n", r1, r2, r3))
		}
	}
}

func do_mob_create(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if len(args) < 2 {
		entity.Send("\r\nSyntax: mcreate <filename> <mobname>\r\n")
		return
	}
	room := DB().GetRoom(entity.RoomId(), entity.ShipId())
	if room.ship > 0 {
		entity.Send("\r\n&RMobs's must be created on planets. So we can use the area name in the filepath.&d\r\n&xsorry...&d\r\n")
		return
	}
	filename := args[0]
	// no need to use .yml in the filename, the path of the item will be /data/items/<area>/<file_name>.yml
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		filename = strings.TrimSuffix(filename, ".yaml")
		filename = strings.TrimSuffix(filename, ".yml")
	}
	mob := new(CharData)
	mob.Id = DB().GetNextMobVnum()
	mob.OId = mob.Id
	mob.Name = strings.TrimSpace(strings.Join(args[1:], " "))
	mob.Desc = "An unfinished creature stands here staring blankly off into the distance."
	mob.Brain = "generic"
	mob.Gender = "n"
	mob.Hp = []int{10, 10}
	mob.Mp = []int{0, 0}
	mob.Mv = []int{10, 10}
	mob.Skills = make(map[string]int)
	mob.Languages = make(map[string]int)
	mob.Languages["basic"] = 100
	mob.Speaking = "basic"
	mob.Race = "Human"
	mob.Equipment = make(map[string]*ItemData)
	mob.Inventory = make([]*ItemData, 0)
	mob.Stats = []int{5, 5, 5, 5, 5, 5}
	mob.Flags = make([]string, 0)
	mob.Flags = append(mob.Flags, "npc")
	mob.Filename = sprintf("data/mobs/%s/%s.yml", strings.ToLower(strings.ReplaceAll(room.Area.Name, " ", "")), filename)
	mob.State = ENTITY_STATE_NORMAL
	mob.Gold = 0
	mob.Bank = 0
	mob.Keywords = make([]string, 0)
	words := strings.Split(strings.ToLower(mob.Name), " ")
	for i, s := range words {
		if len(s) < 3 {
			words[i] = ""
			// simple words like a, or an, or the aren't good keywords.
			// You can use the itemtype as a keyword to add a proper one after creation.
		}
	}
	keywords := strings.Split(strings.Join(words, " "), " ")
	mob.Keywords = append(mob.Keywords, keywords...)
	mob.Level = 1
	mob.XP = 0
	mob.Progs = make(map[string]string)
	mob.Title = sprintf("a %s %s", mob.Race, mob.Name)
	mob.Room = room.Id
	mob.Ship = room.ship

	DB().SaveMob(mob)
	DB().LoadMob(mob.Filename)
	DB().SpawnEntity(mob)
	entity.Send("\r\n&YMob Create. Ok.&d\r\n")
}

func do_mob_spawn(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\nSyntax: mspawn <mob_vnum>\r\n")
		return
	}
	vnum, err := strconv.Atoi(args[0])
	if err != nil {
		entity.Send("\r\n&RUnable to parse argument as number.&d\r\n")
		return
	}
	mob := DB().mobs[uint(vnum)]
	if mob == nil {
		entity.Send("\r\n&RUnable to find mob with vnum: &W%d&d\r\n", vnum)
		return
	} else {
		room := entity.GetRoom()
		room.Area.Mobs = append(room.Area.Mobs, MobSpawn{
			Mob:  mob.OId,
			Room: room.Id,
		})
	}
	DB().SpawnEntity(mob)
	entity.Send("\r\n&YMob Spawn. Ok.&d\r\n")
}
func do_mob_set(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if len(args) < 2 {
		entity.Send("\r\nSyntax: mset <mob> <field> <value>\r\n")
		return
	}
	room := entity.GetRoom()
	var target Entity
	found := false
	for _, e := range room.GetEntities() {
		if found {
			break
		}
		for _, k := range e.GetCharData().Keywords {
			if strings.HasPrefix(k, strings.ToLower(args[0])) {
				target = e
				found = true
				break
			}
		}
	}
	if target == nil {
		entity.Send("\r\n&RUnable to find mob &W%s&R here.&d\r\n", args[0])
		return
	}
	tch := target.GetCharData()
	switch strings.ToLower(args[1]) {
	case "name":
		tch.Name = strings.TrimSpace(strings.Join(args[2:], " "))
	case "desc":
		tch.Desc = consolify(strings.TrimSpace(strings.Join(args[2:], " ")))
	case "race":
		found := false
		for _, r := range race_list {
			if strings.EqualFold(r, args[2]) {
				found = true
				tch.Race = r
			}
		}
		if !found {
			entity.Send("\r\n&RInvalid race.&d\r\n")
			return
		}
	case "keywords":
		switch args[2] {
		case "add":
			tch.Keywords = append(tch.Keywords, args[3])
		case "rm":
			ki := -1
			for i, k := range tch.Keywords {
				if k == args[3] {
					ki = i
				}
			}
			if ki > -1 {
				ret := make([]string, 0)
				ret = append(ret, tch.Keywords[:ki]...)
				ret = append(ret, tch.Keywords[ki+1:]...)
				tch.Keywords = ret
			}
		default:
			entity.Send("\r\nSyntax mset <mob> keywords <add/rm> <keyword>\r\n")
			return
		}
	case "gender":
		g := strings.ToLower(args[2][0:1])
		if g != "m" && g != "n" && g != "f" {
			entity.Send("\r\n&RInvalid gender.&d\r\n")
			return
		}
		tch.Gender = g
	case "level":
		lvl, _ := strconv.Atoi(args[2])
		if lvl <= 0 && lvl >= 100 {
			entity.Send("\r\n&RInvalid level. 1-99.&d\r\n")
			return
		}
		tch.Level = uint(lvl)
		tch.XP = get_xp_for_level(tch.Level)
	case "xp":
		xp, _ := strconv.Atoi(args[2])
		if xp <= 0 {
			entity.Send("\r\n&RInvalid xp. Must be positive number.&d\r\n")
			return
		}
		tch.XP = uint(xp)
		tch.Level = get_level_for_xp(tch.XP)
	case "money":
		gp, _ := strconv.Atoi(args[2])
		if gp <= 0 {
			entity.Send("\r\n&RInvalid money. Must be positive number.&d\r\n")
			return
		}
		tch.Gold = uint(gp)
	case "str":
		value, _ := strconv.Atoi(args[2])
		tch.Stats[0] = value
	case "int":
		value, _ := strconv.Atoi(args[2])
		tch.Stats[1] = value
	case "dex":
		value, _ := strconv.Atoi(args[2])
		tch.Stats[2] = value
	case "wis":
		value, _ := strconv.Atoi(args[2])
		tch.Stats[3] = value
	case "con":
		value, _ := strconv.Atoi(args[2])
		tch.Stats[4] = value
	case "cha":
		value, _ := strconv.Atoi(args[2])
		tch.Stats[5] = value
	case "hp":
		value, _ := strconv.Atoi(args[2])
		tch.Hp[0] = value
		tch.Hp[1] = value
	case "mp":
		value, _ := strconv.Atoi(args[2])
		tch.Mp[0] = value
		tch.Mp[1] = value
	case "mv":
		value, _ := strconv.Atoi(args[2])
		tch.Mv[0] = value
		tch.Mv[1] = value
	case "skill":
		skill := args[2]
		found := false
		for _, s := range skill_list {
			if s == skill {
				found = true
			}
		}
		if !found {
			entity.Send("\r\n&RInvalid skill.&d\r\n")
			return
		}
		if len(args) != 4 {
			entity.Send("\r\nSyntax: mset <mob> skill <skillname> <value>")
			return
		}
		value, _ := strconv.Atoi(args[3])
		tch.Skills[skill] = value
	case "languages":
		language := args[2]
		found := false
		for _, l := range Languages {
			if strings.EqualFold(l.Name, language) {
				found = true
			}
		}
		if !found {
			entity.Send("\r\n&RInvalid language.&d\r\n")
			return
		}
		if len(args) != 4 {
			entity.Send("\r\nSyntax: mset <mob> languages <language> <value>")
			return
		}
		value, _ := strconv.Atoi(args[3])
		tch.Languages[language] = value
	case "speaking":
		language := args[2]
		found := false
		for _, l := range Languages {
			if strings.EqualFold(l.Name, language) {
				found = true
			}
		}
		if !found {
			entity.Send("\r\n&RInvalid language.&d\r\n")
			return
		}
		tch.Speaking = language
	case "brain":
		tch.Brain = args[2]
		if tch.Brain == "generic" {
			tch.AI = MakeGenericBrain(tch)
		}
	case "flags":
		flag := args[2]
		found := false
		for i, f := range tch.Flags {
			if f == flag {
				ret := make([]string, 0)
				ret = append(ret, tch.Flags[:i]...)
				ret = append(ret, tch.Flags[i+1:]...)
				tch.Flags = ret
				found = true
			}
		}
		if !found {
			tch.Flags = append(tch.Flags, flag)
		}
	default:
		entity.Send("\r\nField values are:\r\n")
		entity.Send("-----------------------------------------\r\n")
		entity.Send("name, desc, keywords, race, gender, level, xp, money,\r\n")
		entity.Send("str, int, dex, wis, cha, con, hp, mp, mv, skill,\r\n")
		entity.Send("languages, speaking, brain\r\n")
	}
	tch.Id = tch.OId // make it an original mob. Not a clone.
	DB().SaveMob(tch)
	DB().LoadMob(tch.Filename)
	entity.Send("\r\n&YMob Set. Ok.&d\r\n")
}

func do_mob_remove(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if len(args) < 1 {
		entity.Send("\r\nSyntax: mremove <mob>\r\n")
		return
	}
	room := entity.GetRoom()
	var target Entity
	found := false
	for _, e := range room.GetEntities() {
		if found {
			break
		}
		for _, k := range e.GetCharData().Keywords {
			if strings.HasPrefix(k, strings.ToLower(args[0])) {
				target = e
				found = true
				break
			}
		}
	}
	if target == nil {
		// let's try by vnum from the mobs list (maybe we want to remove a mob entirely!)
		vnum, err := strconv.Atoi(args[0])
		if err != nil {
			entity.Send("\r\n&RUnable to find mob &W%s&R here.&d\r\n", args[0])
			return
		}
		for _, m := range DB().mobs {
			if m.OId == uint(vnum) || m.Id == uint(vnum) {
				target = m
			}
		}
	}
	if target == nil {
		entity.Send("\r\n&RUnable to find mob &W%s&R here.&d\r\n", args[0])
		return
	}
	tch := target.GetCharData()
	DB().RemoveEntity(target, false)
	delete(DB().mobs, tch.OId)
	err := os.Remove(tch.Filename)
	ErrorCheck(err)
	for _, a := range DB().areas {
		for i, msp := range a.Mobs {
			if msp.entity == target || msp.Mob == tch.OId {
				// remove the mobspawn that has this mob listed
				ret := make([]MobSpawn, 0)
				ret = append(ret, a.Mobs[:i]...)
				ret = append(ret, a.Mobs[i+1:]...)
				a.Mobs = ret
			}
		}
	}
	entity.Send("\r\n&YMob Delete. Ok.&d\r\n")
}

func do_mob_stat(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\nSyntax: mstat <mob>\r\n")
		return
	}
	is_vnum := true
	vnum, e := strconv.Atoi(args[0])
	if e != nil {
		is_vnum = false
	}
	keyword := args[0]
	room := entity.GetRoom()
	var target Entity
	if is_vnum {
		if _, ok := DB().mobs[uint(vnum)]; ok {
			target = DB().mobs[uint(vnum)]
		} else {
			entity.Send("\r\n&RUnable to find mob with vnum &W%d&d\r\n", vnum)
			return
		}
	} else {
		for _, e := range room.GetEntities() {
			for _, k := range e.GetCharData().Keywords {
				if strings.HasPrefix(k, keyword) {
					target = e
				}
			}
		}
	}
	if target == nil {
		entity.Send("\r\n&RUnable to find mob with keyword &W%s&d\r\n", keyword)
		return
	}
	tch := target.GetCharData()
	entity.Send("\r\n%s\r\n", MakeTitle("Mob Stats", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
	entity.Send("&GFilename: &W%s&d\r\n", tch.Filename)
	entity.Send("&GID: &W%d &GOID: &W%d&d\r\n", tch.Id, tch.OId)
	entity.Send("&GName: &W%-26s &GLevel: &W%d&d\r\n", tch.Name, tch.Level)
	entity.Send("&GTitle: &W%-26s &GXP: &W%d&d\r\n", tch.Title, tch.XP)
	entity.Send("&GRace: &W%s &GGender: &W%s&d\r\n", tch.Race, capitalize(get_gender_for_code(tch.Gender)))
	entity.Send("&GHp: &W%d&Y/&W%d &GMp: &W%d&Y/&W%d &GMv: &W%d&Y/&W%d&d\r\n", tch.Hp[0], tch.Hp[1], tch.Mp[0], tch.Mp[1], tch.Mv[0], tch.Mv[1])
	entity.Send("&GSTR: &W%d &GINT: &W%d &GDEX: &W%d &GWIS: &W%d &GCON: &W%d &GCHA: &W%d\r\n", tch.Stats[0], tch.Stats[1], tch.Stats[2], tch.Stats[3], tch.Stats[4], tch.Stats[5])
	entity.Send("&GMoney: &W%d &GKeywords: &W%v&d\r\n", tch.Gold, tch.Keywords)
	entity.Send("&GFlags: &W%v &GState: &W%s &GBrain: &W%s&d\r\n", tch.Flags, tch.State, tch.Brain)
	entity.Send("&GSkills: &W%+v&d\r\n", tch.Skills)
	entity.Send("&GEquipment: &W%+v&d\r\n", tch.Equipment)
	entity.Send("&GInventory: &W%v&d\r\n", tch.Inventory)
	entity.Send("&GLanguages: &W%+v&d\r\n", tch.Languages)
	entity.Send("&GSpeaking: &W%s&d\r\n", tch.Speaking)
	entity.Send("&GPrograms:&d\r\n")
	player := entity.(*PlayerProfile)
	for evt, prog := range tch.Progs {

		entity.Send("&c%s&c: |\r\n&W", evt)
		player.Client.Raw([]byte(prog)) // used to bypass colorization and any formatting.
		entity.Send("&d\r\n")
	}

	entity.Send("\r\n")
}

func do_mob_find(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\n%s\r\n", MakeTitle("Mobs", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
		c := 0
		for _, mob := range DB().mobs {
			r := mob.GetCharData()
			n := sprintf("&Y[&W%d&Y]&d%-26s", r.Id, tstring(r.Name, 23))
			entity.Send("%-30s", n)
			c++
			if c > 0 && c%6 == 0 {
				entity.Send("\r\n")
			}
		}
	}
	entity.Send("\r\n")
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
			for i := p.GetCharData().Level; i < uint(l); i++ {
				entity_advance_level(p)
			}
		}
	}
	entity.Send("\r\n&YAdvance. Ok.&d\r\n")
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
			db.SetRoom(next_id, next_room)
			for i, r := range room.Area.Rooms {
				if r.Id == next_room.Id {
					room.Area.Rooms[i] = *next_room
				}
			}
		}
		entity.Send("\r\n&GDug a room to the %s&d\r\n", dir)
	}
}

func do_ship_create(entity Entity, args ...string) {
	if len(args) < 1 {
		entity.Send("\r\nSyntax: screate <ship_type> <name>\r\n")
		return
	}
	room := entity.GetRoom()
	if !room_is_landable(room) {
		entity.Send("\r\n&RShips can only be created in spaceports, shipyards, and hangars.&d\r\n")
		return
	}
	ship_type := args[1]
	id := DB().GetNextShipVnum()
	ship := &ShipData{
		Id:            id,
		OId:           id,
		Name:          strings.Join(args[2:], " "),
		Desc:          "A prototype ship",
		Type:          ship_type,
		LocationId:    entity.GetRoom().Id,
		CurrentSystem: "Somewhere",
		ShipyardId:    entity.GetRoom().Id,
		Rooms:         make(map[uint]*RoomData),
		Owner:         entity.GetCharData().Name,
		Modules:       make(map[string]uint),
		Position:      []float32{0.0, 0.0},
		HighSlots:     make([]*ItemData, 0),
		LowSlots:      make([]*ItemData, 0),
		Cockpit:       1,
		Ramp:          1,
		EngineRoom:    1,
		CargoRoom:     1,
		Blueprint:     0,
		MaxSpeed:      1,
		Hp:            []uint{100, 100},
		Sp:            []uint{0, 0},
	}
	ship.Rooms[1] = &RoomData{
		Id:   1,
		Name: "A prototype cockpit",
		Desc: "A stripped down prototype cockpit, barely able to maintain flight.",
		ship: ship.Id,
	}
	DB().SaveShip(ship)
	DB().LoadShip(sprintf("data/ships/%s.yml", ship.Name))
	DB().SpawnShip(ship)

	entity.Send("\r\n&YShip Create. Ok.&d\r\n")

}

func do_ship_remove(entity Entity, args ...string) {
	if len(args) < 1 {
		entity.Send("\r\nSyntax: sremove <shipname> [prototype?]\r\n")
		entity.Send("-----------------------------------------------------------------\r\n")
		entity.Send("prototype is a bool, 1 or 0 and is optional. If supplied, and is\r\n")
		entity.Send("1, then the ship's prototype data will be removed as well.\r\n")
		return
	}
	room := entity.GetRoom()
	l := len(args)
	prototype := false
	if args[l-1] == "1" {
		prototype = true
	}
	ship_name := strings.Join(args[:l], " ")
	if prototype {
		ship_name = strings.Join(args[:l-2], " ")
	}
	var ship Ship
	for _, s := range room.GetShips() {
		if s.GetData().Name == ship_name {
			ship = s
		}
	}
	if ship == nil {
		entity.Send("\r\n&RUnable to locate ship!&d\r\n")
		return
	}
	DB().RemoveShip(ship)
	e := os.Remove(sprintf("data/ships/%s.yml", strings.ToLower(strings.ReplaceAll(ship.GetData().Name, " ", "_"))))
	ErrorCheck(e)
	if prototype {
		DB().RemoveShipPrototype(ship)
		e = os.Remove(sprintf("data/ships/prototypes/%s.yml", strings.ToLower(strings.ReplaceAll(ship.GetData().Name, " ", "_"))))
		ErrorCheck(e)
	}
	entity.Send("\r\n&YRemove Ship. Ok.&d\r\n")
}

func do_ship_set(entity Entity, args ...string) {

}

func do_ship_stat(entity Entity, args ...string) {

}
