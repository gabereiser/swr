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
	"fmt"
	"strconv"
	"strings"
)

func do_say(entity Entity, args ...string) {
	words := strings.Join(args, " ")
	words = strings.TrimSpace(words)
	if words == "" {
		entity.Send("\r\n&RSay what?&d\r\n")
		return
	}
	speaker := entity.GetCharData()
	if entity_unspeakable_state(entity) {
		entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
		return
	}
	if entity.IsPlayer() {
		entity.Send("You say \"%s\"\n", words)
	}
	for _, ex := range speaker.GetRoom().GetEntities() {
		if ex == nil {
			continue
		}
		if entity_unspeakable_state(ex) {
			continue
		}
		if ex != entity {
			listener := ex.GetCharData()
			lwords := language_spoken(speaker, listener, words)
			ex.Send("%s says \"%s\"\n", speaker.Name, lwords)
			if ex.GetCharData().AI != nil {
				ex.GetCharData().AI.OnSay(speaker, lwords)
			}
		}
	}
	go room_prog_exec(entity, "say", words)
}

func do_shout(entity Entity, args ...string) {
	words := strings.Join(args, " ")
	words = strings.TrimSpace(words)
	if words == "" {
		entity.Send("\r\n&RShout what?&d\r\n")
	}
	speaker := entity.GetCharData()
	if entity_unspeakable_state(entity) {
		entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
		return
	}
	if entity.IsPlayer() {
		entity.Send("You shout \"%s\"!\n", words)
	}
	yell(entity, words, speaker.RoomId(), 0, make([]uint, 0))
}

func yell(entity Entity, words string, roomId uint, dist uint, visited []uint) {
	if dist > 3 {
		return
	}
	visited_already := false
	for _, rId := range visited {
		if rId == roomId {
			visited_already = true
		}
	}
	if !visited_already {
		speaker := entity.GetCharData()
		room := DB().GetRoom(roomId, speaker.Ship)
		for _, ex := range room.GetEntities() {
			if ex == nil {
				continue
			}
			if entity_unspeakable_state(ex) {
				continue
			}
			if ex != entity {
				if ex.IsPlayer() {
					listener := ex.GetCharData()
					ex.Send("Someone shouts \"%s\"!\n", language_spoken(speaker, listener, words))
				} else {
					ex.Send("Someone shouts \"%s\"!\n", words)
				}
			}
		}
		visited = append(visited, roomId)
		for d, e := range room.Exits {
			if exf, ok := room.ExitFlags[d]; ok {
				if exf.Closed {
					continue // prevents yells from going through doors...  may revisit this in the future.
				}
			}
			yell(entity, words, e, dist+1, visited)
		}
	}
}

func do_emote(entity Entity, args ...string) {
	emote := strings.Join(args, " ")
	speaker := entity.GetCharData()
	if entity_unspeakable_state(entity) {
		entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
		return
	}
	speaker.GetRoom().SendToRoom(sprintf("&d%s %s&d\r\n", speaker.Name, emote))
}
func do_say_comlink(entity Entity, args ...string) {
	words := strings.Join(args, " ")
	words = strings.TrimSpace(words)
	speaker := entity.GetCharData()
	speaker_freq := entity.(*PlayerProfile).Frequency
	if entity_unspeakable_state(entity) {
		entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
		return
	}
	if words == "" {
		entity.Send("\r\n%s\r\n", MakeTitle("Comlink Status", ANSI_TITLE_STYLE_SYSTEM, ANSI_TITLE_ALIGNMENT_LEFT))
		entity.Send("&GComlink&d: %-32s\r\n\r\n", "PIC//113 Kuat Systems Intercom")
		entity.Send("&G----[&W%s&G]----&d\r\n", MakeTunerBar(speaker_freq, 50))
		entity.Send("&G Freq&d: &Y%s mhz&d\r\n", speaker_freq)
		return
	}
	if entity.IsPlayer() {
		entity.Send("You're comlink hums after you say &W\"%s\"&d\r\n", words)
	}
	db := DB()
	for _, ex := range db.entities {
		if ex == nil {
			continue
		}
		if entity_unspeakable_state(ex) {
			continue
		}
		if ex != entity {
			if ex.IsPlayer() {
				listener := ex.GetCharData()
				listener_freq := ex.(*PlayerProfile).Frequency
				if listener_freq == speaker_freq {
					ex.Send("&CYou're comlink crackles to life with a voice that says...&d\r\n\"&W%s&Y:&d %s\"\r\n", speaker.Name, language_spoken(speaker, listener, words))
				}
			}
		}
	}
}

//lint:ignore U1000 useful code
func do_broadcast_comlink(freq string, message string) {
	db := DB()
	for _, ex := range db.entities {
		if ex == nil {
			continue
		}
		if entity_unspeakable_state(ex) {
			continue
		}
		if ex.IsPlayer() {
			listener_freq := ex.(*PlayerProfile).Frequency
			if listener_freq == freq {
				ex.Send("&CYou're comlink crackles to life with a voice that says...&d\r\n\"%s\"\r\n", message)
			}
		}
	}
}

func do_tune_frequency(entity Entity, args ...string) {
	if entity.IsPlayer() {
		player := entity.(*PlayerProfile)
		if entity_unspeakable_state(entity) {
			entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
			return
		}
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
			freq_parts := strings.Split(freq_str, ".")
			if freq_parts[1] != "000" && freq_parts[1] != "250" && freq_parts[1] != "500" && freq_parts[1] != "750" {
				entity.Send("\r\n&RInvalid frequency. Frequencies are 250hz increment on the mhz interstellar bands.&d")
				entity.Send("example: 150.000 425.250 320.750 200.500\r\n")
				return
			}
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
	if entity_unspeakable_state(entity) {
		entity.Send("\r\n&dYou are %s.&d\r\n", entity_unspeakable_reason(entity))
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
