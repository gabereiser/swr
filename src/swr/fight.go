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
				if strings.HasPrefix(strings.ToLower(k), strings.ToLower(args[0])) {
					found = true
					if ch.State != ENTITY_STATE_DEAD && ch.State != ENTITY_STATE_UNCONSCIOUS {
						e.SetAttacker(entity)
						entity.SetAttacker(e)
						entity.Send("\r\n&RYou begin fighting &w%s&R!!&d\r\n", ch.Name)
						break
					} else {
						entity.Send("\r\n&RYou can't fight what can't fight back.&d\r\n")
					}

				}
			}
			if found {
				break
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
		if e != nil {
			if e.IsFighting() {
				target := e.GetCharData().Attacker
				do_combat(e, target)
			}
		}
	}
	for _, e := range db.entities {
		if e != nil {
			e.Prompt()
		}
	}
}

func do_combat(attacker Entity, defender Entity) {

	if attacker == nil || defender == nil {
		return
	}
	ach := attacker.GetCharData()
	dch := defender.GetCharData()

	for _, flag := range ach.Flags {
		if flag == "nofight" {
			return
		}
	}
	hit_chance := roll_dice("1d20")
	damage := uint(0)
	ach_weapon := "fists"
	if attacker.Weapon() != nil {
		ach_weapon = attacker.Weapon().GetData().Name
	}
	if dch.State == ENTITY_STATE_UNCONSCIOUS && defender.IsPlayer() {
		attacker.StopFighting()
		return
	}
	if dch.Attacker == nil && dch.Mv[0] > 0 {
		dch.Attacker = attacker
		dch.State = ENTITY_STATE_FIGHTING
	}

	if hit_chance > dch.ArmorAC() || hit_chance == 20 {
		if ach.Mv[0] <= 0 {
			attacker.Send("\r\n&YYou are exhausted.&d\r\n")
			attacker.StopFighting()
			ach.State = ENTITY_STATE_NORMAL
			return
		}
		if roll_dice("1d20") <= 10-(entity_get_skill_value(ach, "aerobics")/10) {
			ach.Mv[0]--
			if roll_dice("1d20") == 20 {
				entity_add_skill_value(attacker, "aerobics", 1)
			}
		}
		skill := "martial-arts"
		if ach.Weapon() != nil {
			skill = item_get_weapon_skill(ach.Weapon())
		}
		damage = ach.DamageRoll(skill)
		defender.ApplyDamage(damage)
	}
	xp_base := 1
	attacker.Send(get_damage_string(damage, "You", dch.Name, fmt.Sprintf("your %s", ach_weapon)))
	defender.Send(get_damage_string(damage, ach.Name, "you", fmt.Sprintf("their %s", ach_weapon)))
	if attacker.GetCharData().State == ENTITY_STATE_DEAD {
		attacker.Send("\r\n&W%s &Rhas killed you.&d\r\n", dch.Name)
		defender.Send("\r\n&RYou have killed &W%s&d\r\n", ach.Name)
		attacker.StopFighting()
		defender.StopFighting()
		xp_base = int(dch.Level) * 250
		make_corpse(attacker)
		entity_add_xp(defender, xp_base)
		entity_lose_xp(attacker, xp_base)
		return
	}
	if attacker.GetCharData().State == ENTITY_STATE_UNCONSCIOUS {
		attacker.Send("\r\n&W%s &Rhas knocked you out.&d\r\n", dch.Name)
		defender.Send("\r\n&RYou have knocked out &W%s&d\r\n", ach.Name)
		attacker.StopFighting()
		defender.StopFighting()
		xp_base = int(dch.Level) * 100
		entity_add_xp(defender, xp_base)
		entity_lose_xp(attacker, xp_base)
		return
	}
	if defender.GetCharData().State == ENTITY_STATE_DEAD {
		defender.Send("\r\n&R%s has killed you.&d\r\n", ach.Name)
		attacker.Send("\r\n&RYou have killed &W%s&d\r\n", dch.Name)
		attacker.StopFighting()
		defender.StopFighting()
		xp_base = int(ach.Level) * 250
		make_corpse(defender)
		entity_lose_xp(defender, xp_base)
		entity_add_xp(attacker, xp_base)
		return
	}
	if defender.GetCharData().State == ENTITY_STATE_UNCONSCIOUS {
		defender.Send("\r\n&W%s &Rhas knocked you out.&d\r\n", ach.Name)
		attacker.Send("\r\n&RYou have knocked out &W%s&d\r\n", dch.Name)
		attacker.StopFighting()
		defender.StopFighting()
		xp_base = int(ach.Level) * 100
		entity_lose_xp(defender, xp_base)
		entity_add_xp(attacker, xp_base)
		return
	}
	if roll_dice("1d10") == 10 {
		entity_add_skill_value(attacker, item_get_weapon_skill(ach.Weapon()), 1)
	}
	xp := rand_min_max(1, 5) * int(ach.Level)
	entity_add_xp(attacker, int(dch.Level)*xp)
}

