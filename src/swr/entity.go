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

const (
	ENTITY_STAT_STR = iota
	ENTITY_STAT_INT
	ENTITY_STAT_DEX
	ENTITY_STAT_WIS
	ENTITY_STAT_CON
	ENTITY_STAT_CHA
)

const (
	ENTITY_STATE_NORMAL      = "normal"
	ENTITY_STATE_AGGRO       = "aggro"
	ENTITY_STATE_SEDATED     = "sedated"
	ENTITY_STATE_UNCONSCIOUS = "unconscious"
	ENTITY_STATE_SLEEPING    = "sleeping"
	ENTITY_STATE_SITTING     = "sitting"
	ENTITY_STATE_PILOTING    = "piloting"
	ENTITY_STATE_GUNNING     = "gunning"
	ENTITY_STATE_EDITING     = "editing"
	ENTITY_STATE_CRAFTING    = "crafting"
	ENTITY_STATE_DEAD        = "dead"
)

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
	IsFighting() bool
	SetAttacker(entity *Entity)
	GetCharData() *CharData
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
	Stats     []uint16        `yaml:"stats,flow"`        // str, int, dex, wis, con, cha
	Skills    map[string]int  `yaml:"skills,flow,omitempty"`
	Languages map[string]int  `yaml:"languages,flow,omitempty"`
	Speaking  string          `yaml:"speaking,omitempty"`
	Equipment map[string]Item `yaml:"equipment,flow,omitempty"` // Key is a EQUIPMENT_WEAR_LOC_* const, Value is an item.
	Inventory []Item          `yaml:"inventory,omitempty"`
	State     string          `yaml:"state,omitempty"`
	Brain     string          `yaml:"brain,omitempty"`
	AI        Brain           `yaml:"-"`
	Attacker  *Entity         `yaml:"-"`
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

func (c *CharData) IsFighting() bool {
	return c.Attacker != nil
}

func (c *CharData) SetAttacker(entity *Entity) {
	c.Attacker = entity
}

func (c *CharData) ArmorAC() uint {
	str := uint(c.Stats[ENTITY_STAT_STR])
	dex := uint(c.Stats[ENTITY_STAT_DEX])

	ac_armor := uint(0)
	for _, i := range c.Equipment {
		item := i.(ItemData)
		ac_armor += uint(item["ac"].(int))
	}
	return ac_armor + (dex / 10) + (str / 10)
}

func (c *CharData) DamageRoll() uint {
	str := uint(c.Stats[ENTITY_STAT_STR])
	dex := uint(c.Stats[ENTITY_STAT_DEX])

	dmg := uint(0)
	if i, ok := c.Equipment["weapon"]; ok {
		item := i.(ItemData)
		dmg := item["dmg"].(string)
		roll_dice(dmg)
	}
	return dmg + (str / 10) + (dex / 10)
}

func (c *CharData) ApplyDamage(damage uint) {
	c.Hp[0] -= uint16(damage)
	if c.Hp[0] < 0 {
		c.State = "unconcious"
		if c.Hp[0] < (2 * c.Hp[1]) {
			c.State = "dead"
			c.Send("\r\n&RYou have died.&d\r\n")
			return
		}
		c.Send("\r\n&YYou have been knocked unconscious...&d\r\n")
		return
	}
}

func (c *CharData) GetCharData() *CharData {
	return c
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
func (p *PlayerProfile) IsFighting() bool {
	return p.Char.Attacker != nil
}
func (p *PlayerProfile) SetAttacker(entity *Entity) {
	p.Char.Attacker = entity
}

func (p *PlayerProfile) GetCharData() *CharData {
	return &p.Char
}

func player_prompt(player *PlayerProfile) string {
	prompt := "\r\n"
	prompt += fmt.Sprintf("&Y[&GHp:&W%d&Y/&G%d&Y]&d ", player.CurrentHp(), player.MaxHp())
	prompt += fmt.Sprintf("&Y[&GMv:&W%d&Y/&G%d&Y]&d ", player.CurrentMv(), player.MaxMv())
	if player.IsFighting() {
		attacker := (*player.Char.Attacker)
		hp := attacker.MaxHp()
		chp := attacker.CurrentHp()
		third := hp / 3
		if chp < third {
			prompt += fmt.Sprintf("&w[&R%s&w]&d\n", MakeProgressBar(int(chp), int(hp), 15))
		} else if chp < third*2 {
			prompt += fmt.Sprintf("&w[&Y%s&w]&d\n", MakeProgressBar(int(chp), int(hp), 15))
		} else {
			prompt += fmt.Sprintf("&w[&G%s&w]&d\n", MakeProgressBar(int(chp), int(hp), 15))
		}
	}
	return prompt
}
