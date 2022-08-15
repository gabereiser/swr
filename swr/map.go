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

import "strings"

/*
+----------------------+
| @--@--@--@--@--@--@  |
|    |     |           |
|    |     |           |
|    @--@--@           |
*/

const (
	MAP_ROOM      = "@"
	MAP_EXIT_NS   = "|"
	MAP_EXIT_EW   = "-"
	MAP_EXIT_NWSE = "\\"
	MAP_EXIT_SWNE = "/"
)
const MAPSIZE = 10

func build_map(room *RoomData) string {
	m := [MAPSIZE][MAPSIZE]string{}
	for y := 0; y < MAPSIZE; y++ {
		for x := 0; x < MAPSIZE; x++ {
			m[x][y] = " "
		}
	}
	m[MAPSIZE/2][MAPSIZE/2] = "&R@&W"
	m = walk_rooms(m, room, MAPSIZE/2, MAPSIZE/2, 0)
	for y := 0; y < MAPSIZE; y++ {
		for x := 0; x < MAPSIZE; x++ {
			if (x == 0 || x == MAPSIZE-1) && (y == MAPSIZE-1 || y == 0) {
				m[x][y] = "&g+&W"
			} else if (x == 0 || x == MAPSIZE-1) && (y != 0 && y < MAPSIZE-1) {
				m[x][y] = "&g|&W"
			} else if (y == MAPSIZE-1 || y == 0) && (x != 0 && x != MAPSIZE-1) {
				m[x][y] = "&g-&W"
			}
		}
	}
	buf := "&W"
	for y := 0; y < MAPSIZE; y++ {
		for x := 0; x < MAPSIZE; x++ {
			buf += sprintf("&W%s&W", m[x][y])
		}
		buf += "&W\r\n&W"
	}
	return strings.TrimSpace(buf)
}

func map_in_bounds(x int, y int) bool {
	if x > 0 && x < MAPSIZE && y < MAPSIZE && y > 0 {
		return true
	}
	return false
}
func walk_rooms(m [MAPSIZE][MAPSIZE]string, room *RoomData, curX int, curY int, depth int) [MAPSIZE][MAPSIZE]string {
	if map_in_bounds(curX, curY) && depth > 0 && !(curX == MAPSIZE/2 && curY == MAPSIZE/2) {
		m[curX][curY] = "&Y@&W"
	}
	depth++
	if depth > 6 {
		return m
	}
	db := DB()
	for d, e := range room.Exits {
		switch d {
		case "east":
			x := curX + 1
			y := curY
			if map_in_bounds(x, y) {
				m[x][y] = sprintf("&d%s&W", MAP_EXIT_EW)
				m = walk_rooms(m, db.GetRoom(e, room.ship), x+1, y, depth)
			}
		case "west":
			x := curX - 1
			y := curY
			if map_in_bounds(x, y) {
				m[x][y] = sprintf("&d%s&W", MAP_EXIT_EW)
				m = walk_rooms(m, db.GetRoom(e, room.ship), x-1, curY, depth)
			}
		case "north":
			x := curX
			y := curY - 1
			if map_in_bounds(x, y) {
				m[x][y] = sprintf("&d%s&W", MAP_EXIT_NS)
				m = walk_rooms(m, db.GetRoom(e, room.ship), x, y-1, depth)
			}
		case "south":
			x := curX
			y := curY + 1
			if map_in_bounds(x, y) {
				m[x][y] = sprintf("&d%s&W", MAP_EXIT_NS)
				m = walk_rooms(m, db.GetRoom(e, room.ship), x, y+1, depth)
			}
		case "northwest":
			x := curX - 1
			y := curY - 1
			if map_in_bounds(x, y) {
				m[x][y] = sprintf("&d%s&W", MAP_EXIT_NWSE)
				m = walk_rooms(m, db.GetRoom(e, room.ship), x-1, y-1, depth)
			}
		case "southeast":
			x := curX + 1
			y := curY + 1
			if map_in_bounds(x, y) {
				m[x][y] = sprintf("&d%s&W", MAP_EXIT_NWSE)
				m = walk_rooms(m, db.GetRoom(e, room.ship), x+1, y+1, depth)
			}
		case "northeast":
			x := curX + 1
			y := curY - 1
			if map_in_bounds(x, y) {
				m[x][y] = sprintf("&d%s&W", MAP_EXIT_SWNE)
				m = walk_rooms(m, db.GetRoom(e, room.ship), x+1, y-1, depth)
			}
		case "southwest":
			x := curX - 1
			y := curY + 1
			if map_in_bounds(x, y) {
				m[x][y] = sprintf("&d%s&W", MAP_EXIT_SWNE)
				m = walk_rooms(m, db.GetRoom(e, room.ship), x-1, y+1, depth)
			}
		}
	}
	return m
}
