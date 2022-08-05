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
	ENTITY_STAT_STR = iota // [0] Strength
	ENTITY_STAT_INT        // [1] Intelligence
	ENTITY_STAT_DEX        // [2] Dexterity
	ENTITY_STAT_WIS        // [3] Wisdom
	ENTITY_STAT_CON        // [4] Constitution
	ENTITY_STAT_CHA        // [5] Charisma
)

const (
	ENTITY_STATE_NORMAL      = "normal"      // everything's fine.
	ENTITY_STATE_FIGHTING    = "fighting"    // conner mcgreggor mode
	ENTITY_STATE_SEDATED     = "sedated"     // fear and loathing
	ENTITY_STATE_UNCONSCIOUS = "unconscious" // comatose
	ENTITY_STATE_SLEEPING    = "sleeping"    // because we all need rest
	ENTITY_STATE_SITTING     = "sitting"     // can't move, but you can see (and gain a bonus to hp regen)
	ENTITY_STATE_PILOTING    = "piloting"    // chuck yaeger / wedge antilles
	ENTITY_STATE_GUNNING     = "gunning"     // millenium falcon turret mode
	ENTITY_STATE_EDITING     = "editing"     // set when the admin is editing a description so that his char is safe from harm
	ENTITY_STATE_CRAFTING    = "crafting"    // set when crafting the next big item.
	ENTITY_STATE_DEAD        = "dead"        // EOL...
)

type Entity interface {
	RoomId() uint                        // the current room id.
	ShipId() uint                        // the current ship id. 0 value means they currently aren't on a ship.
	IsPlayer() bool                      // is the entity a player? or a mob?
	Send(str string, any ...interface{}) // send sprintf style.
	Event(evt string)                    // trigger a non-AI event (long speak prog?)
	Prompt()                             // shows the player prompt at next tick.
	CurrentHp() int                      // current hitpoints
	MaxHp() int                          // maximum hitpoints
	CurrentMv() int                      // current movement
	MaxMv() int                          // maximum movement
	IsFighting() bool                    // is fighting someone?
	StopFighting()                       // stop fighting whomever it's fighting.
	SetAttacker(entity Entity)           // starting fighting an entity. Next tick will commence combat.
	ApplyDamage(damage uint)             // apply damage to this entity. If it's hp < 0 or otherwise unconscious then the [CharData.State] will change.
	GetCharData() *CharData              // get the backing [CharData]
	Weapon() Item                        // get the current weapon. nil means he's bruce lee and fights with his fists. bravo good sir, bravo.
	FindItem(keyword string) Item        // find an item on this entity by keyword. If multiple are found, it will return the first found.
	GetShip() Ship                       // get the current ship the player is in. Not that they own it, but physically inside.
	GetRoom() *RoomData                  // get the room the entity is in.
}

