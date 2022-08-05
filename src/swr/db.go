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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	mobs            map[uint]*CharData   // used as templates for spawning [entities]
	items           map[uint]*ItemData   // used as templates for spawning [items]
	ships           []Ship
	ship_prototypes map[uint]*ShipData // used as templates for spawning [ships]
	starsystems     []Starsystem       // Planets (star systems)
	helps           []*HelpData
}

func DB() *GameDatabase {
	if _db == nil {
		log.Printf("Starting Database.")
		_db = new(GameDatabase)
		_db.m = &sync.Mutex{}
		_db.clients = make([]*MudClient, 0, 64)
		_db.entities = make([]Entity, 0)
		_db.areas = make(map[string]*AreaData)
		_db.rooms = make(map[uint]*RoomData)
		_db.mobs = make(map[uint]*CharData)
		_db.items = make(map[uint]*ItemData)
		_db.ships = make([]Ship, 0)
		_db.ship_prototypes = make(map[uint]*ShipData)
		_db.starsystems = make([]Starsystem, 0)
		_db.helps = make([]*HelpData, 0)
		log.Printf("Database Started.")
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
	d.clients = append(d.clients, client)
}

func (d *GameDatabase) RemoveClient(client *MudClient) {
	for _, e := range d.entities {
		if e == nil {
			continue
		}
		if e.IsPlayer() {
			p := e.(*PlayerProfile)
			if p.Client != nil {
				if p.Client == client {
					d.RemoveEntity(e)
				}
			}
		}
	}
	index := -1
	for i, c := range d.clients {
		if c == nil {
			continue
		}
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
	if entity == nil {
		return
	}
	index := -1
	for i, e := range d.entities {
		if e == nil {
			continue
		}
		if e == entity {
			index = i
		}
	}
	if index > -1 {
		ret := make([]Entity, len(d.entities)-1)
		ret = append(ret, d.entities[:index]...)
		ret = append(ret, d.entities[index+1:]...)
		d.entities = ret
	} else {
		ErrorCheck(Err(fmt.Sprintf("Can't find entity %s to remove.", entity.GetCharData().Name)))
	}
}
func (d *GameDatabase) RemoveShip(ship Ship) {
	index := -1
	for i, s := range d.ships {
		if s == nil {
			continue
		}
		if s == ship {
			index = i
		}
	}
	if index > -1 {
		ret := make([]Ship, len(d.ships)-1)
		ret = append(ret, d.ships[:index]...)
		ret = append(ret, d.ships[index+1:]...)
		d.ships = ret
	} else {
		ErrorCheck(Err(fmt.Sprintf("Can't find ship %s", ship.GetData().Name)))
	}
}

func (d *GameDatabase) RemoveArea(area *AreaData) {
	if a, ok := d.areas[area.Name]; ok {
		for _, r := range a.Rooms {
			delete(d.rooms, r.Id)
		}
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
	err := filepath.Walk("data/items",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				d.LoadItem(path)
			}
			return nil
		})
	ErrorCheck(err)
	log.Printf("%d items loaded.", len(d.items))
}

func (d *GameDatabase) LoadItem(path string) {
	fp, err := ioutil.ReadFile(path)
	ErrorCheck(err)
	item := new(ItemData)
	err = yaml.Unmarshal(fp, item)
	ErrorCheck(err)
	item.Filename = path
	d.items[item.Id] = item
}

func (d *GameDatabase) LoadMobs() {
	err := filepath.Walk("data/mobs",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fp, err := ioutil.ReadFile(path)
				ErrorCheck(err)
				ch := new(CharData)
				err = yaml.Unmarshal(fp, ch)
				ErrorCheck(err)
				ch.Filename = path
				d.mobs[ch.Id] = ch
			}
			return nil
		})
	ErrorCheck(err)
	log.Printf("%d mobs loaded.", len(d.mobs))
}

// The Mother of all save functions
func (d *GameDatabase) Save() {
	d.Lock()
	defer d.Unlock()
	d.SaveAreas()
	d.SaveItems()
	d.SaveShips()
}

func (d *GameDatabase) SaveAreas() {
	for _, area := range d.areas {
		d.SaveArea(area)
	}
}
func (d *GameDatabase) SaveItems() {
	for _, item := range d.items {
		d.SaveItem(item)
	}
}

func (d *GameDatabase) SaveShips() {
	for _, ship := range d.ship_prototypes {
		buf, err := yaml.Marshal(ship)
		ErrorCheck(err)
		err = ioutil.WriteFile(fmt.Sprintf("data/ships/prototypes/%s.yml", strings.ToLower(strings.ReplaceAll(ship.Name, " ", ""))), buf, 0755)
		ErrorCheck(err)
	}
	for _, ship := range d.ships {
		d.SaveShip(ship)
	}
}

func (d *GameDatabase) SaveShip(ship Ship) {
	buf, err := yaml.Marshal(ship)
	ErrorCheck(err)
	err = ioutil.WriteFile(fmt.Sprintf("data/ships/%s.yml", strings.ToLower(strings.ReplaceAll(ship.GetData().Name, " ", ""))), buf, 0755)
	ErrorCheck(err)
}

func (d *GameDatabase) SaveArea(area *AreaData) {
	buf, err := yaml.Marshal(area)
	ErrorCheck(err)
	err = ioutil.WriteFile(fmt.Sprintf("data/areas/%s.yml", area.Name), buf, 0755)
	ErrorCheck(err)
}
func (d *GameDatabase) SaveItem(item *ItemData) {
	buf, err := yaml.Marshal(item)
	ErrorCheck(err)
	err = ioutil.WriteFile(item.Filename, buf, 0755)
	ErrorCheck(err)
}

