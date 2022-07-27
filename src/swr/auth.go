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
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func auth_do_welcome(client Client) {
	client.Send(Color().ClearScreen())
	client.Send("\r\n\r\n&CA long time ago in a galaxy far, far away...&d\r\n\r\n\r\n[press ENTER]")
	_ = client.Read()
	welcome, err := ioutil.ReadFile("data/sys/welcome")
	ErrorCheck(err)
	client.Send(string(welcome))
	auth_do_login(client)
}

func auth_do_login(client Client) {
Login:
	client.Send("\r\n&GHolonet Login:&d ")
	username := client.Read()
	if username == "" {
		goto Login
	}
	sanitized := strings.ToLower(username)
	if strings.Contains(sanitized, " ") {
		client.Send("\r\n}RInvalid login!&d &RSpace aren't allowed.&d\r\n")
		goto Login
	}
	path := fmt.Sprintf("data/accounts/%s/%s.yml", sanitized[0:1], sanitized)
	log.Printf("Loading player %s", sanitized)
	if FileExists(path) {
		player := DB().ReadPlayerData(path)
		client.Send("\r\n&GPassword:&d ")
		telnet_disable_local_echo(client.(*MudClient).Con)
		password := client.Read()
		telnet_enable_local_echo(client.(*MudClient).Con)
		if encrypt_string(password) != player.Password {
			fmt.Printf("%s\n", password)
			client.Send("\r\n}RInvalid password!&d\r\n")
			goto Login
		}
		if encrypt_string(password) == player.Password {
			client.Send(fmt.Sprintf("\r\n&GAccess granted! Welcome %s.&d\r\n", player.Char.Name))
			time.Sleep(1 * time.Second)
			client.Send(Color().ClearScreen())
			player.LastSeen = time.Now()
			player.Client = client
			DB().SavePlayerData(player)
			room := DB().GetRoom(player.Char.Room)
			room.SendToRoom(fmt.Sprintf("\r\n&P%s&d has arrived.\r\n", player.Char.Name))
			// see if player is already in the game...
			p := DB().GetPlayerEntityByName(player.Char.Name)
			if p == nil {
				DB().AddEntity(player)
			} else {
				client.Send("\r\nReconnecting to player...\r\n")
				player = p.(*PlayerProfile)
				player.Client = client
				player.LastSeen = time.Now()
			}
			ServerQueue <- MudClientCommand{
				Entity:  player,
				Command: "look",
			}
		}
	} else {
		client.Send(fmt.Sprintf("\r\n&RHrm, it seems there isn't a record of &Y%s&R in the galactic databank.\r\n\r\n&RAre you new? &G[&Wy&G/&Wn&G]&d ", username))
		are_new := strings.ToLower(client.Read())
		if strings.HasPrefix(are_new, "y") {
			player := new(PlayerProfile)
			player.Char = CharData{}
			player.Char.Name = username
			auth_do_new_player(client, player)
		} else {
			goto Login
		}
	}
}