type CharData struct {
	Id        uint                 `yaml:"id"`                      // instance id. Will always be unique to a spawn.
	OId       uint                 `yaml:"mobId,omitempty"`         // type id. What kind of mob is it? check [GameDatabase.Mobs]
	Room      uint                 `yaml:"room,omitempty"`          // room id.
	Ship      uint                 `yaml:"ship,omitempty"`          // ship id. if 0, entity is not on a ship
	Name      string               `yaml:"name"`                    // character name
	Filename  string               `yaml:"-"`                       // mob filename as used in ./data/mobs/<areaname>/<filename>.yml
	Keywords  []string             `yaml:"keywords,flow,omitempty"` // keywords to refer to this mob
	Title     string               `yaml:"title,omitempty"`         // titles granted
	Desc      string               `yaml:"desc"`                    // description of mob
	Race      string               `yaml:"race,omitempty"`          // race name from [race_list]
	Gender    string               `yaml:"gender,omitempty"`        // single char gender, lowercase. m/f/n
	Level     uint                 `yaml:"level,omitempty"`         // character level. 100 is max level.
	XP        uint                 `yaml:"xp,omitempty"`            // character xp.
	Gold      uint                 `yaml:"gold,omitempty"`          // character money on hand.
	Bank      uint                 `yaml:"bank,omitempty"`          // character money in bank.
	Hp        []int                `yaml:"hp,flow"`                 // Hit Points [0] Current [1] Max : len = 2
	Mp        []int                `yaml:"mp,flow"`                 // Magic Points [0] Current [1] Max : len = 2
	Mv        []int                `yaml:"mv,flow"`                 // Move Points [0] Current [1] Max : len = 2
	Stats     []int                `yaml:"stats,flow"`              // str, int, dex, wis, con, cha
	Skills    map[string]int       `yaml:"skills"`                  // skills map. skills are indexed by skill name with a value of 0-100.
	Languages map[string]int       `yaml:"languages"`               // languages known. languages are indexed by language name with a value of 0-100
	Speaking  string               `yaml:"speaking"`                // what language are we speaking? should match a key in [CharData.Languages]
	Equipment map[string]*ItemData `yaml:"equipment"`               // equipment map.key is a EQUIPMENT_WEAR_LOC_* const, value is an [ItemData].
	Inventory []*ItemData          `yaml:"inventory"`               // inventory list. multiple items (with different id's) can be stored.
	State     string               `yaml:"state,omitempty"`         // character state as defined as an ENTITY_STATE_* const
	Brain     string               `yaml:"brain,omitempty"`         // character brain. essentially the ai class. generic is the default.
	Progs     map[string]string    `yaml:"progs,omitempty"`         // mob progs. key is an event ("greet", "enter", "death"...) and the value is motherfucking javascript.
	Flags     []string             `yaml:"flags,omitempty"`         // list of flags. See [entity_flags] for values.
	AI        Brain                `yaml:"-"`                       // actual AI interface. instantiated upon spawn.
	Attacker  Entity               `yaml:"-"`                       // who is this mob fighting?
}

// Returns true if the entity is a *PlayerProfile, false if just a *CharData mob.
func (*CharData) IsPlayer() bool {
	return false
}

// Entities don't need screens...
func (*CharData) Prompt() {
}

// Entities don't have connections
func (c *CharData) Send(m string, args ...interface{}) {
}

// Room ID of the entity
func (c *CharData) RoomId() uint {
	return c.Room
}

// Ship ID of the entity. If 0, entity is on planet.
func (c *CharData) ShipId() uint {
	return c.Ship
}

// Triggers a mob prog to fire.
func (c *CharData) Event(evt string) {
	e := mud_prog_exec(c.Progs[evt], c, evt)
	ErrorCheck(e)
}

// Current HP
func (c *CharData) CurrentHp() int {
	return c.Hp[0]
}

// Maximum HP
func (c *CharData) MaxHp() int {
	return c.Hp[1]
}

// Current MV
func (c *CharData) CurrentMv() int {
	return c.Mv[0]
}

// Maximum MV
func (c *CharData) MaxMv() int {
	return c.Mv[1]
}

// Base weight of a person based on race (ignoring gender for sake of gender equality and body positivity ;)
func (c *CharData) base_weight() int {
	switch c.Race {
	case "Wookiee": // big guy eats steaks for sure
		return 105
	case "Hutt": // because there's no fatter
		return 425
	case "Ewok":
	case "Jawa": // the tinyist of hero's
		return 25
	case "Droid": // metal weighs more than bone, facts.
	case "Protocol Droid":
		return 95
	case "Assassin Droid": // so do guns...
	case "Gladiator Droid":
		return 245
	}
	return 75
}

// Calculates the current weight of the mob taking into account their inventory (equipment isn't factored, yet...)
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

	// str / 10 * base_weight + (level * 5) + dex / 10 * base_weight  : because math is awesome.
	return ((c.Stats[0] / 10) * weight) + int(c.Level*5) + ((c.Stats[2] / 10) * weight)
}

