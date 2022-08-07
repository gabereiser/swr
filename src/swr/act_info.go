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
	"strings"
	"time"
)

func do_quit(entity Entity, args ...string) {
	if entity.IsPlayer() {
		if entity.IsFighting() {
			entity.Send("\r\n&RYou can't quit while fighting!&d\r\n")
			return
		}
		player := entity.(*PlayerProfile)
		DB().SavePlayerData(player)
		entity.Send("\r\n&CThe world slowly fades away as you close your eyes and leave the game...&d\r\n\r\n")
		entity.GetCharData().State = ENTITY_STATE_SLEEPING
		ScheduleFunc(func() {
			player.Send("\r\n%s Thank you for playing! %s\r\n", EMOJI_ALERT, EMOJI_ALERT)
			time.Sleep(100 * time.Millisecond)
			player.Client.Close()
		}, false, 1)
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
			entity.Send(sprintf("&W%-67s&G [ &WLevel %2d&G ]\r\n", player.Char.Title, player.Char.Level))
			total++
		}
	}
	entity.Send("\r\n")
	entity.Send(MakeTitle(sprintf("%d Online", total), ANSI_TITLE_STYLE_NORMAL, ANSI_TITLE_ALIGNMENT_RIGHT))
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
		player.Send("&c│ Weight: &G%3d kg&p(%4d kg)&c                  │&d▒\r\n", char.CurrentWeight(), char.MaxWeight())
		player.Send("&c│ Inventory: &G%3d&p(%3d)&c                      │&d▒\r\n", char.CurrentInventoryCount(), char.MaxInventoryCount())
		player.Send("&c│ Kills: &G%-5d    &cPlayer Kills: &G%-5d&c      │&d▒\r\n", player.Kills, player.PKills)
		player.Send("&c├─( Equipment )────────────────────────────┤▒&d\r\n")
		player.Send("&c│       Head: &d%-20s&c         │&d▒\r\n", entity_get_equipment_for_slot(player, "head"))
		player.Send("&c│      Torso: &d%-20s&c         │&d▒\r\n", entity_get_equipment_for_slot(player, "torso"))
		player.Send("&c│      Waist: &d%-20s&c         │&d▒\r\n", entity_get_equipment_for_slot(player, "waist"))
		player.Send("&c│       Legs: &d%-20s&c         │&d▒\r\n", entity_get_equipment_for_slot(player, "legs"))
		player.Send("&c│       Feet: &d%-20s&c         │&d▒\r\n", entity_get_equipment_for_slot(player, "feet"))
		player.Send("&c│      Hands: &d%-20s&c         │&d▒\r\n", entity_get_equipment_for_slot(player, "hands"))
		player.Send("&c│                                          │&d▒\r\n")
		player.Send("&c│     &RWeapon: &d%-20s&c         │&d▒\r\n", entity_get_equipment_for_slot(player, "weapon"))
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
		if item == nil {
			continue
		}
		player.Send("&c│ %-34s│&d▒\r\n", item.GetData().Name)
	}
	player.Send("&c└───────────────────────────────────┘&d▒\r\n")
	player.Send(" ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒&d\r\n")
}

func do_levels(entity Entity, args ...string) {
	ch := entity.GetCharData()
	entity.Send("\r\n%s\r\n", MakeTitle("Levels / Experience", ANSI_TITLE_STYLE_NORMAL, ANSI_TITLE_ALIGNMENT_LEFT))
	entity.Send("&YLevel: &W%3d&Y Exp: &W%d&Y/&W%d&d\r\n\r\n", ch.Level, ch.XP, get_xp_for_level(ch.Level))
	level := ch.Level
	if level > 100-5 {
		level -= 5
	}
	m := " " // used to mark our level in the list.
	for i := 1; i < 6; i++ {
		if uint(i) == ch.Level {
			m = "}b>&d&Y" // pretty >...
		}
		entity.Send("%s&YLevel: &W%3d&Y Exp: &W%d&d\r\n", m, level+uint(i), get_xp_for_level(level+uint(i)))
	}
}

func do_description(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\n&CSyntax: description <string>.&d\r\n")
		return
	}
	ch := entity.GetCharData()
	ch.Desc = consolify(args[0])
	entity.Send("\r\n&YDescription set.&d\r\n")
}

