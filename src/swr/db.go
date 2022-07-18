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
	m        *sync.Mutex
	clients  []*MudClient
	entities []Entity
	areas    []AreaData
	rooms    []RoomData
	helps    []HelpData
}

func DB() *GameDatabase {
	if _db == nil {
		_db = new(GameDatabase)
		_db.m = &sync.Mutex{}
		_db.clients = make([]*MudClient, 0, 64)
		_db.entities = make([]Entity, 0)
		_db.areas = make([]AreaData, 0)
		_db.rooms = make([]RoomData, 0)
	}
	return _db
}

func (d *GameDatabase) Lock() {
	d.m.Lock()
}

func (d *GameDatabase) Unlock() {
	d.m.Unlock()
}

func (d *GameDatabase) RemoveIndex(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
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
		d.helps = append(d.helps, *help)
	}
	log.Printf("%d help files loaded.\n", len(flist))
}

func (d *GameDatabase) LoadAreas() {
	flist, err := ioutil.ReadDir("data/areas")
	ErrorCheck(err)
	for _, area_file := range flist {
		d.LoadArea(area_file.Name())
	}
	log.Printf("%d areas loaded.\n", len(flist))
}

func (d *GameDatabase) LoadArea(name string) {
	fpath := fmt.Sprintf("data/areas/%s", name)
	fp, err := ioutil.ReadFile(fpath)
	ErrorCheck(err)
	area := new(AreaData)
	err = yaml.Unmarshal(fp, area)
	ErrorCheck(err)
	for vnum, r := range area.Rooms {
		r.Id = uint(vnum)
		d.rooms = append(d.rooms, r)
	}
	d.areas = append(d.areas, *area)
}

func (d *GameDatabase) LoadItems() {

}

func (d *GameDatabase) LoadMobs() {

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
		d.SaveArea(&area)
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
	name := strings.ToLower(player.Name())
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
	for _, a := range d.areas {
		for vnum, r := range a.Rooms {
			if uint(vnum) == roomId {
				return &r
			}
		}
	}

	panic(fmt.Sprintf("RoomId %d not found!", roomId))

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

func (d *GameDatabase) GetHelp(help string) []HelpData {
	ret := []HelpData{}
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
