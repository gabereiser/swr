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

import "fmt"

type StarData struct {
	Name     string    `yaml:"name"`
	Type     string    `yaml:"type"`
	Radius   int       `yaml:"radius"`
	Position []float32 `yaml:"position,flow"`
}
type StarSystemData struct {
	Name     string                 `yaml:"name"`
	Sector   string                 `yaml:"sector"`
	Grid     string                 `yaml:"grid"`
	Position []float32              `yaml:"position,flow"`
	Stars    map[int]StarData       `yaml:"stars"`
	Orbits   map[int]OribitalObject `yaml:"orbits"`
}

type OribitalObject struct {
	Name       string                 `yaml:"name"`
	Type       string                 `yaml:"type"`
	Radius     uint                   `yaml:"radius"`
	Position   []float32              `yaml:"position,flow"`
	Spaceports []uint16               `yaml:"spaceports,flow,omitempty"`
	Market     map[string]interface{} `yaml:"market,omitempty"`
}

type ShipData = map[string]interface{}
type ShipRoom = RoomData

type Ship interface {
	Name() string
	RoomVNum() uint
	InSpace() bool
	GetPosition() []float32
	SetPosition(x float32, y float32)
	GetSpeed() int
	SetSpeed(speed int)
	GetHeading() []float32
	SetHeading(x float32, y float32)
	InHyperspace() bool
	GetRadar() []interface{}
	GetStatus()
	GetOwner() Entity
	GetCrafter() Entity
	GetPilot() Entity
	SetOwner(entity Entity)
	SetCrafter(entity Entity)
	SetPilot(entity Entity)
	GetShields() []uint16
	GetHp() []uint16
	GetShieldState() int
	GetRooms() []*ShipRoom
	GetCockpitRooms() []*ShipRoom
	GetEngineRooms() []*ShipRoom
	GetCargoRooms() []*ShipRoom
	GetRampRooms() []*ShipRoom
	SetRooms(rooms []*ShipRoom)
	SetCockpitRooms(rooms []*ShipRoom)
	SetEngineRooms(rooms []*ShipRoom)
	SetCargoRooms(rooms []*ShipRoom)
	SetRampRooms(rooms []*ShipRoom)
}
type Starsystem interface {
	GetData() *StarSystemData
}

func (s *StarSystemData) GetData() *StarSystemData {
	return s
}

func do_starsystems(entity Entity, args ...string) {
	db := DB()
	entity.Send("\r\n")
	entity.Send(MakeTitle("Star Systems", ANSI_TITLE_STYLE_NORMAL, ANSI_TITLE_ALIGNMENT_CENTER))

	for _, s := range db.starsystems {
		starsystem := s.GetData()
		entity.Send("&Y┌──────────────────────────────────────────────────────────────────────────────┐&d\r\n")
		entity.Send(fmt.Sprintf("&Y│ System: &W%-42s &YSector:&d%16s   &Y│\r\n", starsystem.Name, starsystem.Sector))
		for _, o := range starsystem.Orbits {
			entity.Send(fmt.Sprintf("&Y│           &W└%-34s    &YPosition: &g%4.2f, %4.2f  &Y│\r\n", o.Name, starsystem.Position[0], starsystem.Position[1]))
		}
		entity.Send("&Y└──────────────────────────────────────────────────────────────────────────────┘&d\r\n")
	}
	entity.Send("\r\n")
}
