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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

func FileExists(filename string) bool {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("%s does not exist!", filename)
		return false
	}
	defer file.Close()
	return true
}

var _db *GameDatabase

type HelpData struct {
	Name     string   `yaml:"name"`
	Keywords []string `yaml:"keywords,flow"`
	Desc     string   `yaml:"desc"`
	Level    uint     `yaml:"level"`
}

type GameDatabase struct {
	m               *sync.Mutex
	clients         []*MudClient
	entities        []Entity
	areas           map[string]*AreaData // pointers to the [AreaData] of the game.
	rooms           map[uint]*RoomData   // pointers to the room structs in [AreaData]
	mobs            map[string]*CharData // used as templates for spawning [entities]
	items           map[uint]*ItemData   // used as templates for spawning [items]
	ships           []Ship
	ship_prototypes map[string]*ShipData // used as templates for spawning [ships]
	starsystems     []Starsystem         // Planets (star systems)
	helps           []*HelpData
}

func DB() *GameDatabase {
	if _db == nil {
		_db = new(GameDatabase)
		_db.m = &sync.Mutex{}
		_db.clients = make([]*MudClient, 0, 64)
		_db.entities = make([]Entity, 0)
		_db.areas = make(map[string]*AreaData)
		_db.rooms = make(map[uint]*RoomData)
		_db.mobs = make(map[string]*CharData)
		_db.items = make(map[uint]*ItemData)
		_db.ships = make([]Ship, 0)
		_db.ship_prototypes = make(map[string]*ShipData)
		_db.starsystems = make([]Starsystem, 0)
		_db.helps = make([]*HelpData, 0)
	}
	return _db
}

func (d *GameDatabase) Lock() {
	d.m.Lock()
}

func (d *GameDatabase) Unlock() {
	d.m.Unlock()
}

func (d *GameDatabase) AddClient(client *MudClient) {
	d.Lock()
	defer d.Unlock()
	d.clients = append(d.clients, client)
}

func (d *GameDatabase) RemoveClient(client *MudClient) {
	d.Lock()
	defer d.Unlock()
	index := -1
	for i, c := range d.clients {
		if c.Id == client.Id {
			index = i
		}
	}
	if index > -1 {
		ret := make([]*MudClient, len(d.clients)-1)
		ret = append(ret, d.clients[:index]...)
		ret = append(ret, d.clients[index+1:]...)
		d.clients = ret
	}
}

func (d *GameDatabase) RemoveEntity(entity Entity) {
	d.Lock()
	defer d.Unlock()
	index := -1
	for i, e := range d.entities {
		if e == entity {
			index = i
		}
	}
	if index > -1 {
		ret := make([]Entity, len(d.entities)-1)
		ret = append(ret, d.entities[:index]...)
		ret = append(ret, d.entities[index+1:]...)
		d.entities = ret
	}
}

// The Mother of all load functions
func (d *GameDatabase) Load() {
	d.Lock()
	defer d.Unlock()

	// Load Help files
	d.LoadHelps()

	// Load Areas
	d.LoadAreas()

	// Load Items
	d.LoadItems()

	// Load Planets / Star Systems
	d.LoadPlanets()

	// Load Mobs
	d.LoadMobs()

	// Load Progs
	d.LoadMudProgs()
}

func (d *GameDatabase) LoadHelps() {
	flist, err := ioutil.ReadDir("docs")
	ErrorCheck(err)
	for _, help_file := range flist {
		fpath := fmt.Sprintf("docs/%s", help_file.Name())
		fp, err := ioutil.ReadFile(fpath)
		ErrorCheck(err)
		help := new(HelpData)
		err = yaml.Unmarshal(fp, help)
		ErrorCheck(err)
		d.helps = append(d.helps, help)
	}
	log.Printf("%d help files loaded.\n", len(flist))
}

func (d *GameDatabase) LoadAreas() {
	flist, err := ioutil.ReadDir("data/areas")
	ErrorCheck(err)
	count := 0
	for _, area_file := range flist {
		if strings.HasSuffix(area_file.Name(), "yml") {
			d.LoadArea(area_file.Name())
			count++
		}
	}
	log.Printf("%d areas loaded.\n", count)
}

func (d *GameDatabase) LoadArea(name string) {
	fpath := fmt.Sprintf("data/areas/%s", name)
	fp, err := ioutil.ReadFile(fpath)
	ErrorCheck(err)
	area := new(AreaData)
	err = yaml.Unmarshal(fp, area)
	ErrorCheck(err)
	for i := range area.Rooms {
		room := area.Rooms[i]
		room.Area = area
		d.rooms[room.Id] = &room
		time.Sleep(1 * time.Millisecond)
	}
	d.areas[area.Name] = area
}

func (d *GameDatabase) LoadPlanets() {
	flist, err := ioutil.ReadDir("data/planets")
	ErrorCheck(err)
	for _, f := range flist {
		fpath := fmt.Sprintf("data/planets/%s", f.Name())
		fp, err := ioutil.ReadFile(fpath)
		ErrorCheck(err)
		p := new(StarSystemData)
		err = yaml.Unmarshal(fp, p)
		ErrorCheck(err)
		d.starsystems = append(d.starsystems, p)
	}
}

