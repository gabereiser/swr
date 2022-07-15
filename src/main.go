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
package main

import (
	"fmt"
	"time"

	swr "github.com/gabereiser/swr"
)

func init() {
	fmt.Printf("%d %d %d %d\r\n", int('a')-97, int('A')-65, int('z')-97, int('Z')-65)
	if version != swr.GetVersion() {
		panic(fmt.Sprintf("Version Mismatch! %s != %s", version, swr.GetVersion()))
	}
	fmt.Println(`SWR  Copyright (C) 2022
This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it
under certain conditions; see LICENSE for details.

-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-`)

	time.Sleep(1 * time.Second)
	swr.Init()
}

func main() {
	swr.Main()
}
