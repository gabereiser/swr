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
	"time"
)

var race_list = []string{
	"Human", "Wookiee", "Twi'lek", "Rodian", "Hutt", "Mon Calamari", "Noghri",
	"Gamorrean", "Jawa", "Adarian", "Ewok", "Verpine", "Defel", "Trandoshan",
	"Hapan", "Quarren", "Shistavanen", "Falleen", "Ithorian", "Devaronian", "Gotal", "Droid",
	"Firrerreo", "Barabel", "Bothan", "Togorian", "Dug", "Kubaz", "Selonian", "Gran", "Yevetha", "Gand",
	"Duros", "Coynite", "Sullustan", "Protocol Droid", "Assassin Droid", "Gladiator Droid", "Astromech Droid",
	"Interrogation Droid", "Sarlacc", "Saurin", "Snivvian", "Gand", "Gungan", "Weequay", "Bith",
	"Ortolan", "Snit", "Cerean", "Ugnaught", "Taun Taun", "Bantha", "Tusken",
	"Gherkin", "Zabrak", "Dewback", "Rancor", "Ronto",
	"Monster",
}

type Entity interface {
	RoomId() uint
	Name() string
	IsPlayer() bool
	Send(str string, any ...interface{})
	Event(evt string)
	Prompt()
	CurrentHp() uint16
	MaxHp() uint16
	CurrentMv() uint16
	MaxMv() uint16
}

type CharData struct {
	Room      uint            `yaml:"room,omitempty"`
	CharName  string          `yaml:"name"`
	Keywords  []string        `yaml:"keywords,flow,omitempty"`
	Title     string          `yaml:"title,omitempty"`
	Desc      string          `yaml:"desc"`
	Race      string          `yaml:"race,omitempty"`
	Gender    string          `yaml:"gender,omitempty"`
	Level     uint16          `yaml:"level,omitempty"`
	XP        uint            `yaml:"xp,omitempty"`
	Gold      uint            `yaml:"gold,omitempty"`
	Bank      uint            `yaml:"bank,omitempty"`
	Hp        []uint16        `yaml:"hp,flow"`           // Hit Points [0] Current [1] Max : len = 2
	Mp        []uint16        `yaml:"mp,flow"`           // Magic Points [0] Current [1] Max : len = 2
	Mv        []uint16        `yaml:"mv,flow,omitempty"` // Move Points [0] Current [1] Max : len = 2
	Stats     []uint16        `yaml:"stats,flow"`
	Skills    map[string]int  `yaml:"skills,flow,omitempty"`
	Languages map[string]int  `yaml:"languages,flow,omitempty"`
	Speaking  string          `yaml:"speaking,omitempty"`
	Equipment map[string]Item `yaml:"equipment,flow,omitempty"` // Key is a EQUIPMENT_WEAR_LOC_* const, Value is an item.
	Inventory []Item          `yaml:"inventory,omitempty"`
	State     string          `yaml:"state,omitempty"`
	Brain     string          `yaml:"brain,omitempty"`
	AI        Brain           `yaml:"-"`
}

func (*CharData) IsPlayer() bool {
	return false
}

func (*CharData) Prompt() {
}

func (*CharData) Send(m string, args ...interface{}) {
}

func (c *CharData) RoomId() uint {
	return c.Room
}

func (c *CharData) Name() string {
	return c.CharName
}

func (c *CharData) Event(evt string) {
}

func (c *CharData) CurrentHp() uint16 {
	return c.Hp[0]
}

func (c *CharData) MaxHp() uint16 {
	return c.Hp[1]
}

func (c *CharData) CurrentMv() uint16 {
	return c.Mv[0]
}

func (c *CharData) MaxMv() uint16 {
	return c.Mv[1]
}

func (c *CharData) CurrentWeight() int {
	weight := 75
	for _, item := range c.Inventory {
		weight += item.GetWeight()
	}
	return weight
}

func (c *CharData) CurrentInventoryCount() int {
	return len(c.Inventory)
}

type PlayerProfile struct {
	Char       CharData  `yaml:"char,inline"`
	Email      string    `yaml:"email,omitempty"`
	Password   string    `yaml:"password,omitempty"`
	Priv       int       `yaml:"priv,omitempty"`
	LastSeen   time.Time `yaml:"last_seen,omitempty"`
	Banned     bool      `yaml:"banned,omitempty"`
	Client     Client    `yaml:"-"`
	NeedPrompt bool      `yaml:"-"`
}

func (*PlayerProfile) IsPlayer() bool {
	return true
}

func (p *PlayerProfile) Send(m string, any ...interface{}) {
	if p.Client != nil {
		p.Client.Sendf(m, any...)
		p.NeedPrompt = true
	}
}

func (p *PlayerProfile) RoomId() uint {
	return p.Char.Room
}

func (p *PlayerProfile) Name() string {
	return p.Char.CharName
}

func (p *PlayerProfile) Event(evt string) {
}

func (p *PlayerProfile) CurrentHp() uint16 {
	return p.Char.Hp[0]
}

func (p *PlayerProfile) MaxHp() uint16 {
	return p.Char.Hp[1]
}

func (p *PlayerProfile) CurrentMv() uint16 {
	return p.Char.Mv[0]
}

func (p *PlayerProfile) MaxMv() uint16 {
	return p.Char.Mv[1]
}
func (p *PlayerProfile) Prompt() {
	if p.NeedPrompt {
		prompt := player_prompt(p)
		p.Send("%s\r\n", prompt)
		p.NeedPrompt = false
	}
}

func player_prompt(player *PlayerProfile) string {
	prompt := "\r\n"
	prompt += fmt.Sprintf("&Y[&GHp:&W%d&Y/&G%d&Y]&d ", player.CurrentHp(), player.MaxHp())
	prompt += fmt.Sprintf("&Y[&GMv:&W%d&Y/&G%d&Y]&d ", player.CurrentMv(), player.MaxMv())
	return prompt
}
