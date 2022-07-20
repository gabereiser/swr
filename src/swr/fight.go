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
	"strings"
)

func do_kill(entity Entity, args ...string) {
	do_fight(entity, args...)
}

func do_fight(entity Entity, args ...string) {
	if len(args) < 1 {
		entity.Send("\r\n&RFight who?&d\r\n")
		return
	}
	db := DB()
	if entity.IsFighting() {
		entity.Send("\r\n&RYou are already fighting!&d\r\n")
		return
	} else {
		state := entity.GetCharData().State
		if state == ENTITY_STATE_CRAFTING {
			entity.Send("\r\n&RYou can't fight while working!&d\r\n")
			return
		}
		if state == ENTITY_STATE_DEAD {
			entity.Send("\r\n&RYou are dead!&d\r\n")
			return
		}
		if state == ENTITY_STATE_GUNNING {
			entity.Send("\r\n&RYou can't gun and fight at the same time!&d\r\n")
			return
		}
		if state == ENTITY_STATE_PILOTING {
			entity.Send("\r\n&RYou can't fly and fight at the same time!&d\r\n")
			return
		}
		if state == ENTITY_STATE_SEDATED {
			entity.Send("\r\nYou feel too relaxed!\r\n")
			return
		}
		if state == ENTITY_STATE_SLEEPING {
			entity.Send("\r\nYou are asleep!\r\n")
			return
		}
		if state == ENTITY_STATE_UNCONSCIOUS {
			entity.Send("\r\n&RYou are unconscious!&d\r\n")
			return
		}
		found := false
		for _, e := range db.GetEntitiesInRoom(entity.RoomId()) {
			if e == entity {
				continue
			}
			ch := e.GetCharData()
			for _, k := range ch.Keywords {
				if strings.HasPrefix(k, args[0]) {
					e.SetAttacker(entity)
					entity.SetAttacker(e)
					entity.Send("\r\n&RYou begin fighting &w%s&R!!&d\r\n", ch.Name)
					found = true
					do_combat(entity, e)
				}
			}
		}
		if !found {
			entity.Send("\r\n&dThey aren't here.\r\n")
		}
	}
}

func processCombat() {
	db := DB()
	db.Lock()
	defer db.Unlock()
	for _, e := range db.entities {
		if e.IsFighting() {
			target := e.GetCharData().Attacker
			do_combat(e, target)
		}
	}
}

func do_combat(attacker Entity, defender Entity) {

	ach := attacker.GetCharData()
	dch := defender.GetCharData()

	hit_chance := roll_dice("1d20")
	damage := uint(0)
	if hit_chance > dch.ArmorAC() {
		damage = ach.DamageRoll()
		dch.ApplyDamage(damage)

	}
	attacker.Send(get_damage_string(damage, "You", dch.Name, "an object."))
	defender.Send(get_damage_string(damage, dch.Name, "you", "an object."))

	make_corpse(attacker)
	make_corpse(defender)
}

func get_damage_string(damage uint, attacker string, defender string, weapon string) string {
	if damage > 20 {
		return fmt.Sprintf("\r\n&R%s ANNIHILATES %s with %s.&d\r\n", attacker, defender, weapon)
	} else if damage > 15 {
		return fmt.Sprintf("\r\n&R%s EVICERATES %s with %s.&d\r\n", attacker, defender, weapon)
	} else if damage > 10 {
		return fmt.Sprintf("\r\n&R%s BRUISES %s with %s.&d\r\n", attacker, defender, weapon)
	} else if damage > 5 {
		return fmt.Sprintf("\r\n&R%s HITS %s with %s.&d\r\n", attacker, defender, weapon)
	} else if damage > 1 {
		return fmt.Sprintf("\r\n&R%s NICKS %s with %s.&d\r\n", attacker, defender, weapon)
	} else {
		return fmt.Sprintf("\r\n&R%s MISSES %s.&d\r\n", attacker, defender)
	}
}

func make_corpse(entity Entity) {
	ch := entity.GetCharData()
	if entity.CurrentHp() < entity.MaxHp()*2 {
		corpse := &ItemData{
			Id:     gen_item_id(),
			Name:   fmt.Sprintf("A corpse of a %s %s", ch.Race, ch.Gender),
			Desc:   fmt.Sprintf("A bloody corpse of a %s %s lies here in rot.", ch.Race, ch.Gender),
			Type:   ITEM_TYPE_CORPSE,
			Value:  int(ch.Gold),
			Weight: ch.CurrentWeight(),
			AC:     0,
			Items:  make([]Item, 0),
		}
		items := make([]Item, 0)
		for _, item := range ch.Equipment {
			if !entity.IsPlayer() && roll_dice("1d4") == 4 {
				items = append(items, item_clone(item))
			} else if entity.IsPlayer() {
				items = append(items, item_clone(item))
			}
		}
		for i := range ch.Inventory {
			item := ch.Inventory[i]
			items = append(items, item_clone(item))
		}
		corpse.Items = items
		room := DB().GetRoom(ch.Room)
		room.AddItem(corpse)
		if entity.IsPlayer() {
			entity.Send("\r\n&RYou have been killed.")
			respawn_entity(entity, 100)
			entity.Prompt()
		} else {
			DB().RemoveEntity(entity)
		}

	}
}

// Respawn an entity in a room.
func respawn_entity(entity Entity, roomId uint) {
	ch := entity.GetCharData()
	ch.Room = roomId
	ch.Hp[0] = ch.Hp[1]
	ch.State = ENTITY_STATE_NORMAL
}
