package swr

import "sync"

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
