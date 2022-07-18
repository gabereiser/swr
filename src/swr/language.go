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
	"io/ioutil"
	"log"
	"math/rand"
	"strings"

	"gopkg.in/yaml.v3"
)

type Language struct {
	Race     string
	Name     string
	Alphabet string
}

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var Languages = []Language{}

func LanguageLoad() {
	if len(Languages) == 0 {
		log.Printf("Loading languages.")
		flist, err := ioutil.ReadDir("data/languages")
		ErrorCheck(err)
		for _, file := range flist {
			fp, err := ioutil.ReadFile("data/languages/" + file.Name())
			ErrorCheck(err)
			l := new(Language)
			yaml.Unmarshal(fp, l)
			for _, race := range race_list {
				if strings.EqualFold(race, l.Race) {
					Languages = append(Languages, *l)
				}
			}

		}
		log.Printf("%d languages loaded.", len(Languages))
	}
	ScheduleFunc(language_decay, true, 60*60)
}
func language_get_by_name(name string) *Language {
	for _, l := range Languages {
		if l.Name == name {
			return &l
		}
	}
	return nil
}

func language_get_rune(r rune, language *Language) string {
	r_index := strings.IndexRune(alphabet, r)
	if r_index > -1 {
		return language.Alphabet[r_index : r_index+1]
	} else {
		return string(r)
	}
}
func language_spoken(speaker *CharData, listener *CharData, words string) string {
	spoken_language := language_get_by_name(speaker.Speaking)
	if _, ok := listener.Languages[spoken_language.Name]; !ok {
		// Never heard that language before so we'll add it to the listeners' languages.
		listener.Languages[spoken_language.Name] = 0

	}
	l := listener.Languages[spoken_language.Name]
	// language-ify the word based on the characters of the word.
	word := ""
	for _, s := range words {
		// randomly see if our listeners' knowledge of a language is greater than
		// the speakers'. If the listener's better, return word as english
		if rand.Intn(speaker.Languages[spoken_language.Name]) < l || l == 100 {
			word += string(s)
			continue
		}

		r := language_get_rune(s, spoken_language)
		word += r
		// repeat the rune because we haven't quite learned the language yet (so it looks funky)
		if rand.Intn(speaker.Languages[spoken_language.Name]) > l && rand.Intn(l+1) < 5 {
			r_index := strings.IndexRune(alphabet, rune(r[0]))
			if r_index != -1 {
				word += spoken_language.Alphabet[r_index : r_index+1]
			}

		}
	}
	// give a little knowledge of the language
	if rand.Intn(5) == 0 {
		if listener.Languages[spoken_language.Name] != 100 {
			listener.Languages[spoken_language.Name]++
			listener.Send(fmt.Sprintf("&cYou gain a little knowledge of the %s language.&d\r\n", spoken_language.Name))
		}
	}
	return word
}

func language_decay() {
	for _, entity := range DB().entities {
		if entity.IsPlayer() {
			player := entity.(*PlayerProfile)
			lost := false
			for language, level := range player.Char.Languages {
				if language != player.Char.Speaking && language != strings.ToLower(player.Char.Race) && language != "basic" {
					if level != 100 {
						if rand.Intn(5) == 0 && strings.ToLower(player.Char.Race) != "droid" {
							player.Char.Languages[language] = level - 1
							lost = true
						}
					}
				}
			}
			if lost {
				player.Send("&XYou've forgotten a little bit of language knowledge.&d")
			}
		}
	}
}
func Capitalize(str string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(str[0:1]), str[1:])
}
