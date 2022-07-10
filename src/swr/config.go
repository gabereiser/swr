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
	"encoding/json"
	"io/ioutil"
	"log"
)

type Configuration struct {
	Name string
	Addr string
}

var _config *Configuration

func Config() *Configuration {
	if _config == nil {
		fp, err := ioutil.ReadFile("./data/sys/config.json")
		ErrorCheck(err)
		err = json.Unmarshal(fp, &_config)
		ErrorCheck(err)
		log.Printf("Configuration loaded.")
	}
	return _config
}
