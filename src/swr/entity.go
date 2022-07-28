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
	"math"
	"strings"
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
	ENTITY_STATE_FIGHTING    = "fighting"
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
	IsPlayer() bool
	Send(str string, any ...interface{})
	Event(evt string)
	Prompt()
	CurrentHp() int
	MaxHp() int
	CurrentMv() int
	MaxMv() int
	IsFighting() bool
	StopFighting()
	SetAttacker(entity Entity)
	ApplyDamage(damage uint)
	GetCharData() *CharData
	Weapon() Item
	FindItem(keyword string) Item
}

type CharData struct {
	Id        uint                 `yaml:"id"`
	OId       uint                 `yaml:"mobId,omitempty"`
	Room      uint                 `yaml:"room,omitempty"`
	Name      string               `yaml:"name"`
	Keywords  []string             `yaml:"keywords,flow,omitempty"`
	Title     string               `yaml:"title,omitempty"`
	Desc      string               `yaml:"desc"`
	Race      string               `yaml:"race,omitempty"`
	Gender    string               `yaml:"gender,omitempty"`
	Level     uint                 `yaml:"level,omitempty"`
	XP        uint                 `yaml:"xp,omitempty"`
	Gold      uint                 `yaml:"gold,omitempty"`
	Bank      uint                 `yaml:"bank,omitempty"`
	Hp        []int                `yaml:"hp,flow"`    // Hit Points [0] Current [1] Max : len = 2
	Mp        []int                `yaml:"mp,flow"`    // Magic Points [0] Current [1] Max : len = 2
	Mv        []int                `yaml:"mv,flow"`    // Move Points [0] Current [1] Max : len = 2
	Stats     []int                `yaml:"stats,flow"` // str, int, dex, wis, con, cha
	Skills    map[string]int       `yaml:"skills,flow,omitempty"`
	Languages map[string]int       `yaml:"languages,flow,omitempty"`
	Speaking  string               `yaml:"speaking,omitempty"`
	Equipment map[string]*ItemData `yaml:"equipment,flow,omitempty"` // Key is a EQUIPMENT_WEAR_LOC_* const, Value is an item.
	Inventory []*ItemData          `yaml:"inventory,omitempty"`
	State     string               `yaml:"state,omitempty"`
	Brain     string               `yaml:"brain,omitempty"`
	Progs     map[string]string    `yaml:"progs,omitempty"`
	Flags     []string             `yaml:"flags,omitempty"`
	AI        Brain                `yaml:"-"`
	Attacker  Entity               `yaml:"-"`
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

func (c *CharData) Event(evt string) {
}

func (c *CharData) CurrentHp() int {
	return c.Hp[0]
}

func (c *CharData) MaxHp() int {
	return c.Hp[1]
}

func (c *CharData) CurrentMv() int {
	return c.Mv[0]
}

func (c *CharData) MaxMv() int {
	return c.Mv[1]
}

func (c *CharData) base_weight() int {
	switch c.Race {
	case "Wookiee":
		return 105
	case "Hutt":
		return 425
	case "Ewok":
	case "Jawa":
		return 25
	case "Droid":
	case "Protocol Droid":
		return 95
	case "Assassin Droid":
	case "Gladiator Droid":
		return 245
	}
	return 75
}

func (c *CharData) CurrentWeight() int {
	weight := c.base_weight()
	for _, item := range c.Inventory {
		if item == nil {
			continue
		}
		weight += item.GetWeight()
	}
	return weight
}

func (c *CharData) MaxWeight() int {
	weight := c.base_weight()

	// str / 10 * base_weight + (level * 5) + dex / 10 * base_weight
	return ((c.Stats[0] / 10) * weight) + int(c.Level*5) + ((c.Stats[2] / 10) * weight)
}

func (c *CharData) CurrentInventoryCount() int {
	count := 0
	for range c.Inventory {
		count++
	}
	return count
}

func (c *CharData) MaxInventoryCount() int {
	return (int(c.Level) * 3) + c.Stats[0]
}

func (c *CharData) IsFighting() bool {
	return c.State == ENTITY_STATE_FIGHTING
}

func (c *CharData) StopFighting() {
	if c.State == ENTITY_STATE_FIGHTING {
		c.State = ENTITY_STATE_NORMAL
		c.Attacker = nil
		c.Send("\r\n&dYou stop fighting.\r\n")
	}
}

func (c *CharData) SetAttacker(entity Entity) {
	c.Attacker = entity
	c.State = ENTITY_STATE_FIGHTING
}

func (c *CharData) ArmorAC() int {
	str := c.Stats[ENTITY_STAT_STR]
	dex := c.Stats[ENTITY_STAT_DEX]

	ac_armor := 0
	for _, i := range c.Equipment {
		item := i.GetData()
		ac_armor += item.AC
	}
	return ac_armor + (dex / 10) + (str / 10)
}

func (c *CharData) DamageRoll(skillName string) uint {
	skill := uint(c.Skills[skillName])
	str := uint(c.Stats[ENTITY_STAT_STR])
	dex := uint(c.Stats[ENTITY_STAT_DEX])
	d := "1d4"
	i := c.Weapon()
	if i != nil {
		item := i.GetData()
		d = *item.Dmg
		skill += get_weapon_skill_stat(*item.WeaponType, str, dex)
	} else {
		skill += get_weapon_skill_stat("martial-arts", str, dex)
	}
	if skill == 0 {
		skill = 1
	}
	dmg := uint(roll_dice(d)) + uint(roll_dice(fmt.Sprintf("1d%d", skill)))
	return dmg
}

func (c *CharData) Weapon() Item {
	if i, ok := c.Equipment["weapon"]; ok {
		item := i
		return item
	}
	return nil
}

func (c *CharData) FindItem(keyword string) Item {
	for i := range c.Inventory {
		item := c.Inventory[i]
		if item == nil {
			continue
		}
		keys := item.GetKeywords()
		for k := range keys {
			key := keys[k]
			if strings.HasPrefix(key, keyword) {
				return item
			}
		}
	}
	return nil
}

func (c *CharData) RemoveItem(item Item) {
	idx := -1
	for id := range c.Inventory {
		i := c.Inventory[id]
		if i == nil {
			continue
		}
		if i == item {
			idx = id
		} else if i.IsContainer() {
			i.GetData().RemoveItem(item)
		}
	}
	if idx > -1 {
		ret := make([]*ItemData, 0)
		for id, i := range c.Inventory {
			if id == idx {
				continue
			}
			if i == nil {
				continue
			}
			ret = append(ret, i)
		}
		c.Inventory = ret
	}
}

func (c *CharData) GetItem(item_id uint) Item {
	for id := range c.Inventory {
		i := c.Inventory[id]
		if i.GetId() == item_id || i.GetData().OId == item_id {
			return i
		}
	}
	return nil
}

func (c *CharData) ApplyDamage(damage uint) {
	c.Hp[0] -= int(damage)
	if c.Hp[0] <= 0 {
		c.State = ENTITY_STATE_DEAD
		if c.AI != nil {
			c.AI.OnDeath()
		}
	}
}

func (c *CharData) GetCharData() *CharData {
	return c
}

type PlayerProfile struct {
	Char        CharData  `yaml:"char,inline"`
	Email       string    `yaml:"email,omitempty"`
	Password    string    `yaml:"password,omitempty"`
	Priv        int       `yaml:"priv,omitempty"`
	LastSeen    time.Time `yaml:"last_seen,omitempty"`
	Banned      bool      `yaml:"banned,omitempty"`
	Frequency   string    `yaml:"freq"`
	Client      Client    `yaml:"-"`
	NeedPrompt  bool      `yaml:"-"`
	LastCommand string    `yaml:"-"`
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

func (p *PlayerProfile) Event(evt string) {
}

func (p *PlayerProfile) CurrentHp() int {
	return p.Char.Hp[0]
}

func (p *PlayerProfile) MaxHp() int {
	return p.Char.Hp[1]
}

func (p *PlayerProfile) CurrentMv() int {
	return p.Char.Mv[0]
}

func (p *PlayerProfile) MaxMv() int {
	return p.Char.Mv[1]
}
func (p *PlayerProfile) Prompt() {
	if p.NeedPrompt && p.Char.State != ENTITY_STATE_DEAD {
		prompt := player_prompt(p)
		p.Send("%s\r\n", prompt)
		p.NeedPrompt = false
	}
}
func (p *PlayerProfile) IsFighting() bool {
	return p.Char.State == ENTITY_STATE_FIGHTING
}
func (p *PlayerProfile) StopFighting() {
	if p.Char.State == ENTITY_STATE_FIGHTING {
		p.Char.State = ENTITY_STATE_NORMAL
		p.Char.Attacker = nil
		p.Send("\r\n&dYou stop fighting.\r\n")
	}
}
func (p *PlayerProfile) SetAttacker(entity Entity) {
	p.Char.Attacker = entity
	p.Char.State = ENTITY_STATE_FIGHTING
}

func (p *PlayerProfile) GetCharData() *CharData {
	return &p.Char
}
func (p *PlayerProfile) ApplyDamage(damage uint) {
	p.Char.Hp[0] -= int(damage)
	if p.Char.Hp[0] <= 0 {
		p.Char.State = ENTITY_STATE_UNCONSCIOUS
		if p.Char.Hp[0] <= -(p.Char.Hp[1]) {
			p.Char.State = ENTITY_STATE_DEAD
			p.Send(sprintf("\r\n&RYou have died.&d %s\r\n", EMOJI_SKULL))
			return
		}
		p.Send("\r\n&YYou have been knocked unconscious...&d\r\n")
		return
	}
}

func (p *PlayerProfile) Weapon() Item {
	if i, ok := p.Char.Equipment["weapon"]; ok {
		item := i
		return item
	}
	return nil
}

func (p *PlayerProfile) FindItem(keyword string) Item {
	return p.Char.FindItem(keyword)
}

func entity_clone(entity Entity) Entity {
	ch := entity.GetCharData()
	c := &CharData{
		Id:        gen_npc_char_id(),
		OId:       ch.Id,
		Room:      ch.Room,
		Name:      ch.Name,
		Keywords:  make([]string, 0),
		Flags:     make([]string, 0),
		Title:     ch.Title,
		Desc:      ch.Desc,
		Race:      ch.Race,
		Gender:    ch.Gender,
		Level:     ch.Level,
		XP:        ch.XP,
		Gold:      ch.Gold,
		Bank:      ch.Bank,
		Speaking:  ch.Speaking,
		Hp:        make([]int, 2),
		Mp:        make([]int, 2),
		Mv:        make([]int, 2),
		Stats:     make([]int, 6),
		Skills:    make(map[string]int),
		Equipment: make(map[string]*ItemData),
		Inventory: make([]*ItemData, 0),
		Languages: make(map[string]int),
		AI:        ch.AI,
		State:     ch.State,
		Brain:     ch.Brain,
		Attacker:  ch.Attacker,
	}
	for i := range ch.Keywords {
		k := ch.Keywords[i]
		c.Keywords = append(c.Keywords, k)
	}
	for f := range ch.Flags {
		k := ch.Flags[f]
		c.Flags = append(c.Flags, k)
	}
	c.Hp[0] = ch.Hp[0]
	c.Hp[1] = ch.Hp[1]
	c.Mp[0] = ch.Mp[0]
	c.Mp[1] = ch.Mp[1]
	c.Mv[0] = ch.Mv[0]
	c.Mv[1] = ch.Mv[1]
	for i := range ch.Stats {
		s := ch.Stats[i]
		c.Stats[i] = s
	}
	for wearLoc, item := range ch.Equipment {
		c.Equipment[wearLoc] = item_clone(item).GetData()
	}
	for language, level := range ch.Languages {
		c.Languages[language] = level
	}
	for i := range ch.Inventory {
		item := ch.Inventory[i]
		if item == nil {
			continue
		}
		c.Inventory = append(c.Inventory, item)
	}
	for s, v := range ch.Skills {
		c.Skills[s] = v
	}
	return c
}

func player_prompt(player *PlayerProfile) string {
	mc := player.Client.(*MudClient)
	if mc.Closed {
		return sprintf("Thank you for playing! %s\r\n", EMOJI_ALERT)
	}
	if player.Char.State == ENTITY_STATE_DEAD {
		return "&R*DEAD&d\r\n"
	}
	prompt := "\r\n"
	prompt += fmt.Sprintf("&Y[&GHp:&W%d&Y/&G%d&Y]&d ", player.CurrentHp(), player.MaxHp())
	prompt += fmt.Sprintf("&Y[&GMv:&W%d&Y/&G%d&Y]&d ", player.CurrentMv(), player.MaxMv())
	if player.IsFighting() {
		attacker := player.Char.Attacker
		hp := attacker.MaxHp()
		chp := attacker.CurrentHp()
		third := hp / 3
		if chp < third {
			prompt += fmt.Sprintf("&w[&R%s&w]&d", MakeProgressBar(int(chp), int(hp), 15))
		} else if chp < third*2 {
			prompt += fmt.Sprintf("&w[&Y%s&w]&d", MakeProgressBar(int(chp), int(hp), 15))
		} else {
			prompt += fmt.Sprintf("&w[&G%s&w]&d", MakeProgressBar(int(chp), int(hp), 15))
		}
	}
	return sprintf("%s\r\n", prompt)
}

func processEntities() {
	db := DB()

	for i := range db.entities {
		e := db.entities[i]
		if e == nil {
			continue
		}
		ch := e.GetCharData()
		switch ch.State {
		case ENTITY_STATE_NORMAL:
			if roll_dice("1d10") >= 10-entity_get_skill_value(ch, "healing") {
				processHealing(e)
			}
		case ENTITY_STATE_SITTING:
			if roll_dice("1d10") >= 8-entity_get_skill_value(ch, "healing") {
				processHealing(e)
				processHealing(e)

			}
		case ENTITY_STATE_SLEEPING:
			processHealing(e)
		case ENTITY_STATE_UNCONSCIOUS:
			if roll_dice("1d10") >= 5-entity_get_skill_value(ch, "healing") {
				processHealing(e)
			}
		case ENTITY_STATE_SEDATED:
			ch.Mv[0]--
			if ch.Mv[0] < 0 {
				ch.Mv[0] = 0
			}
			if roll_dice("1d10") == 10 {
				e.Send("\r\n&cYou feel extremely relaxed.&d\r\n")
				e.Prompt()
			}
		}
		if roll_dice("1d10") == 10 {
			if ch.Mp[0] < ch.Mp[1] {
				ch.Mp[0]++
				if ch.Mp[0] > ch.Mp[1] {
					ch.Mp[0] = ch.Mp[1]
				}
			}
			if ch.Mv[0] < ch.Mv[1] {
				ch.Mv[0]++
				if ch.Mv[0] > ch.Mv[1] {
					ch.Mv[0] = ch.Mv[1]
				}
			}
		}

	}
}
func processHealing(entity Entity) {
	ch := entity.GetCharData()
	ch.Hp[0]++
	if ch.Hp[0] > ch.Hp[1] {
		ch.Hp[0] = ch.Hp[1]
	}
	ch.Mv[0]++
	if ch.Mv[0] > ch.Mv[1] {
		ch.Mv[0] = ch.Mv[1]
	}
	if ch.State == ENTITY_STATE_UNCONSCIOUS {
		if ch.Hp[0] > 0 {
			ch.State = ENTITY_STATE_NORMAL
			entity.Send("\r\n&YYou awake from unconsciousness.&d\r\n")
			entity.Prompt()
		}
	}
	entity.Prompt()
}

func entity_add_xp(entity Entity, xp int) {
	// mobs don't earn xp
	if !entity.IsPlayer() {
		return
	}
	ch := entity.GetCharData()
	level := ch.Level
	x := int(ch.XP)
	x += xp
	if x <= 0 {
		x = 0
	}
	ch.XP = uint(x)
	ch.Level = get_level_for_xp(ch.XP) + 1
	if ch.Level != level {
		entity.Send("\r\n}YYou have gained a level!&d\r\n")
		entity.Send("\r\n&YYou are now level &W%d&d.\r\n", ch.Level)
		// reset current life stats as a reward and gain a little extra
		ch.Hp[1] = 10 + (int(ch.Level) / 2)
		ch.Mv[1] = 10 + (int(ch.Level) / 2)
		if ch.Mp[1] > 0 { // you don't get force unless you got force
			ch.Mp[1] = 10 + (int(ch.Level) / 2)
		}
		ch.Hp[0] = ch.Hp[1]
		ch.Mp[0] = ch.Mp[1]
		ch.Mv[0] = ch.Mv[1]
	}
	entity.Send("\r\n&dYou gained &w%d&d xp.\r\n", xp)
}
func entity_lose_xp(entity Entity, xp int) {
	if !entity.IsPlayer() {
		return
	}
	ch := entity.GetCharData()
	level := ch.Level
	x := int(ch.XP)
	x -= xp
	if x <= 0 {
		x = 0
	}
	ch.XP = uint(x)
	ch.Level = get_level_for_xp(ch.XP) + 1
	entity.Send("\r\n&dYou lost &w%d&d xp.\r\n", xp)
	if ch.Level != level {
		entity.Send("\r\n}RYou have lost a level!&d\r\n")
		entity.Send("\r\n&RYou are now level &W%d&d.\r\n", ch.Level)
	}
}

func get_level_for_xp(xp uint) uint {
	return uint(math.Sqrt(float64(xp) / 500))
}

func get_xp_for_level(level uint) uint {
	return uint(math.Pow(float64(level), 2)) * 500
}

func entity_unspeakable_state(entity Entity) bool {
	if entity == nil {
		return true
	}
	state := entity.GetCharData().State
	switch state {
	case ENTITY_STATE_DEAD:
	case ENTITY_STATE_UNCONSCIOUS:
	case ENTITY_STATE_SLEEPING:
		return true
	}
	return false
}

func entity_unspeakable_reason(entity Entity) string {
	state := entity.GetCharData().State
	switch state {
	case ENTITY_STATE_DEAD:
		return "dead"
	case ENTITY_STATE_UNCONSCIOUS:
		return "unconscious"
	case ENTITY_STATE_SLEEPING:
		return "sleeping"
	}
	return "none"
}

func entity_pickup_item(entity Entity, item Item) bool {
	ch := entity.GetCharData()
	if item.IsCorpse() {
		entity.Send("\r\n&RYou can't carry a corpse.&d\r\n")
		return false
	}
	if item.GetData().Weight+ch.CurrentWeight() > ch.MaxWeight() {
		entity.Send("\r\n&RYou can't carry any more weight!&d\r\n")
		return false
	}
	if ch.CurrentInventoryCount() >= ch.MaxInventoryCount() {
		entity.Send("\r\n&RYou can't carry any more items!&d\r\n")
		return false
	}
	ch.Inventory = append(ch.Inventory, item.GetData())
	return true
}

func entity_get_skill_value(ch *CharData, skill string) int {
	if v, ok := ch.Skills[strings.ToLower(skill)]; ok {
		return v
	}
	return 0
}

func entity_add_skill_value(entity Entity, skill string, value int) {
	ch := entity.GetCharData()
	ch.Skills[skill] += value
	if ch.Skills[skill] >= 100 {
		ch.Skills[skill] = 100
	} else {
		entity.Send("\r\n&CYou gain some knowledge of %s.&d\r\n", skill)
	}

}

func entity_get_equipment_for_slot(entity Entity, wearLoc string) string {
	if o, ok := entity.GetCharData().Equipment[wearLoc]; ok {
		return o.Name
	}
	return "None"
}
