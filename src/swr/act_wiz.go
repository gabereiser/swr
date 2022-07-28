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

import "strconv"

func do_area_create(entity Entity, args ...string) {

}

func do_area_set(entity Entity, args ...string) {

}

func do_area_remove(entity Entity, args ...string) {

}

func do_area_reset(entity Entity, args ...string) {

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
	room := DB().GetRoom(target.RoomId())
	for _, e := range room.GetEntities() {
		if e == nil {
			continue
		}
		if e != target {
			e.Send("\r\n&C%s&d has left.\r\n")
		}
	}
	target.GetCharData().Room = uint(room_id)
	room = DB().GetRoom(uint(room_id))
	for _, e := range room.GetEntities() {
		if e == nil {
			continue
		}
		if e != target {
			e.Send("\r\n&C%s&d has appeared.\r\n")
		}
	}
	entity.Send("\r\n&YOk.&d\r\n")
}
