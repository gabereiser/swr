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
		pch := entity.GetCharData()
		for _, e := range db.GetEntitiesInRoom(entity.RoomId()) {
			ch := e.GetCharData()
			for _, k := range ch.Keywords {
				if strings.HasPrefix(k, args[0]) {
					e.SetAttacker(&entity)
					entity.SetAttacker(&e)
					entity.Send("\r\n&RYou begin fighting &w%s&R!!&d\r\n", e.Name())
					do_combat(entity, e)
				}
			}
		}
	}
}

func process_combat() {
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
	attacker.Send(get_damage_string(damage, "You", dch.CharName, "an object."))
	defender.Send(get_damage_string(damage, dch.CharName, "you", "an object."))
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
