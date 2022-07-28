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
	"math/rand"
	"strconv"
	"time"

	"github.com/robertkrimen/otto"
)

type Brain interface {
	OnSpawn()
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
		mud_prog_bind(vm, any...)
		_, err := vm.Run(pg)
		ErrorCheck(err)
		return err
	}
	return Err("%s is not a program of [%d]%s", prog, ch.Id, ch.Name)
}
func mud_prog_bind(vm *otto.Otto, any ...interface{}) {
	any_len := len(any)
	if any_len > 0 {
		for i := 0; i < any_len; i++ {
			if _, ok := any[i].(Entity); ok {
				err := vm.Set("$n", any[i].(Entity).GetCharData().Name)
				ErrorCheck(err)
			} else if _, ok := any[i].(string); ok {
				err := vm.Set("$s", any[i].(string))
				ErrorCheck(err)
			} else if _, ok := any[i].(Item); ok {
				err := vm.Set("$i", any[i].(Item).GetData().Name)
				ErrorCheck(err)
			}
		}
	}
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
	vm.Set("shout", func(call otto.FunctionCall) otto.Value {
		do_shout(entity, call.Argument(0).String())
		return otto.Value{}
	})
	// emote("sits down");
	vm.Set("emote", func(call otto.FunctionCall) otto.Value {
		do_emote(entity, call.Argument(0).String())
		return otto.Value{}
	})
	vm.Set("echo", func(call otto.FunctionCall) otto.Value {
		entity.Send(call.Argument(0).String())
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
	vm.Set("random", func(call otto.FunctionCall) otto.Value {
		arg, _ := call.Argument(0).ToInteger()
		value, _ := otto.ToValue(rand.Intn(int(arg)))
		return value
	})
	vm.Set("sprintf", func(call otto.FunctionCall) otto.Value {
		format, _ := call.Argument(0).ToString()
		vm_args_list := call.ArgumentList[1:]
		var args []interface{} = make([]interface{}, 0)
		for _, arg := range vm_args_list {
			if arg.IsBoolean() {
				a, _ := arg.ToBoolean()
				args = append(args, a)
				continue
			}
			if arg.IsNumber() {
				a, _ := arg.ToInteger()
				args = append(args, a)
				continue
			}
			if arg.IsString() {
				a, _ := arg.ToString()
				args = append(args, a)
				continue
			}
		}
		value, _ := otto.ToValue(sprintf(format, args...))
		return value
	})
	vm.Set("look", func(call otto.FunctionCall) otto.Value {
		do_look(entity)
		return otto.Value{}
	})
	vm.Set("kill", func(call otto.FunctionCall) otto.Value {
		target, _ := call.Argument(0).ToString()
		do_fight(entity, target)
		return otto.Value{}
	})
	vm.Set("stand", func(call otto.FunctionCall) otto.Value {
		do_stand(entity)
		return otto.Value{}
	})
	return vm
}
