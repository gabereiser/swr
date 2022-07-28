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
	"log"
	"os"
)

func Init() {
	// Ensure that the player directories exists
	for _, p := range "abcdefghijklmnopqrstuvwxyz" {
		_ = os.MkdirAll(fmt.Sprintf("data/accounts/%s", string(p)), 0755)
	}
	_ = os.MkdirAll("backup", 0755)
	// Start the scheduler
	Scheduler()
}

func Main() {

	log.Printf("Starting version %s\n", version)
	is_skill("martial-arts")
	DB().Load()
	DB().ResetAll()
	CommandsLoad()
	LanguageLoad()
	StartBackup()
	EditorStart()
	ServerStart(Config().Addr)

	DB().Save()
}

func GetVersion() string {
	return version
}
