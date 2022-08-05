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
)

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
	Name       string    `yaml:"name"`
	Type       string    `yaml:"type"`
	Radius     uint      `yaml:"radius"`
	Position   []float32 `yaml:"position,flow"`
	Spaceports []uint16  `yaml:"spaceports,flow,omitempty"`
}

const (
	SHIP_MODULE_HYPERDRIVE  = "hyperdrive"
	SHIP_MODULE_CARGO       = "cargo"
	SHIP_MODULE_SHIELD      = "shield"
	SHIP_MODULE_DOCKING_BAY = "dock"
	SHIP_MODULE_AI          = "ai"
	SHIP_MODULE_RADAR       = "radar"
	SHIP_MODULE_ENGINE      = "engine"
	SHIP_MODULE_ENGINE_2    = "engine2"
	SHIP_MODULE_ENGINE_3    = "engine3"
	SHIP_MODULE_ENGINE_4    = "engine4"
	SHIP_MODULE_TURRET      = "turret"
	SHIP_MODULE_TURRET_2    = "turret2"
	SHIP_MODULE_TURRET_3    = "turret3"
	SHIP_MODULE_TURRET_4    = "turret4"
	SHIP_MODULE_WEAPON      = "weapon"
	SHIP_MODULE_WEAPON_2    = "weapon2"
	SHIP_MODULE_WEAPON_3    = "weapon3"
	SHIP_MODULE_WEAPON_4    = "weapon4"
)
const (
	SHIP_ROOM_FLAGS_COCKPIT    = "cockpit"
	SHIP_ROOM_FLAGS_ENGINEROOM = "engineroom"
	SHIP_ROOM_FLAGS_RAMP       = "ramp"
	SHIP_ROOM_FLAGS_TURRET     = "turret"
	SHIP_ROOM_FLAGS_CARGO      = "cargo"
)

type ShipData struct {
	Id               uint               `yaml:"id"`     // instance id when loaded into memory
	OId              uint               `yaml:"shipId"` // global unique type id
	Name             string             `yaml:"name"`
	Desc             string             `yaml:"desc"`
	Type             string             `yaml:"type"`       // class of the ship
	LocationId       uint               `yaml:"locationId"` // roomId where we are docked or where we took off from if in space.
	InSpace          bool               `yaml:"-"`
	CurrentSystem    string             `yaml:"currentSystem,omitempty"` // name of the current system its in
	ShipyardId       uint               `yaml:"shipyardId"`              // where the ship came from.
	Permission       uint               `yaml:"permission"`              // 0 - owner, 1 - group, 2 - guild/clan, 3 - faction, 4 - public
	Simulator        bool               `yaml:"simulator"`               // is it a sim?
	Owner            string             `yaml:"owner,omitempty"`         // who owns this ship?
	Crafter          string             `yaml:"crafter,omitempty"`       // who made this ship?
	Rooms            map[uint]*RoomData `yaml:"rooms"`                   // the ship needs rooms... at the very least a cockpit with a hatch
	Modules          map[string]uint    `yaml:"modules"`                 // ship module healths
	HighSlots        []*ItemData        `yaml:"highSlots,omitempty"`     // loaded ship slots (not all ships have slots)
	LowSlots         []*ItemData        `yaml:"lowSlots,omitempty"`
	Blueprint        uint               `yaml:"blueprintId"`   // the item that an engineer needs to make this ship
	Ramp             uint               `yaml:"ramp"`          // the room that is the ramp (entrance/exit) to the ship
	Cockpit          uint               `yaml:"cockpit"`       // the rooms that are considered cockpits
	EngineRoom       uint               `yaml:"engineRoom"`    // engine rooms (technicians will love these)
	CargoRoom        uint               `yaml:"cargoRoom"`     // storages for kessel runs
	Pilot            Entity             `yaml:"-"`             // who is currently piloting
	CoPilot          Entity             `yaml:"-"`             // who is currently copiloting
	Target           Ship               `yaml:"-"`             // who is this ship targeting
	Position         []float32          `yaml:"position,flow"` // where is this ship in space? InSpace will be true when in space.
	Heading          float32            `yaml:"-"`             // where are we going?
	Speed            float32            `yaml:"-"`             // how fast are we going?
	MaxSpeed         float32            `yaml:"speed"`         // max speed the ship can go (base)
	InHyper          bool               `yaml:"-"`             // are we in hyperspace? (used by the prompt to tell us how long we have to reach our destination)
	HyperDestination Starsystem         `yaml:"-"`             // where are we going in hyperspace?
	HyperOrigin      Starsystem         `yaml:"-"`             // where did we come from?
	HyperTimeUntil   uint               `yaml:"-"`             // time in seconds until we exit hyperspace.
	Hp               []uint             `yaml:"hp,flow"`       // ship hitpoints as an array of uint's. [0] is current hp, [1] is max hp. Always a len() of 2.
	Sp               []uint             `yaml:"sp,flow"`       // ship shield points as an array of uint's. [0] is current sp, [1] is max sp. Always a len() of 2.
}

type Ship interface {
	GetData() *ShipData
}

func (s *ShipData) GetData() *ShipData {
	return s
}

