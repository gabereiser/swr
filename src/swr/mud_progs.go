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
	"strings"
	"time"

	"github.com/robertkrimen/otto"
)

type Brain interface {
	// OnSpawn event handler. Executes the brains "spawn" program. When the entity enters the game.
	OnSpawn()
	// OnGreet event handler. Executes the brains "greet" program. When another entity enters the room.
	OnGreet(entity Entity)

	// OnMove event handler. Executes the brains "move" program. Not to be confused with [GenericBrain.Move()].
	// This is triggered when another entity moves from the room.
	OnMove(entity Entity)

	// OnKill event handler. Executes the brains "kill" program.
	OnKill(entity Entity)

	// OnDeath event handler. Executes the brains "death" program.
	OnDeath()

	// OnDrop event handler. Executes the brains "drop" program.
	OnDrop(entity Entity, item Item)

	// OnHeal event handler. Executes the brains "heal" program.
	OnHeal(entity Entity)

	// OnGive event handler. Executes the brains "give" program.
	OnGive(entity Entity, quantity int, item Item)

	// OnSay event handler. Executes the brains "say" program.
	OnSay(entity Entity, words string)

	// Update is the main logic tree for AI and [GenericBrain]. It will figure out which action to take on the controlling entity.
	Update()
}

type GenericBrain struct {
	Entity Entity
	vm     *otto.Otto
}

// MakeGenericBrain creates a [*GenericBrain] instance and wraps the entity in it. Effectively passing control to the brain.
func MakeGenericBrain(entity Entity) *GenericBrain {
	brain := new(GenericBrain)
	brain.Entity = entity
	brain.vm = mud_prog_init(entity)
	return brain
}

func (b *GenericBrain) OnSpawn() {
	go mud_prog_exec(b.vm, "spawn", b.Entity)
}
func (b *GenericBrain) OnDeath() {
	go mud_prog_exec(b.vm, "death", b.Entity)
}
func (b *GenericBrain) OnKill(entity Entity) {
	go mud_prog_exec(b.vm, "kill", b.Entity, entity)
}
func (b *GenericBrain) OnMove(entity Entity) {
	go mud_prog_exec(b.vm, "move", b.Entity, entity)
}
func (b *GenericBrain) OnGreet(entity Entity) {
	go mud_prog_exec(b.vm, "greet", b.Entity, entity)
}
func (b *GenericBrain) OnDrop(entity Entity, item Item) {
	go mud_prog_exec(b.vm, "drop", b.Entity, entity, item)
}
func (b *GenericBrain) OnGive(entity Entity, quantity int, item Item) {
	go mud_prog_exec(b.vm, "give", b.Entity, entity, quantity, item)
}
func (b *GenericBrain) OnHeal(entity Entity) {
	go mud_prog_exec(b.vm, "heal", b.Entity, entity)
}
func (b *GenericBrain) OnSay(entity Entity, words string) {
	go mud_prog_exec(b.vm, "say", b.Entity, entity, words)
}

/* Update is called every server tick, it's the main logic tree for AI and {GenericBrain}
 */
func (b *GenericBrain) Update() {
	if b.Entity.GetCharData().State == ENTITY_STATE_NORMAL {
		move := true
		for _, f := range b.Entity.GetCharData().Flags {
			if strings.ToLower(f) == "sentinel" {
				move = false
			}
		}
		if roll_dice("1d30") == 30 && move {
			// let's try to move...
			b.Move()
		}
	}
}

// Move makes the brain perform a move action.
// It will move it's controlling entity to another room using do_direction,
// same as a player.
func (b *GenericBrain) Move() {
	room := b.Entity.GetRoom()
	total_exits := len(room.Exits)
	exit := rand_min_max(0, total_exits)
	count := 0
	for i, e := range room.Exits {
		if count == exit {
			// only move if the room has an exit
			// this prevents mobs getting stuck in "turbolift" rooms
			to_room := DB().GetRoom(e, room.ship)
			if len(to_room.Exits) > 0 {
				do_direction(b.Entity, i)
			}
		}
		count++
	}
}

