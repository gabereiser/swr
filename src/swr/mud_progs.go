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
	"strconv"
	"time"

	"github.com/robertkrimen/otto"
)

type Brain interface {
	OnSpawn()
	OnEnter(entity Entity)
	OnGreet(entity Entity)
	OnMove(entity Entity)
	OnKill(entity Entity)
	OnDeath()
	OnDrop(entity Entity, item Item)
	OnHeal(entity Entity)
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

func (b *GenericBrain) OnSpawn() {
	go mud_prog_exec("spawn", b.Entity)
}
func (b *GenericBrain) OnDeath() {
	go mud_prog_exec("death", b.Entity)
}
func (b *GenericBrain) OnKill(entity Entity) {
	go mud_prog_exec("kill", b.Entity, entity)
}
func (b *GenericBrain) OnMove(entity Entity) {
	go mud_prog_exec("move", b.Entity, entity)
}
func (b *GenericBrain) OnEnter(entity Entity) {
	go mud_prog_exec("enter", b.Entity, entity)
}
func (b *GenericBrain) OnGreet(entity Entity) {
	go mud_prog_exec("greet", b.Entity, entity)
}
func (b *GenericBrain) OnDrop(entity Entity, item Item) {
	go mud_prog_exec("drop", b.Entity, entity, item)
}
func (b *GenericBrain) OnGive(entity Entity, item Item) {
	go mud_prog_exec("give", b.Entity, entity, item)
}
func (b *GenericBrain) OnHeal(entity Entity) {
	go mud_prog_exec("heal", b.Entity, entity)
}
func (b *GenericBrain) OnSay(entity Entity, words string) {
	go mud_prog_exec("say", b.Entity, entity, words)
}
func (b *GenericBrain) Update() {

}

func mud_prog_exec(prog string, entity Entity, any ...interface{}) error {
	ch := entity.GetCharData()
	if pg, ok := ch.Progs[prog]; ok {
		vm := mud_prog_init(entity)
		any_len := len(any)
		if any_len > 0 {
			err := vm.Set("$n", any[0].(Entity))
			ErrorCheck(err)
			if any_len > 1 {
				if _, ok := any[1].(string); ok {
					err = vm.Set("$s", any[1].(string))
					ErrorCheck(err)
				} else {
					err = vm.Set("$i", any[1].(Item))
					ErrorCheck(err)
				}
			}
		}
		_, err := vm.Run(pg)
		ErrorCheck(err)
		return err
	}
	return Err("%s is not a program of [%d]%s", prog, ch.Id, ch.Name)
}
func mud_prog_init(entity Entity) *otto.Otto {
	vm := otto.New()
	err := vm.Set("$me", entity.GetCharData().Name)
	if err != nil {
		panic(err)
	}
	// say("hello");
	vm.Set("say", func(call otto.FunctionCall) otto.Value {
		do_say(entity, call.Argument(0).String())
		return otto.Value{}
	})
	// emote("sits down");
	vm.Set("emote", func(call otto.FunctionCall) otto.Value {
		do_emote(entity, call.Argument(0).String())
		return otto.Value{}
	})
	// transfer($n, 100);  - $n is the player, 100 is the room_id
	vm.Set("transfer", func(call otto.FunctionCall) otto.Value {
		entity_name := call.Argument(0).String()
		room_value, _ := call.Argument(1).ToInteger()
		do_transfer(entity, entity_name, strconv.Itoa(int(room_value)))
		return otto.Value{}
	})
	vm.Set("delay", func(call otto.FunctionCall) otto.Value {
		t, _ := call.Argument(0).ToInteger()
		time.Sleep(time.Duration(t) * time.Second)
		return otto.Value{}
	})
	return vm
}