func (d *GameDatabase) GetPlayer(name string) *PlayerProfile {
	d.Lock()
	defer d.Unlock()
	var player *PlayerProfile
	for i := range d.entities {
		e := d.entities[i]
		if e != nil {
			if e.IsPlayer() {
				if e.GetCharData().Name == name {
					player = e.(*PlayerProfile)
				}
			}
		}
	}
	// Player isn't online
	if player == nil {
		path := fmt.Sprintf("data/accounts/%s/%s.yml", strings.ToLower(name[0:1]), strings.ToLower(name))
		player = d.ReadPlayerData(path)
	}
	return player
}

func (d *GameDatabase) ReadPlayerData(filename string) *PlayerProfile {
	fp, err := ioutil.ReadFile(filename)
	ErrorCheck(err)
	p_data := new(PlayerProfile)
	err = yaml.Unmarshal(fp, p_data)
	if err != nil {
		ErrorCheck(err)
		return nil
	}
	if p_data.Char.Equipment == nil {
		p_data.Char.Equipment = make(map[string]*ItemData)
	}
	if p_data.Char.Inventory == nil {
		p_data.Char.Inventory = make([]*ItemData, 0)
	}
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

func (d *GameDatabase) GetPlayerEntityByName(name string) Entity {
	for _, e := range d.entities {
		if e == nil {
			continue
		}
		if e.IsPlayer() {
			p := e.(*PlayerProfile)
			if strings.EqualFold(p.Char.Name, name) {
				return p
			}
		}
	}
	return nil
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

func (d *GameDatabase) AddShip(ship Ship) {
	d.Lock()
	defer d.Unlock()
	d.ships = append(d.ships, ship)
}

func (d *GameDatabase) SpawnEntity(entity Entity) Entity {
	e := entity_clone(entity)
	e.GetCharData().State = ENTITY_STATE_NORMAL
	if e.GetCharData().AI == nil {
		e.GetCharData().AI = MakeGenericBrain(e)
		e.GetCharData().AI.OnSpawn()
	}
	d.AddEntity(e)
	return e
}

func (d *GameDatabase) SpawnShip(ship Ship) Ship {
	s := ship_clone(ship)
	d.AddShip(s)
	return s
}

func (d *GameDatabase) GetShip(shipId uint) Ship {
	for _, ship := range d.ships {
		if ship.GetData().Id == shipId {
			return ship
		}
	}
	return nil
}
func (d *GameDatabase) GetShipsInSystem(system string) []Ship {
	ret := make([]Ship, 0)
	for _, ship := range d.ships {
		s := ship.GetData()
		if s.CurrentSystem == system {
			if s.InSpace {
				ret = append(ret, ship)
			}
		}
	}
	return ret
}
func (d *GameDatabase) GetShipsInRoom(roomId uint) []Ship {
	ret := make([]Ship, 0)
	for _, s := range d.ships {
		if s.GetData().LocationId == roomId && !s.GetData().InSpace {
			ret = append(ret, s)
		}
	}
	return ret
}
func (d *GameDatabase) GetEntity(entity Entity) Entity {
	for _, e := range d.entities {
		if e == entity {
			return e
		}
	}
	return nil
}
func (d *GameDatabase) GetEntitiesInRoom(roomId uint, shipId uint) []Entity {
	ret := make([]Entity, 0)
	for _, entity := range d.entities {
		if entity == nil {
			continue
		}
		if entity.RoomId() == roomId && entity.ShipId() == shipId {
			ret = append(ret, entity)
		}
	}
	return ret
}

func (d *GameDatabase) GetRoom(roomId uint, shipId uint) *RoomData {
	if shipId > 0 {
		for _, s := range d.ships {
			if s.GetData().Id == shipId {
				for _, r := range s.GetData().Rooms {
					if r.Id == roomId {
						return r
					}
				}
			}
		}
	} else {
		for _, r := range d.rooms {
			if r == nil {
				continue
			}
			if r.Id == roomId {
				return r
			}
		}
	}
	return nil
}

func (d *GameDatabase) GetNextRoomVnum(roomId uint, shipId uint) uint {
	if shipId > 0 {
		for _, s := range d.ships {
			if s.GetData().Id == shipId {
				lastVnum := uint(0)
				for _, r := range s.GetData().Rooms {
					if r.Name == "A void" { // return the first void prototype room...
						return r.Id
					}
					if r.Id > lastVnum {
						lastVnum = r.Id
					}
				}
				return lastVnum + 1
			}
		}
	} else {
		lastVnum := uint(0)
		// let's get the room's area
		room := d.GetRoom(roomId, shipId)
		for _, r := range room.Area.Rooms {
			if r.Name == "A void" { // return the first void prototype room
				return r.Id
			}
			if r.Id > lastVnum {
				lastVnum = r.Id
			}
		}
		return lastVnum + 1
	}
	return 0
}
func (d *GameDatabase) GetNextItemVnum() uint {
	for i := uint(1); i < (^uint(0)); i++ {
		if _, ok := d.items[i]; !ok {
			return i
		}
	}
	return 0
}

func (d *GameDatabase) GetItem(itemId uint) Item {
	for _, i := range d.items {
		if i == nil {
			continue
		}
		if i.GetId() == itemId {
			return i
		}
	}
	return nil
}

func (d *GameDatabase) GetMob(mobId uint) Entity {
	if m, ok := d.mobs[mobId]; ok {
		return m
	}
	return nil
}

func (d *GameDatabase) GetEntityForClient(client Client) Entity {
	for _, e := range d.entities {
		if e == nil {
			continue
		}
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

func (d *GameDatabase) ResetAll() {
	for area_name, area := range d.areas {
		log.Printf("Resetting Area %s", area_name)
		area_reset(area)
	}
}
