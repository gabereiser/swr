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
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

func main() {
	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		url = "127.0.0.1:5000"
	}
	addr, err := net.ResolveTCPAddr("tcp", url)
	if err != nil {
		panic(fmt.Errorf("error parsing connection URL"))
	}
	con, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(fmt.Errorf("error, %v", err))
	}
	defer con.Close()
	fmt.Println("Connected...")

	var waitGroup sync.WaitGroup

	waitGroup.Add(2)
	fmt.Print("\x1b[")
	go func() {
		defer waitGroup.Done()
		read(con)
	}()
	go func() {
		defer waitGroup.Done()
		write(con)
	}()
	waitGroup.Wait()

}

func read(con *net.TCPConn) {
	con.SetReadBuffer(4096)
	f, err := con.File()
	if err != nil {
		panic(fmt.Errorf("error getting file descriptor: %v", err))
	}
	for {
		d := make([]byte, 1)
		_, err := f.Read(d)
		if err == io.EOF {
			os.Exit(1)
		}
		os.Stdout.Write(d)
	}
}

func write(con *net.TCPConn) {
	con.SetWriteBuffer(1)
	f, err := con.File()
	if err != nil {
		panic(fmt.Errorf("error getting file descriptor: %v", err))
	}
	for {
		i := make([]byte, 1)
		_, err := os.Stdin.Read(i)
		if err == io.EOF {
			os.Exit(1)
		}
		f.Write(i)
	}
}