func auth_do_new_player(client Client, player *PlayerProfile) {
	// ch is a new Character. Allocated but unassigned in the game world.
	// complete initialization, associate, and load into the game as that
	// character.
Name:
	client.Send("\r\n&GCharacter Name:&d ")
	name := client.Read()
	if strings.ContainsAny(name, " `~,./?<>;:'\"[]}{\\|+_-=!@#$%^&*()") {
		client.Send("}RSpaces and special characters are not allowed.&d\r\n\r\n")
		goto Name
	}
	client.Sendf("\r\n&GYou will be known as &W%s&G. Is that ok? [&Wy&G/&Wn&G] &d", name)
	name_confirm := client.Read()
	if !strings.HasPrefix(strings.ToLower(name_confirm), "y") {
		goto Name
	}
Password:
	client.Sendf("\r\n&GWelcome &W%s&G.\r\n&GPlease enter a &Wpassword&G:&d ", name)
	telnet_disable_local_echo(client.(*MudClient).Con)
	password := client.Read()
	client.Send("\r\n&GRepeat your &Wpassword&G:&d ")
	password2 := client.Read()
	telnet_enable_local_echo(client.(*MudClient).Con)
	if password != password2 {
		client.Send("\r\n}RError! Password mismatch!&d\r\n")
		goto Password
	}
Email:
	client.Send("\r\n&GPlease enter your email &X(we won't spam you)&G:&d ")
	email := client.Read()
	if !strings.Contains(email, "@") {
		client.Send("\r\n}RError, an email address is needed for account recovery purposes.&d\r\n")
		goto Email
	}

Race:
	client.Send("\r\n&G-=-=-=-=-=-=-=-=-=-=-=( &WChoose Your Race &G)=-=-=-=-=-=-=-=-=-=-=-=-=-=-&d\r\n\r\n")

	buf := ""
	for i, race := range race_list {
		buf += fmt.Sprintf("&Y[&w%2d&Y] &W%-12s\t", i+1, race)
		if (i+1)%3 == 0 {
			buf += "\r\n"
		}
		if i >= 31 {
			break
		}
	}
	client.Send(buf)
	client.Sendf("\r\n&GRace Selection [1-%d]:&d ", 32)
	race := client.Read()
	r_index, err := strconv.Atoi(race)
	if err != nil {
		client.Send("}RUnable to parse race, please use a number!&d")
		ErrorCheck(err)
		goto Race
	}
	if r_index > 32 {
		client.Send("}RNumber outside of bounds. Try again.&d")
		goto Race
	}
	race = race_list[r_index-1]
Gender:
	client.Send("\r\n\r\n&G-=-=-=-=-=-=-=-=-=-=-=( &WChoose Your Gender &G)=-=-=-=-=-=-=-=-=-=-=-=-=-&d")
	client.Send("\r\n&GYour character needs a gender. You can be &Wmale&G, &Wfemale&G, or &Wnon-binary&G/&Wneutral&G.\r\n")
	client.Send("&W[&GM&W/&GF&W/&GN&W]:&d ")
	gender := strings.ToLower(client.Read())
	if gender[0:1] != "m" && gender[0:1] != "f" && gender[0:1] != "n" {
		client.Send("\r\n}RGender not recognized, please try again.&d\r\n")
		goto Gender
	}
	switch gender[0:1] {
	case "m":
	case "M":
		gender = "Male"
	case "f":
	case "F":
		gender = "Female"
	default:
		gender = "Non"
	}

	client.Send("\r\n\r\n&G-=-=-=-=-=-=-=-=-=-=-=( &WStats &G)=-=-=-=-=-=-=-=-=-=-=-=-=-&d")
	stats := make([]int, 6)
Stats:
	stats[0] = rand_min_max(1, 6) + rand_min_max(1, 6) + rand_min_max(1, 6)
	stats[1] = rand_min_max(3, 6) + rand_min_max(3, 6) + rand_min_max(3, 6)
	stats[2] = rand_min_max(3, 6) + rand_min_max(3, 6) + rand_min_max(3, 6)
	stats[3] = rand_min_max(3, 6) + rand_min_max(3, 6) + rand_min_max(3, 6)
	stats[4] = rand_min_max(3, 6) + rand_min_max(3, 6) + rand_min_max(3, 6)
	stats[5] = rand_min_max(3, 6) + rand_min_max(3, 6) + rand_min_max(3, 6)
	client.Send(fmt.Sprintf("\r\n\r\n&YSTR: &w%d  &YINT: &w%d  &YDEX: &w%d  &YWIS: &w%d  &YCON: &w%d  &YCHA: &w%d\r\n\r\n", stats[0], stats[1], stats[2], stats[3], stats[4], stats[5]))
	client.Send("&GAre these ok? &G[&Wy&G/&Wn&G]&d ")
	if !strings.HasPrefix(strings.ToLower(client.Read()), "y") {
		goto Stats
	}
	name = strings.ToLower(name)
	name_t := strings.ToUpper(name[0:1])
	name = name_t + name[1:]
	player.Char = CharData{}
	player.Char.Id = gen_player_char_id()
	player.Char.Name = name
	player.Char.Room = 100
	player.Char.Race = race
	player.Char.Gender = gender
	player.Char.Title = fmt.Sprintf("%s the %s", player.Char.Name, player.Char.Race)
	player.Char.Level = 1
	player.Char.XP = 0
	player.Char.Gold = 0
	player.Char.Stats = stats
	player.Char.Skills = map[string]int{"kick": 1, "beg": 1, "search": 1}
	player.Char.Hp = []int{10, 10}
	player.Char.Mp = []int{0, 0}
	player.Char.Mv = []int{10, 10}
	player.Char.Equipment = make(map[string]*ItemData)
	player.Char.Inventory = make([]*ItemData, 0)
	player.Char.Keywords = []string{name, race}
	player.Char.Bank = 0
	player.Char.Brain = "client"
	if race == "Human" {
		player.Char.Speaking = "basic"
	} else {
		player.Char.Speaking = strings.ToLower(race)
	}
	player.Char.Languages = make(map[string]int)
	player.Char.Languages["basic"] = 100
	player.Char.Languages[player.Char.Speaking] = 100

	player.LastSeen = time.Now()
	player.Password = encrypt_string(password)
	player.Email = email
	player.Banned = false
	player.Frequency = tune_random_frequency()
	player.Priv = 1

	client.Sendf("\r\n\r\n&GYou are about to create the character &W%s the %s&G.\r\nAre you ok with this? [&Wy&G/&Wn&G]&d ", name, race)
	k := client.Read()
	if strings.ToLower(k[0:1]) != "y" {
		client.Send("\r\nGoodbye.\r\n\r\n&xThe terminal view fades away and all you see is black.&d\r\n")
		client.Close()
		return
	}
	DB().SavePlayerData(player)
	player.Client = client
	room := DB().GetRoom(player.Char.Room)
	room.SendToRoom(fmt.Sprintf("\r\n&P%s&d has arrived.\r\n", player.Char.Name))
	DB().AddEntity(player)
	client.Send(Color().ClearScreen())
	client.Send("\r\nEntering game world...\r\n")
	ServerQueue <- MudClientCommand{
		Entity:  player,
		Command: "look",
	}

}

func encrypt_string(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