func (s *ShipData) JumpHyperspace(target Starsystem, pos []float32) {
	if s.Pilot != nil || s.CoPilot != nil {
		db := DB()
		s.HyperDestination = target

		for _, star := range db.starsystems {
			if star.GetData().Name == s.CurrentSystem {
				s.HyperOrigin = star
			}
		}
		p1 := s.HyperOrigin.GetData().Position
		p2 := target.GetData().Position
		distance := distance_between_points(p1, p2)
		s.HyperTimeUntil = uint(math.Round(float64(distance)))
		s.Heading = float32(rand_min_max(0, 360))
	}
}

func ship_clone(ship Ship) *ShipData {
	sp := ship.GetData()
	if sp != nil {
		s := new(ShipData)
		s.Id = gen_ship_id()
		s.OId = sp.OId
		s.Name = sp.Name
		s.Desc = sp.Desc
		s.Type = sp.Type
		s.LocationId = sp.LocationId
		s.CurrentSystem = sp.CurrentSystem
		s.ShipyardId = sp.ShipyardId
		s.Permission = sp.Permission
		s.Simulator = sp.Simulator
		s.Owner = sp.Owner
		s.Crafter = sp.Crafter
		s.Rooms = make(map[uint]*RoomData)
		for i, r := range sp.Rooms {
			s.Rooms[i] = r
		}
		s.Modules = make(map[string]uint)
		for i, m := range sp.Modules {
			s.Modules[i] = m
		}
		s.HighSlots = make([]*ItemData, 0)
		s.HighSlots = append(s.HighSlots, sp.HighSlots...)
		s.LowSlots = make([]*ItemData, 0)
		s.LowSlots = append(s.LowSlots, sp.LowSlots...)
		s.Blueprint = sp.Blueprint
		s.Ramp = sp.Ramp
		s.Cockpit = sp.Cockpit
		s.EngineRoom = sp.EngineRoom
		s.CargoRoom = sp.CargoRoom
		s.Pilot = sp.Pilot
		s.CoPilot = sp.CoPilot
		s.Target = sp.Target
		s.Position = make([]float32, 2)
		s.Position[0] = sp.Position[0]
		s.Position[1] = sp.Position[1]
		s.Heading = sp.Heading
		s.Speed = sp.Speed
		s.MaxSpeed = sp.MaxSpeed
		s.HyperOrigin = nil
		s.HyperDestination = nil
		s.HyperTimeUntil = 0
		return s
	}
	return sp
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
		entity.Send(fmt.Sprintf("&Y│ System: &W%-42s &YSector:&d%17s  &Y│\r\n", starsystem.Name, starsystem.Sector))
		p := false
		for _, o := range starsystem.Orbits {
			if !p {
				entity.Send(fmt.Sprintf("&Y│           &W└%-34s    &YPosition: &g%-18s&Y│\r\n", o.Name, sprintf("%4.2f, %4.2f", starsystem.Position[0], starsystem.Position[1])))
				p = true
			} else {
				entity.Send(fmt.Sprintf("&Y│           &W└%-64s  &Y│\r\n", o.Name))
			}

		}
		entity.Send("&Y└──────────────────────────────────────────────────────────────────────────────┘&d\r\n")
	}
	entity.Send("\r\n")
}

func do_board_ship(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	if len(args) == 0 {
		entity.Send("\r\nSyntax: board <shipname>\r\n")
		return
	}
	room := DB().GetRoom(entity.RoomId(), entity.ShipId())
	shipname := strings.Join(args, " ")
	ships := DB().GetShipsInRoom(entity.RoomId())
	for _, s := range ships {
		ship := s.GetData()
		if strings.HasPrefix(ship.Name, shipname) {
			ch := entity.GetCharData()
			if ch.Mv[0] <= 0 {
				entity.Send("\r\n&YYou are exhausted.&d\r\n")
				return
			}
			ch.Room = ship.Ramp
			ch.Ship = ship.Id
			ship_ramp := DB().GetRoom(ship.Ramp, ship.Id)
			ship_ramp.SendToOthers(entity, sprintf("\r\n%s has boarded the ship.\r\n", ch.Name))
			room.SendToOthers(entity, sprintf("\r\n%s left boarding a ship.\r\n", ch.Name))
			entity.Send("\r\nYou board the ship.\r\n")
		}
	}
}

func do_leave_ship(entity Entity, args ...string) {
	if entity == nil {
		return
	}
	room := DB().GetRoom(entity.RoomId(), entity.ShipId())
	ship := DB().GetShip(entity.ShipId())
	if ship.GetData().InSpace {
		entity.Send("\r\n&RYou can't leave the airlock in space.&d\r\n")
		return
	}
	if room.Id != ship.GetData().Ramp {
		entity.Send("\r\n&cPlease make your way to the ramp to leave the ship.&d\r\n")
		return
	} else {
		to_room := DB().GetRoom(ship.GetData().LocationId, 0)
		ch := entity.GetCharData()
		ch.Room = to_room.Id
		ch.Ship = 0
		room.SendToOthers(entity, sprintf("\r\n%s has left the ship.\r\n", ch.Name))
		to_room.SendToOthers(entity, sprintf("\r\n%s has arrived.\r\n", ch.Name))
		entity.Send("\r\nYou leave the ship.")
		ServerQueue <- MudClientCommand{
			Entity:  entity,
			Command: "look",
		}
	}
}
