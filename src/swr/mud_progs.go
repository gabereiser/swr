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

type Brain interface {
	OnSpawn()
	OnEnter(entity Entity)
	OnGreet(entity Entity)
	OnMove()
	OnKill(entity Entity)
	OnDeath()
	OnDrop(entity Entity, item Item)
	OnHeal(entity Entity, item Item)
	OnGive(entity Entity, item Item)
	OnSay(entity Entity, words string)
	Update()
}

type GenericBrain struct {
	Entity Entity
}

func MakeGenericBrain(entity Entity) *GenericBrain {
	brain := new(GenericBrain)
	brain.Entity = entity
	return brain
}