/*
	mud_prog_exec takes a string, an entity, and various argument types (entity, item, string) and

initializes a javascript vm, sets the variables, and executes the script. This is the main
function for AI scripts the mud universe calls mudprogs. Any script errors will be visible
in the console.
*/
func mud_prog_exec(vm *otto.Otto, prog string, entity Entity, any ...interface{}) error {
	ch := entity.GetCharData()
	if pg, ok := ch.Progs[prog]; ok {
		mud_prog_bind(vm, any...)
		_, err := vm.Run(pg)
		ErrorCheck(err)
		return err
	}
	return Err("%s is not a program of [%d]%s", prog, ch.Id, ch.Name)
}

/*
	mud_prog_bind sets various variables for the [GenericBrain] AI script executor.

The values are:

	$me - The current entity {string}   (mob)
	$n  - The other entity name {string}   (mob/player)
	$s  - What was said on a 'say' event {string}   (player)
*/
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
			} else if _, ok := any[i].(int); ok {
				err := vm.Set("$v", any[i].(int))
				ErrorCheck(err)
			} else if _, ok := any[i].(Item); ok {
				err := vm.Set("$i", any[i].(Item).GetData().Name)
				ErrorCheck(err)
			}
		}
	}
}

// mud_prog_init initializes a new javascript virtual machine instance for the given entity.
// It binds various mudprog functions useful for scripting mob interactions.
func mud_prog_init(entity Entity) *otto.Otto {
	vm := otto.New()
	vm.SetRandomSource(func() float64 {
		return random_float()
	})
	err := vm.Set("$me", entity.GetCharData().Name)
	if err != nil {
		panic(err)
	}
	// say("hello");
	vm.Set("say", func(call otto.FunctionCall) otto.Value {
		do_say(entity, call.Argument(0).String())
		return otto.Value{}
	})
	// shout("Stop!");
	vm.Set("shout", func(call otto.FunctionCall) otto.Value {
		do_shout(entity, call.Argument(0).String())
		return otto.Value{}
	})
	// emote("sits down");
	vm.Set("emote", func(call otto.FunctionCall) otto.Value {
		do_emote(entity, call.Argument(0).String())
		return otto.Value{}
	})
	// echo("straight to the terminal")
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
	// delay(2);  - delay($n); where $n is an integer. delay will sleep the goroutine for $n seconds.
	vm.Set("delay", func(call otto.FunctionCall) otto.Value {
		t, _ := call.Argument(0).ToInteger()
		time.Sleep(time.Duration(t) * time.Second)
		return otto.Value{}
	})
	// random(10);   - random($n); where $n is an integer. random will return a random number 0<=$n including $n.
	vm.Set("random", func(call otto.FunctionCall) otto.Value {
		arg, _ := call.Argument(0).ToInteger()
		value, _ := otto.ToValue(rand.Intn(int(arg) + 1))
		return value
	})
	// sprintf(fmt, args...);   - same as go's fmt.Sprintf but for javascript?
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
	// look();...  not sure how useful this is to the entity, maybe rework it so it makes the player ($n) perform a do_look...
	vm.Set("look", func(call otto.FunctionCall) otto.Value {
		do_look(entity)
		return otto.Value{}
	})
	// kill($n);  - makes the entity fight $n. Like scott pilgrim.
	vm.Set("kill", func(call otto.FunctionCall) otto.Value {
		target, _ := call.Argument(0).ToString()
		do_fight(entity, target)
		return otto.Value{}
	})
	// stand();  -  makes the entity stand up.
	vm.Set("stand", func(call otto.FunctionCall) otto.Value {
		do_stand(entity)
		return otto.Value{}
	})
	// sit();  -  makes the entity stand up.
	vm.Set("sit", func(call otto.FunctionCall) otto.Value {
		do_sit(entity)
		return otto.Value{}
	})
	vm.Set("give", func(call otto.FunctionCall) otto.Value {
		entity_name, _ := call.Argument(0).ToString()
		id, _ := call.Argument(1).ToInteger()
		item := DB().GetItem(uint(id))
		if item != nil {
			for _, e := range entity.GetRoom().GetEntities() {
				if e.GetCharData().Name == entity_name {
					if e.GetCharData().CurrentInventoryCount() >= e.GetCharData().MaxInventoryCount() {
						v, _ := otto.ToValue(false)
						return v
					}
					e.GetCharData().Inventory = append(e.GetCharData().Inventory, item.(*ItemData))
					e.Send("\r\n&Y have received &W%s&Y.&d\r\n", item.GetData().Name)
				}
			}
		}
		v, _ := otto.ToValue(true)
		return v
	})

	return vm
}