func get_damage_string(damage uint, attacker string, defender string, weapon string) string {
	if damage > 50 {
		return fmt.Sprintf("&R%s **ANNIHILATED** %s with %s for &w%d&R damage.&d\r\n", attacker, defender, weapon, damage)
	} else if damage > 25 {
		return fmt.Sprintf("&R%s *EVICERATED* %s with %s for &w%d&R damage.&d\r\n", attacker, defender, weapon, damage)
	} else if damage > 10 {
		return fmt.Sprintf("&R%s *BLASTED* %s with %s for &w%d&R damage.&d\r\n", attacker, defender, weapon, damage)
	} else if damage > 2 {
		return fmt.Sprintf("&R%s *HIT* %s with %s for &w%d&R damage.&d\r\n", attacker, defender, weapon, damage)
	} else if damage > 1 {
		return fmt.Sprintf("&R%s SCRATCHED %s with %s for &w%d&R damage.&d\r\n", attacker, defender, weapon, damage)
	} else {
		return fmt.Sprintf("&d%s MISSED %s.&d\r\n", attacker, defender)
	}
}

func make_corpse(entity Entity) {
	ch := entity.GetCharData()
	if ch.State == ENTITY_STATE_DEAD {
		corpse := &ItemData{
			Id:       gen_item_id(),
			Name:     fmt.Sprintf("a corpse of a %s %s", get_gender_for_code(ch.Gender), strings.ToLower(ch.Race)),
			Keywords: make([]string, 0),
			Desc:     fmt.Sprintf("A bloody corpse of a %s %s lies here in rot.", get_gender_for_code(ch.Gender), strings.ToLower(ch.Race)),
			Type:     ITEM_TYPE_CORPSE,
			Value:    int(ch.Gold),
			Weight:   ch.base_weight(),
			AC:       0,
			Items:    make([]Item, 0),
		}
		if !entity.IsPlayer() {
			corpse.Name = fmt.Sprintf("a corpse of %s", entity.GetCharData().Name)
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
		for k := range ch.Keywords {
			key := ch.Keywords[k]
			corpse.Keywords = append(corpse.Keywords, key)
		}
		corpse.Keywords = append(corpse.Keywords, "corpse")
		corpse.Items = items
		room := DB().GetRoom(ch.Room)
		room.AddItem(corpse)
		if entity.IsPlayer() {
			entity.Send("\r\n&RYou have been killed.&d\r\n\r\n&yRespawning in 30 seconds...&d\r\n")
			ScheduleFunc(func() {
				respawn_entity(entity, 100)
			}, false, 5)
		} else {
			DB().RemoveEntity(entity)
		}
		log.Printf("Entity %s [%d] has been killed.", ch.Name, ch.Id)
	}
}

// Respawn an entity in a room.
func respawn_entity(entity Entity, roomId uint) {
	ch := entity.GetCharData()
	ch.Attacker = nil
	ch.Room = roomId
	ch.Hp[0] = ch.Hp[1]
	ch.State = ENTITY_STATE_NORMAL
	do_look(entity)
	entity.Prompt()
}