func do_examine(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\n&RExamine what?&d\r\n")
		return
	}
	object_name := args[0]
	room := DB().GetRoom(entity.RoomId(), entity.ShipId())
	object := entity.FindItem(object_name)
	if object == nil {
		object = room.FindItem(object_name)
	}
	if object == nil {
		for _, e := range room.GetEntities() {
			if strings.HasPrefix(e.GetCharData().Name, object_name) {
				entity.Send("&dYou look at &W%s&d and see...\r\n%s\r\n", e.GetCharData().Name, e.GetCharData().Desc)
				entity.Send("&YEquipment:\r\n-------------------------------------&d\r\n")
				if len(e.GetCharData().Equipment) == 0 {
					entity.Send("Nothing\r\n")
				} else {
					entity.Send("&YHead: &d%-26s\r\n", entity_get_equipment_for_slot(entity, "head"))
					entity.Send("&YTorso: &d%-26s\r\n", entity_get_equipment_for_slot(entity, "torso"))
					entity.Send("&YWaist: &d%-26s\r\n", entity_get_equipment_for_slot(entity, "waist"))
					entity.Send("&YLegs: &d%-26s\r\n", entity_get_equipment_for_slot(entity, "legs"))
					entity.Send("&YFeet: &d%-26s\r\n", entity_get_equipment_for_slot(entity, "feet"))
					entity.Send("&YHands: &d%-26s\r\n", entity_get_equipment_for_slot(entity, "hands"))
					entity.Send("&Y--------------------------------------&d\r\n")
					entity.Send("&RWeapon: &d%-26s\r\n", entity_get_equipment_for_slot(entity, "weapon"))
				}
				return
			}
		}
		entity.Send("\r\nCan't find that here.\r\n")
		return
	} else {
		entity.Send("You look at %s and see...\r\n%s\r\n", object.GetData().Name, object.GetData().Desc)
		if object.IsContainer() {
			entity.Send("&YContents:\r\n-------------------------------------&d\r\n")
			for _, o := range object.GetData().Items {
				if o == nil {
					continue
				}
				entity.Send("&Y%-26s&d\r\n", o.GetData().Name)
			}
		}
		return
	}

}

func do_equip(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\n&REquip what?&d\r\n")
		return
	}
	item_name := args[0]
	item := entity.FindItem(item_name)
	if item == nil {
		entity.Send("\r\n&RYou don't have that item!&d\r\n")
		return
	}
	if !item.IsWeapon() && !item.IsWearable() {
		entity.Send("\r\n&RYou can't equip that item!&d\r\n")
		return
	}
	data := item.GetData()

	if data.WearLoc == nil && !item.IsWeapon() {
		entity.Send("\r\n&RBUG: Unable to determine wear location!&d\r\n")
		log.Printf("BUG: Unable to determine wear location for %s\r\n", data.Name)
		return
	}
	wearLoc := "weapon"
	if !item.IsWeapon() {
		wearLoc = *data.WearLoc
	}
	if !item_is_wearable_slot(wearLoc) && wearLoc != "weapon" {
		entity.Send("\r\n&RBUG: Item is wearable but on an invalid slot.&d\r\n")
		log.Printf("BUG: Item OId[%d] is wearable but on an invalid slot.&d", item.GetData().OId)
		return
	}
	entity.GetCharData().Equipment[wearLoc] = data
	entity.Send("\r\n&YYou equip %s&d\r\n", data.Name)
	entity.GetCharData().RemoveItem(item)
	others := DB().GetEntitiesInRoom(entity.RoomId(), entity.ShipId())
	for _, e := range others {
		if e == nil {
			continue
		}
		if e != entity {
			e.Send("%s equips %s&d\r\n", entity.GetCharData().Name, data.Name)
		}
	}
}

func do_remove(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\n&RRemove what?&d\r\n")
		return
	}
	ch := entity.GetCharData()
	var item Item = nil
	for wearLoc, i := range ch.Equipment {
		if strings.HasPrefix(wearLoc, args[0]) {
			item = i
			break
		}
		for _, keyword := range i.Keywords {
			if strings.HasPrefix(keyword, args[0]) {
				item = i
				break
			}
		}
	}
	if item != nil {
		if !entity_pickup_item(entity, item) {
			entity.Send("\r\n&RYou can't carry anymore!&d\r\n")
			return
		}
		delete(ch.Equipment, *item.GetData().WearLoc)
	}
	data := item.GetData()
	others := DB().GetEntitiesInRoom(entity.RoomId(), entity.ShipId())
	for _, e := range others {
		if e == nil {
			continue
		}
		if e != entity {
			e.Send("%s removes %s&d\r\n", entity.GetCharData().Name, data.Name)
		}
	}
	entity.Send("\r\n&YYou remove %s&d\r\n", data.Name)
}

func do_time(entity Entity, args ...string) {
	now := time.Now()
	entity.Send("\r\n&BHolonet Time Synchronization&d\r\n")
	entity.Send("&g----------------------------------------------------------------&d\r\n")
	entity.Send("&cThe Current Local Time is: &Y%s&d\r\n", now.Format(time.RFC822))
	entity.Send("&cThe Current UTC Time is: &Y%s&d\r\n", now.UTC().Format(time.RFC822Z))
	entity.Send("&cThe Server Started at: &Y%s&d\r\n", startup.Format(time.RFC822))
	entity.Send("&CThe Server has been running for &Y%s&d\r\n", time.Since(startup).String())
	entity.Send("\r\n")
}