// Total number of held objects (objects in containers, don't count, that's the point of containers...)
func (c *CharData) CurrentInventoryCount() int {
	count := 0
	for range c.Inventory {
		count++
	}
	return count
}

// How many items can you juggle on your person?
func (c *CharData) MaxInventoryCount() int {
	return (int(c.Level) * 3) + c.Stats[0]
}

// Is fighting someone?
func (c *CharData) IsFighting() bool {
	return c.State == ENTITY_STATE_FIGHTING
}

// Stop fighting
func (c *CharData) StopFighting() {
	if c.State == ENTITY_STATE_FIGHTING {
		c.State = ENTITY_STATE_NORMAL
		c.Attacker = nil
	}
}

// Start fighting someone
func (c *CharData) SetAttacker(entity Entity) {
	c.Attacker = entity
	c.State = ENTITY_STATE_FIGHTING
}

// What's your armor class? AC can't be above 20.
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

// How hard did you hit for your skill and weapon?
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
	skill = umin(1, skill/10)
	dmg := uint(roll_dice(d)) + uint(roll_dice(fmt.Sprintf("%dd4", skill)))
	return dmg
}

// Get's the entity's weapon. If nil, entity is fighting bare handed like chuck norris.
func (c *CharData) Weapon() Item {
	if i, ok := c.Equipment["weapon"]; ok {
		item := i
		return item
	}
	return nil
}

// Find an item on this entity by keyword. Useful for checking existence of keys or player commands on items.
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

// Remove an item from this person. No soup for you.
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

// Get's an item from this entity based on item id (or item type id), not that dissimilar to find, only you know the id.
func (c *CharData) GetItem(item_id uint) Item {
	for id := range c.Inventory {
		i := c.Inventory[id]
		if i.GetId() == item_id || i.GetData().OId == item_id {
			return i
		}
	}
	return nil
}

// Apply damage and play dead
func (c *CharData) ApplyDamage(damage uint) {
	c.Hp[0] -= int(damage)
	if c.Hp[0] <= 0 {
		c.State = ENTITY_STATE_DEAD
		if c.AI != nil {
			c.AI.OnDeath()
		}
	}
}

// Get the ship this entity is onboard. Physically. If nil, player is on planet (or other static location)
func (c *CharData) GetShip() Ship {
	if c.Ship > 0 {
		return DB().GetShip(c.Ship)
	}
	return nil
}

func (c *CharData) GetRoom() *RoomData {
	return DB().GetRoom(c.Room, c.Ship)
}

// Return the backing CharData struct
func (c *CharData) GetCharData() *CharData {
	return c
}

// [PlayerProfile] is an [Entity] that represents the player, not a mob. As such it has a few extra fields...
// [Entity.IsPlayer] will return whether or not an [Entity] is a [*PlayerProfile] or just [*CharData]
type PlayerProfile struct {
	Char        CharData  `yaml:"char,inline"`
	Email       string    `yaml:"email,omitempty"`
	Password    string    `yaml:"password,omitempty"`
	Priv        int       `yaml:"priv,omitempty"`
	LastSeen    time.Time `yaml:"last_seen,omitempty"`
	Banned      bool      `yaml:"banned,omitempty"`
	Frequency   string    `yaml:"freq"`
	Kills       uint      `yaml:"kills"`
	PKills      uint      `yaml:"pkills"`
	Client      Client    `yaml:"-"`
	NeedPrompt  bool      `yaml:"-"`
	LastCommand string    `yaml:"-"`
}

// Is Entity a player?
func (*PlayerProfile) IsPlayer() bool {
	return true
}

// Send player a message through their [Client]
func (p *PlayerProfile) Send(m string, any ...interface{}) {
	if p.Client != nil {
		p.Client.Sendf(m, any...)
		p.NeedPrompt = true
	}
}

// Get's the room id of the player.
func (p *PlayerProfile) RoomId() uint {
	return p.Char.Room
}

// Get's the ship id of the player. Physically. Not the one they own.
func (p *PlayerProfile) ShipId() uint {
	return p.Char.Ship
}

