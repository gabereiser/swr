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
	"os"
	"sync"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

var _db *GameDatabase

type GameDatabase struct {
	m       *sync.Mutex
	clients []*MudClient
}

func DB() *GameDatabase {
	if _db == nil {
		_db = &GameDatabase{
			m:       &sync.Mutex{},
			clients: make([]*MudClient, 0, 64),
		}
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

}

// The Mother of all save functions
func (d *GameDatabase) Save() {

}

func (d *GameDatabase) ReadCharData(filename string) *CharData {
	return nil
}
