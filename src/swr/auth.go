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
	"fmt"
	"io/ioutil"
	"strings"
)

func auth_do_welcome(client Client) {
	welcome, err := ioutil.ReadFile("data/sys/welcome")
	ErrorCheck(err)
	client.Send(string(welcome))
	auth_do_login(client)
}

func auth_do_login(client Client) {
	client.Send("\r\n>>Holonet Login:")
	username := client.Read()
	sanitized := strings.ToLower(username)
	if strings.Contains(sanitized, " ") {
		client.Send("\r\nSpaces aren't allowed.\r\n")
		auth_do_login(client)
	}
	path := fmt.Sprintf("data/accounts/%s/%s", sanitized[0:1], sanitized)
	if FileExists(path) {
		char_data := DB().ReadCharData(path)
		auth_do_password(client, char_data)
	} else {
		client.Send(fmt.Sprintf("\r\nIt seems we have no record of %s.\r\n\r\nAre you new here?", username))
		are_new := strings.ToLower(client.Read())
		if are_new[0:1] == "y" {
			auth_do_new_player(client, new(CharData))
		} else {
			auth_do_login(client)
		}
	}
}

func auth_do_password(client Client, ch *CharData) {
	// ch is the loaded Character, it is not yet associated with a client.
	// verify passwords to associate and load into the game as that
	// character.

}

func auth_do_new_player(client Client, ch *CharData) {
	// ch is a new Character. Allocated but unassigned in the game world.
	// complete initialization, associate, and load into the game as that
	// character.
}