// Trigger an event on a player. Since players are in control of their own characters, not sure how useful this is.
func (p *PlayerProfile) Event(evt string) {
}

// Player's current hitpoints
func (p *PlayerProfile) CurrentHp() int {
	return p.Char.Hp[0]
}

// Player's maximum hitpoints
func (p *PlayerProfile) MaxHp() int {
	return p.Char.Hp[1]
}

// Player's current movement
func (p *PlayerProfile) CurrentMv() int {
	return p.Char.Mv[0]
}

// Player's maximum movement
func (p *PlayerProfile) MaxMv() int {
	return p.Char.Mv[1]
}

// Prompt the player, show's their stats, and readies for input next turn.
func (p *PlayerProfile) Prompt() {
	if p.NeedPrompt && p.Char.State != ENTITY_STATE_DEAD && !p.Client.(*MudClient).Editing {
		prompt := player_prompt(p)
		p.Send("%s\r\n", prompt)
		p.NeedPrompt = false
	}
}

// Is the player fighting something?
func (p *PlayerProfile) IsFighting() bool {
	return p.Char.State == ENTITY_STATE_FIGHTING
}

// Stop fighthing whatever it's fighting.
func (p *PlayerProfile) StopFighting() {
	if p.Char.State == ENTITY_STATE_FIGHTING {
		p.Char.State = ENTITY_STATE_NORMAL
		p.Char.Attacker = nil
		p.Send("\r\n&dYou stop fighting.\r\n")
	}
}

// Start fighting [Entity], combat will commence next turn.
func (p *PlayerProfile) SetAttacker(entity Entity) {
	p.Char.Attacker = entity
	p.Char.State = ENTITY_STATE_FIGHTING
}

// Get the underlying [CharData] pointer.
func (p *PlayerProfile) GetCharData() *CharData {
	return &p.Char
}

// Apply damange and change state for knockout/death.
func (p *PlayerProfile) ApplyDamage(damage uint) {
	p.Char.Hp[0] -= int(damage)
	if p.Char.Hp[0] <= 0 {
		p.Char.State = ENTITY_STATE_UNCONSCIOUS
		if p.Char.Hp[0] <= -(p.Char.Hp[1] * 2) {
			p.Char.State = ENTITY_STATE_DEAD
			p.Send(sprintf("\r\n&RYou have died.&d %s\r\n", EMOJI_SKULL))
			return
		}
		p.Send("\r\n&YYou have been knocked unconscious...&d\r\n")
		return
	}
}

// Get's the player's current weapon. If nil, player is bruce lee...
func (p *PlayerProfile) Weapon() Item {
	if i, ok := p.Char.Equipment["weapon"]; ok {
		item := i
		return item
	}
	return nil
}

// Find's an item on the player by keyword.
func (p *PlayerProfile) FindItem(keyword string) Item {
	return p.Char.FindItem(keyword)
}

// Get's the ship the player is currently on.
func (p *PlayerProfile) GetShip() Ship {
	if p.Char.Ship > 0 {
		return DB().GetShip(p.Char.Ship)
	}
	return nil
}

func (p *PlayerProfile) GetRoom() *RoomData {
	return DB().GetRoom(p.Char.Room, p.Char.Ship)
}

// Clones an entity, generating a new ID, copying over the values, and returns the cloned [Entity]
func entity_clone(entity Entity) Entity {
	ch := entity.GetCharData()
	c := &CharData{
		Id:        gen_npc_char_id(),
		OId:       ch.Id,
		Room:      ch.Room,
		Name:      ch.Name,
		Filename:  ch.Filename,
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
		Progs:     make(map[string]string),
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
	for k, v := range ch.Progs {
		c.Progs[k] = v
	}
	return c
}

// Builds a player prompt to send to the player using pretty ANSI colors and ASCII glyphs.
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

// Called every turn to process the hundreds of entities in the game. Processes state affects as well as health, movement, force regen.
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

// Heals an entity naturally. Called by processEntities and should never be called by a skill/spell/script.
func processHealing(entity Entity) {
	ch := entity.GetCharData()
	if ch.State == ENTITY_STATE_DEAD {
		return
	}
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
			if roll_dice("1d10") >= 10-entity_get_skill_value(ch, "healing") {
				entity_add_skill_value(entity, "healing", 1)
			}
			entity.Prompt()
		}
	}
	entity.Prompt()
}

