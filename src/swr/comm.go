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
	"io/ioutil"
	"log"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var CommandFuncs = map[string]func(Entity, ...string){
	"do_quit":           do_quit,
	"do_qui":            do_qui,
	"do_say":            do_say,
	"do_speak":          do_speak,
	"do_ooc":            do_ooc,
	"do_look":           do_look,
	"do_who":            do_who,
	"do_save":           do_save,
	"do_score":          do_score,
	"do_help":           do_help,
	"do_say_comlink":    do_say_comlink,
	"do_tune_frequency": do_tune_frequency,
	"do_fight":          do_fight,
	"do_kill":           do_kill,
	"do_starsystems":    do_starsystems,
	"do_north":          do_north,
	"do_northwest":      do_northwest,
	"do_northeast":      do_northeast,
	"do_west":           do_west,
	"do_east":           do_east,
	"do_south":          do_south,
	"do_southeast":      do_southeast,
	"do_southwest":      do_southwest,
	"do_up":             do_up,
	"do_down":           do_down,
	"do_stand":          do_stand,
	"do_sleep":          do_sleep,
	"do_sit":            do_sit,
}
var GMCommandFuncs = map[string]func(Entity, ...string){
	"do_area_create": do_area_create,
	"do_area_set":    do_area_set,
	"do_area_remove": do_area_remove,
	"do_area_reset":  do_area_reset,
	"do_room_create": do_room_create,
	"do_room_edit":   do_room_edit,
	"do_room_find":   do_room_find,
	"do_room_remove": do_room_remove,
	"do_room_reset":  do_room_reset,
	"do_room_set":    do_room_set,
	"do_mob_create":  do_mob_create,
	"do_mob_stat":    do_mob_stat,
	"do_mob_find":    do_mob_find,
	"do_mob_remove":  do_mob_remove,
	"do_mob_reset":   do_mob_reset,
	"do_mob_set":     do_mob_set,
	"do_item_create": do_item_create,
	"do_item_stat":   do_item_stat,
	"do_item_find":   do_item_find,
	"do_item_remove": do_item_remove,
	"do_item_set":    do_item_set,
}

var Commands []*Command = make([]*Command, 0)

type Command struct {
	Name     string   `yaml:"name"`
	Keywords []string `yaml:"keywords,flow"`
	Level    uint     `yaml:"level"`
	Func     string   `yaml:"func"`
}

func CommandsLoad() {
	log.Printf("Loading commands list.")
	fp, err := ioutil.ReadFile("data/sys/commands.yml")
	ErrorCheck(err)
	err = yaml.Unmarshal(fp, &Commands)
	ErrorCheck(err)
	log.Printf("Commands successfully loaded.")
}
func command_map_to_func(name string) func(Entity, ...string) {
	if k, ok := CommandFuncs[name]; ok {
		return k
	}
	if k, ok := GMCommandFuncs[name]; ok {
		return k
	}
	return do_nothing
}
func command_fuzzy_match(command string) []Command {
	ret := []Command{}
	for _, com := range Commands {
		for _, keyword := range com.Keywords {
			if len(keyword) < len(command) {
				continue
			}
			match := true
			for i, r := range command {
				if keyword[i] != byte(r) {
					match = false
				}
			}
			if match {
				ret = append(ret, *com)
			}
		}
	}
	return ret
}
func do_command(entity Entity, input string) {
	args := strings.Split(input, " ")
	if entity.IsPlayer() && input == "!" {
		player := entity.(*PlayerProfile)
		args = strings.Split(player.LastCommand, " ")
	}
	if strings.HasPrefix(args[0], "'") {
		args[0] = strings.TrimPrefix(args[0], "'")
		do_say(entity, args...)
		entity.Prompt()
	} else if strings.HasPrefix(args[0], "\"") {
		args[0] = strings.TrimPrefix(args[0], "\"")
		do_say_comlink(entity, args...)
		entity.Prompt()
		return
	}

	commands := command_fuzzy_match(args[0])
	if len(commands) > 0 {
		a := args[1:]
		command_map_to_func(commands[0].Func)(entity, a...)
		entity.Prompt()
	} else {
		commands := command_fuzzy_match(args[0])
		if len(commands) > 0 {
			a := args[1:]
			command_map_to_func(commands[0].Func)(entity, a...)
			entity.Prompt()
		} else {
			if entity.IsPlayer() {
				entity.Send("\r\nHuh?\r\n")
				entity.Prompt()
			}
		}
	}
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		player.LastSeen = time.Now().UTC()
		if input[0:1] != "!" {
			player.LastCommand = input
		}
	}

}
