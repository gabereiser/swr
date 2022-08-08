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
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func assert(expr bool) {
	if !expr {
		buf := make([]byte, 0)
		runtime.Stack(buf, true)
		panic(sprintf("%s\n", string(buf)))
	}
}
func bytes_to_mb(b uint64) float64 {
	return float64(b) / (1024 * 1024)
}

func file_exists(filename string) bool {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("%s does not exist!", filename)
		return false
	}
	defer file.Close()
	return true
}
func FileExists(filename string) bool {
	return file_exists(filename)
}

// Replaces line endings to enforce CR+LF instead of just LF
func telnet_encode(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, "\r\n", "\n"), "\n", "\r\n")
}

func sprintf(format string, any ...interface{}) string {
	return fmt.Sprintf(format, any...)
}

// D20 System Dice Roll Mechanic
// Supported formats are...
// "1d4" for 1 4-sided dice
// "1d20" for 1 20-sided dice
// "2d10+10" for 2 10-sided dice + 10 after the roll.
func roll_dice(d20 string) int {
	mods := strings.Split(d20, "+")
	p := strings.Split(strings.ToLower(mods[0]), "d")
	num_dice, _ := strconv.Atoi(p[0])
	if num_dice < 1 {
		return 0
	}
	sides, _ := strconv.Atoi(p[1])
	if sides < 1 {
		return 0
	}
	roll := 0
	for i := 0; i < num_dice; i++ {
		roll += rand.Intn(sides) + 1
	}
	if len(mods) == 2 {
		mod, _ := strconv.Atoi(mods[1])
		roll += mod
	}
	return roll
}

func rand_min_max(min int, max int) int {
	return min + rand.Intn((max-min)+1)
}

//lint:ignore U1000 int version of umin
func min(min int, value int) int {
	if value < min {
		return value
	}
	return min
}

//lint:ignore U1000 int version of umin
func umin(min uint, value uint) uint {
	if value < min {
		return value
	}
	return min
}
func umax(min uint, value uint) uint {
	if value > min {
		return value
	}
	return min
}
func random_seed(seed int64) {
	rand.Seed(seed)
}
func random_float() float64 {
	return rand.Float64()
}
func gen_player_char_id() uint {
	return uint(rand.Intn(1000000000)) + 900000000 // 900000000-999999999
}
func gen_ship_id() uint {
	return uint(rand.Intn(1000000000)) + 300000000 // 300000000-399999999
}

func gen_npc_char_id() uint {
	return uint(rand.Intn(1000000000)) + 100000000 // 100000000-199999999
}

func gen_item_id() uint {
	return uint(rand.Intn(1000000000)) + 200000000 // 200000000-299999999
}

func tune_random_frequency() string {
	buf := ""
	buf += strconv.Itoa(rand.Intn(3) + 1) // 1,2,3,4
	buf += strconv.Itoa(rand.Intn(9))     // 0-9
	buf += strconv.Itoa(rand.Intn(9))     // 0-9
	buf += "."
	switch rand.Intn(3) {
	case 0:
		buf += "000"
	case 1:
		buf += "250"
	case 2:
		buf += "500"
	case 3:
		buf += "750"
	}
	return buf
}

// capitalize makes the first letter of each word... capitalized.
func capitalize(str string) string {
	parts := strings.Split(str, " ")
	ret := ""
	for _, p := range parts {
		ret += sprintf("%s%s ", strings.ToUpper(p[0:1]), strings.ToLower(p[1:]))
	}
	return strings.TrimSpace(ret)
}

// consolify takes a long string and chops it up by word to limit it to 80 character width.
// useful for terminals and telnet.
func consolify(str string) string {
	if len(str) < 70 {
		return str
	}
	words := strings.Split(str, " ")
	cursor := 1
	buf := ""
	for _, w := range words {
		wlen := len(w)
		if cursor+wlen > 70 {
			buf += "\r\n"
			cursor = 1
		}
		buf += sprintf("%s ", w)
		cursor += wlen + 1 // +1 for the space
	}
	return buf
}

//lint:ignore U1000 useful code
func slice_contains_string(slice []string, value string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, value) {
			return true
		}
	}
	return false
}

// direction_reverse takes a direction string and returns its spacial opposite direction.
// ex: east -> west   north -> south   up -> down
func direction_reverse(direction string) string {
	switch direction {
	case "north":
		return "south"
	case "south":
		return "north"
	case "east":
		return "west"
	case "west":
		return "east"
	case "northwest":
		return "southeast"
	case "northeast":
		return "southwest"
	case "southwest":
		return "northeast"
	case "southeast":
		return "northwest"
	case "up":
		return "down"
	case "down":
		return "up"
	default:
		return "somewhere"
	}
}

// get_gender_for_code takes a gender code (m/f/n) and returns a lowercase printable name.
// Use [capitalize] if you want to make it pretty.
func get_gender_for_code(gender string) string {
	g := strings.ToLower(gender)
	if g[0:1] == "m" {
		return "male"
	}
	if g[0:1] == "f" {
		return "female"
	}
	if g[0:1] == "n" {
		return "neuter"
	}
	return "male"
}

func sqrt32(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}
func pow32(v float32, i int) float32 {
	return float32(math.Pow(float64(v), float64(i)))
}

func distance_between_points(origin []float32, dest []float32) float32 {
	return sqrt32(pow32(dest[0]-origin[0], 2) + pow32(dest[1]-origin[1], 2)*1.0)
}

var ZERO_DISTANCE float32 = distance_between_points([]float32{0.0, 0.0}, []float32{0.0, 0.0})
var MAX_DISTANCE float32 = distance_between_points([]float32{-100000.0, -100000.0}, []float32{100000.0, 100000.0})
