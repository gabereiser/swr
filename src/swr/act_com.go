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
	"strconv"
	"strings"
)

func do_say(entity Entity, args ...string) {
	words := strings.Join(args, " ")
	speaker := entity.GetCharData()
	if entity.IsPlayer() {
		entity.Send(fmt.Sprintf("You say \"%s\"\n", words))
	}
	entities := DB().GetEntitiesInRoom(speaker.RoomId())
	for _, ex := range entities {
		if ex != entity {
			if ex.IsPlayer() {
				listener := ex.GetCharData()
				ex.Send(fmt.Sprintf("%s says \"%s\"\n", speaker.Name, language_spoken(speaker, listener, words)))
			} else {
				ex.Send(fmt.Sprintf("%s says \"%s\"\n", speaker.Name, words))
			}
		}
	}
}
func do_say_comlink(entity Entity, args ...string) {
	words := strings.Join(args, " ")
	speaker := entity.GetCharData()
	if entity.IsPlayer() {
		entity.Send(fmt.Sprintf("You're comlink hums after you say &W\"%s\"&d\r\n", words))
	}
	db := DB()
	for _, ex := range db.entities {
		if ex != entity {
			if ex.IsPlayer() {
				listener := ex.GetCharData()
				ex.Send(fmt.Sprintf("&CYou're comlink crackles to life with a voice that says...&d\r\n\"&W%s&Y:&d %s\"\r\n", speaker.Name, language_spoken(speaker, listener, words)))
			}
		}
	}
}

func do_tune_frequency(entity Entity, args ...string) {
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		if len(args) > 0 {
			freq, err := strconv.ParseFloat(args[0], 32)
			if err != nil {
				entity.Send("\r\n&RError parsing frequency!&d\r\n")
				return
			}
			if freq < 100.000 || freq > 500.000 {
				entity.Send("\r\n&RFrequency out-of-band of your comlink!&d\r\n")
				return
			}
			freq_str := fmt.Sprintf("%3.3f", freq)
			player.Frequency = freq_str
		} else {
			player.Frequency = tune_random_frequency()
		}
		entity.Send("\r\n&YYou're comlink frequency has been set to &W%s&Y.&d\r\n", player.Frequency)
	}
}

func do_speak(entity Entity, args ...string) {
	if len(args) == 0 {
		entity.Send("\r\n&CSyntax: speak <language>&d\r\n")
		return
	}
	ch := entity.GetCharData()
	language := language_get_by_name(args[0])
	if language != nil {
		skill := ch.Languages[language.Name]
		if skill > 0 {
			ch.Speaking = language.Name
			entity.Send("\r\n&YYou are now speaking &W%s&d.\r\n", language.Name)
		} else {
			entity.Send("\r\n&YLanguage not known.&d\r\n")
		}
	} else {
		entity.Send("\r\n&YLanguage not known.&d\r\n")
	}
}