// Adds XP to the entity. If the entity gains a level, calculate new HP/MP/MV maximums.
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
	ch.Level = get_level_for_xp(ch.XP) + 1 // get_level_for_xp is 0 based (0-99)
	if ch.Level > 100 {
		ch.XP = get_xp_for_level(100)
		ch.Level = 100
	} else {
		if ch.Level != level {
			entity_print_level_up(entity)
		}
	}
	entity.Send("\r\n&dYou gained &w%d&d xp.\r\n", xp)
}

func entity_print_level_up(entity Entity) {
	ch := entity.GetCharData()
	entity.Send("\r\n}YYou have gained a level!&d\r\n")
	entity.Send("\r\n&YYou are now level &W%d&d.\r\n", ch.Level)
	// reset current life stats as a reward and gain a little extra
	ch.Hp[1] = 50 + (int(ch.Level) / 2)
	ch.Mv[1] = 50 + (int(ch.Level) / 2)
	if ch.Mp[1] > 0 { // you don't get force unless you got force
		ch.Mp[1] = 50 + (int(ch.Level) / 2)
	}
	ch.Hp[0] = ch.Hp[1]
	ch.Mp[0] = ch.Mp[1]
	ch.Mv[0] = ch.Mv[1]
}

func entity_advance_level(entity Entity) {
	ch := entity.GetCharData()
	ch.Level++
	ch.XP = get_xp_for_level(ch.Level)
	entity_print_level_up(entity)
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

// Returns the level for a given XP amount. !WARNING! 0 based. Level 1 is really 0 so +1 to the return if you're printing this.
func get_level_for_xp(xp uint) uint {
	return uint(math.Sqrt(float64(xp) / 500))
}

// What's the bottom line XP for a level. level+1 return value is target XP for next level.
func get_xp_for_level(level uint) uint {
	return uint(math.Pow(float64(level), 2)) * 500
}

// Can the entity speak (or see? (or breathe?)). Returns true if entity is in an unspeakable state. By passes messages sent to them for some actions.
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

// Why can't they speak? (or see? (or breathe?)). Returns the reason an entity is in an unspeakable state.
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

// Pick up an item off the ground or off a corpse (living or dead). Protects against picking up corpses or objects too heavy to lift.
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

// Returns a 0-100 skill value for a skill.
func entity_get_skill_value(ch *CharData, skill string) int {
	if v, ok := ch.Skills[strings.ToLower(skill)]; ok {
		return v
	}
	return 0
}

// Adds a skill value for a skill for supplied [Entity]
func entity_add_skill_value(entity Entity, skill string, value int) {
	ch := entity.GetCharData()
	if ch == nil {
		return
	}
	ch.Skills[skill] += value
	if ch.Skills[skill] >= 100 {
		ch.Skills[skill] = 100
	} else {
		entity.Send("\r\n&CYou gain some knowledge of %s.&d\r\n", skill)
	}

}

// Returns the name of the item in the equipment slot. Useful for SCORE.
func entity_get_equipment_for_slot(entity Entity, wearLoc string) string {
	if o, ok := entity.GetCharData().Equipment[wearLoc]; ok {
		return o.Name
	}
	return "None"
}

// Award a kill tally for the player's kill. If the victim is a player, award a PKill.
// This will only be called upon death. Not unconscious.
func entity_award_kill(killer Entity, victim Entity) {
	if killer.IsPlayer() {
		kp := killer.(*PlayerProfile)
		if victim.IsPlayer() {
			kp.PKills++
		} else {
			kp.Kills++
		}
	}
}