func (d *GameDatabase) LoadItems() {
	flist, err := ioutil.ReadDir("data/items")
	ErrorCheck(err)
	for _, f := range flist {
		fpath := fmt.Sprintf("data/items/%s", f.Name())
		fp, err := ioutil.ReadFile(fpath)
		ErrorCheck(err)
		item := new(ItemData)
		err = yaml.Unmarshal(fp, item)
		ErrorCheck(err)
		d.items[item.Id] = item
	}
}

func (d *GameDatabase) LoadMobs() {
	flist, err := ioutil.ReadDir("data/mobs")
	ErrorCheck(err)
	for _, f := range flist {
		fpath := fmt.Sprintf("data/mobs/%s", f.Name())
		fp, err := ioutil.ReadFile(fpath)
		ErrorCheck(err)
		ch := new(CharData)
		err = yaml.Unmarshal(fp, ch)
		ErrorCheck(err)
		d.mobs[ch.Name] = ch
	}
}

func (d *GameDatabase) LoadMudProgs() {

}

// The Mother of all save functions
func (d *GameDatabase) Save() {
	d.Lock()
	defer d.Unlock()
	d.SaveAreas()
}

func (d *GameDatabase) SaveAreas() {
	for _, area := range d.areas {
		d.SaveArea(area)
	}
}

func (d *GameDatabase) SaveArea(area *AreaData) {
	buf, err := yaml.Marshal(area)
	ErrorCheck(err)
	err = ioutil.WriteFile(fmt.Sprintf("data/areas/%s.yml", area.Name), buf, 0755)
	ErrorCheck(err)
}

func (d *GameDatabase) ReadPlayerData(filename string) *PlayerProfile {
	fp, err := ioutil.ReadFile(filename)
	ErrorCheck(err)
	p_data := new(PlayerProfile)
	yaml.Unmarshal(fp, p_data)
	return p_data
}

func (d *GameDatabase) SavePlayerData(player *PlayerProfile) {
	name := strings.ToLower(player.Char.Name)
	filename := fmt.Sprintf("data/accounts/%s/%s.yml", name[0:1], name)
	buf, err := yaml.Marshal(player)
	ErrorCheck(err)
	err = ioutil.WriteFile(filename, buf, 0755)
	ErrorCheck(err)
}
func (d *GameDatabase) ReadCharData(filename string) *CharData {
	fp, err := ioutil.ReadFile(filename)
	ErrorCheck(err)
	char_data := new(CharData)
	yaml.Unmarshal(fp, char_data)
	return char_data
}

func (d *GameDatabase) SaveCharData(char_data *CharData, filename string) {
	buf, err := yaml.Marshal(char_data)
	ErrorCheck(err)
	err = ioutil.WriteFile(filename, buf, 0755)
	ErrorCheck(err)
}

func (d *GameDatabase) AddEntity(entity Entity) {
	d.Lock()
	defer d.Unlock()
	d.entities = append(d.entities, entity)
}

func (d *GameDatabase) SpawnEntity(mobName string) Entity {
	d.Lock()
	defer d.Unlock()
	e := entity_clone(d.mobs[mobName])
	d.entities = append(d.entities, e)
	return e
}

func (d *GameDatabase) RemoveCorpse(entity Entity) {
	for _, e := range d.GetEntitiesInRoom(entity.RoomId()) {
		if e != entity {
			e.Send("\r\n&xThe corpse of %s wastes away...&d\r\n", entity.GetCharData().Name)
		}
	}
	d.RemoveEntity(entity)
}

func (d *GameDatabase) GetEntitiesInRoom(roomId uint) []Entity {
	d.Lock()
	defer d.Unlock()
	ret := make([]Entity, 0)
	for _, entity := range d.entities {
		if entity.RoomId() == roomId {
			ret = append(ret, entity)
		}
	}
	return ret
}

func (d *GameDatabase) GetRoom(roomId uint) *RoomData {
	d.Lock()
	defer d.Unlock()
	for _, r := range d.rooms {
		if r.Id == roomId {
			return r
		}
	}
	return nil
}

func (d *GameDatabase) GetEntityForClient(client Client) Entity {
	for _, e := range d.entities {
		if e.IsPlayer() {
			player := e.(*PlayerProfile)
			if player.Client == client {
				return player
			}
		}
	}
	return nil
}

func (d *GameDatabase) GetHelp(help string) []*HelpData {
	ret := []*HelpData{}
	for _, h := range d.helps {
		for _, keyword := range h.Keywords {
			if len(keyword) < len(help) {
				continue
			}
			match := true
			for i, r := range help {
				if keyword[i] != byte(r) {
					match = false
				}
			}
			if match {
				ret = append(ret, h)
			}
		}
	}
	return ret
}
