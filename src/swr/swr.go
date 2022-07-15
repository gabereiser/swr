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
	"flag"
	"fmt"
	"log"
	"os"
)

var editMode = flag.Bool("editmode", false, "Used to run the server in editor mode for offline world building.")

func Init() {
	// Ensure that the player directories exists
	for _, p := range "abcdefghijklmnopqrstuvwxyz" {
		_ = os.MkdirAll(fmt.Sprintf("data/accounts/%s", string(p)), 0755)
	}
	// Start the scheduler
	Scheduler()
}

func Main() {

	flag.Parse()

	log.Printf("Starting version %s\n", version)

	DB().Load()
	CommandsLoad()
	LanguageLoad()

	if *editMode {
		RunEditor()
	} else {
		StartBackup()
		ServerStart(Config().Addr)
	}

	DB().Save()
}

func GetVersion() string {
	return version
}

func RunEditor() {
	log.Printf("Starting editor\n")
	editor := NewEditor()
	editor.Run()
}
