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
	"bufio"
	"encoding/hex"
	"log"
	"net"
	"time"
)

type MudClient struct {
	Id     string
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func (c *MudClient) Send(str string) {
	_, err := c.Writer.Write([]byte(str))
	ErrorCheck(err)
	err = c.Writer.Flush()
	ErrorCheck(err)
}

func (c *MudClient) Read() string {
	s, _, err := c.Reader.ReadLine()
	ErrorCheck(err)
	return string(s)
}

type Client interface {
	Send(str string)
	Read() string
}

func ServerStart(addr string) {
	a, _ := net.ResolveTCPAddr("tcp", addr)
	l, err := net.ListenTCP("tcp", a)
	ErrorCheck(err)
	defer l.Close()
	log.Printf("Listening for connections on %s\n", addr)
	for {
		c, err := l.AcceptTCP()
		if err != nil {
			log.Printf("Error accepting a connection: %v", err)
			continue
		}
		if c != nil {
			go acceptClient(c)
			log.Printf("Accepted client connection from %s", c.RemoteAddr())
		}

	}
}

func acceptClient(con *net.TCPConn) {

	client := new(MudClient)
	client.Id = hex.EncodeToString([]byte(con.RemoteAddr().String()))
	client.Reader = bufio.NewReader(con)
	client.Writer = bufio.NewWriter(con)

	auth_do_welcome(client)

	db := DB()
	db.AddClient(client)
	for {
		_, err := client.Reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading message from %s %v", client.Id, err)
			break
		}

		time.Sleep(time.Duration(1) * time.Millisecond)
	}
	db.RemoveClient(client)
	con.Close()
}
