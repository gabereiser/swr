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
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

const (
	NET_IAC  = byte(255)
	NET_WILL = byte(251)
	NET_WONT = byte(252)
	NET_DO   = byte(253)
	NET_DONT = byte(254)
	NET_ECHO = byte(1)
	NET_GA   = byte(3)
)

var ServerRunning bool = false
var ServerQueue chan MudClientCommand = make(chan MudClientCommand)

type MudClientCommand struct {
	Entity  Entity
	Command string
}

type MudClient struct {
	Id      string
	Con     *net.TCPConn
	Closed  bool
	Idle    int
	Editing bool
	EditPtr *string
	Queue   chan string
}

func (c *MudClient) Send(str string) {
	str = Color().Colorize(str)
	if c.Editing {
		c.Queue <- str
	} else {
		_, e := c.Con.Write([]byte(str))
		if e == io.EOF {
			c.Close()
		}
		if e == io.ErrClosedPipe {
			c.Close()
		}
		if e == net.ErrClosed {
			c.Close()
		}
	}

}

func (c *MudClient) Sendf(format string, any ...interface{}) {
	c.Send(fmt.Sprintf(format, any...))
}

func (c *MudClient) Read() string {
	b := make([]byte, 1)
	buf := ""
	for {
		if c.Closed {
			break
		}
		if c.Editing {
			break
		}
		c.Con.SetReadDeadline(time.Now().Add(1 * time.Hour).Add(1 * time.Second))
		i, err := c.Con.Read(b)
		if err == io.EOF {
			c.Close()
			c.Con.Close()
			return buf
		}
		if err != nil {
			c.Close()
			c.Con.Close()
			return buf
		}
		if i > 0 {
			buf += string(b)
			if strings.HasSuffix(buf, "\n") {
				buf = strings.TrimSuffix(buf, "\r\n")
				buf = strings.TrimSuffix(buf, "\n")
				return strings.TrimSpace(buf)
			}
		}
	}
	return buf
}

func (c *MudClient) Close() {
	c.Closed = true
	c.Con.CloseRead()
}

func (c *MudClient) BufferEditor(str *string) {
	if !c.Editing {
		c.EditPtr = str
		c.Editing = true
	}
}

func (c *MudClient) Raw(buffer []byte) {
	c.Con.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, e := c.Con.Write(buffer)
	if e == io.EOF {
		c.Close()
		return
	}
	if e != nil {
		ErrorCheck(e)
	}
}

type Client interface {
	Raw(buffer []byte)
	Send(str string)
	Sendf(format string, any ...interface{})
	Read() string
	BufferEditor(buf *string)
	Close()
}

func ServerStart(addr string) {
	a, _ := net.ResolveTCPAddr("tcp", addr)
	l, err := net.ListenTCP("tcp", a)
	ErrorCheck(err)
	defer l.Close()
	log.Printf("Listening for connections on %s\n", addr)
	ServerRunning = true
	go processClients()
	go processServerPump()
	for {
		if !ServerRunning {
			break
		}
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
	ServerRunning = false
}
func processClients() {
	for {
		if !ServerRunning {
			break
		}
		cmd := <-ServerQueue
		do_command(cmd.Entity, cmd.Command)
		time.Sleep(500 * time.Millisecond)
	}
}
func processServerPump() {
	for {
		if !ServerRunning {
			break
		}
		processIdleClients()
		processCombat()
		processEntities()
		time.Sleep(1 * time.Second)
	}
	log.Printf("Server Pump has exited!\n")
}
func acceptClient(con *net.TCPConn) {
	//telnet_suppress_ga(con)
	db := DB()
	client := new(MudClient)
	client.Id = hex.EncodeToString([]byte(con.RemoteAddr().String()))
	client.Con = con
	client.Closed = false
	client.Idle = 0

	auth_do_welcome(client)
	if client.Closed {
		return
	}
	db.AddClient(client)
	entity := db.GetEntityForClient(client)
	if entity != nil {
		for {
			if !ServerRunning {
				break
			}
			if client.Closed {
				break
			}
			if client.Editing {
				run_editor(entity, client.EditPtr)
			} else {
				input := client.Read()
				if len(input) > 0 {
					if strings.HasPrefix(input, "editor ") {
						parts := strings.Split(input, " ")[1:]
						buf := strings.Join(parts, " ")
						client.BufferEditor(&buf)
					} else {
						ServerQueue <- MudClientCommand{
							Entity:  entity,
							Command: input,
						}
					}
					client.Idle = 0
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
	db.RemoveEntity(entity)
	log.Printf("Player %s has left the game.", entity.GetCharData().Name)
	db.RemoveClient(client)
	room := DB().GetRoom(entity.RoomId(), entity.ShipId())
	room.SendToRoom(fmt.Sprintf("\r\n&P%s&d has left.\r\n", entity.GetCharData().Name))
	con.Close()
}

func processIdleClients() {
	db := DB()
	db.Lock()
	defer db.Unlock()
	for i := range db.clients {
		client := db.clients[i]
		if client != nil {
			client.Idle++
			minuteSeconds := 60 * 60
			if client.Idle == minuteSeconds-60 {
				client.Sendf("\r\n}YConnection Idle Warning!!&d &wYou have been idle for %d minutes. You're connection will close in 1 minute.&d\r\n", client.Idle/60)
			}
			if client.Idle == minuteSeconds-30 {
				client.Sendf("\r\n}YConnection Idle Warning!!&d &wYou have been idle for %d minutes. You're connection will close in 30 seconds.&d\r\n", client.Idle/60)
			}
			if client.Idle == minuteSeconds-15 {
				client.Sendf("\r\n}YConnection Idle Warning!!&d &wYou have been idle for %d minutes. You're connection will close in 15 seconds.&d\r\n", client.Idle/60)
			}
			if client.Idle > minuteSeconds {
				client.Send("\r\n&xClosing idle connection...&d\r\n")
				client.Close()
			}
		}

	}
}

func run_editor(entity Entity, buffer *string) {
	//client := entity.(*PlayerProfile).Client.(*MudClient)
	//telnet_disable_local_echo(client.Con)
	//telnet_suppress_ga(client.Con)
	value := editor(entity, *buffer)
	if value == "" {
		log.Printf("Editor has no contents, not setting buffer")
	}
	p := entity.(*PlayerProfile)
	c := p.Client.(*MudClient)
	c.Editing = false
	for i := 0; i < len(c.Queue); i++ {
		c.Send(<-c.Queue)
	}

	//telnet_unsuppress_ga(client.Con)
	//telnet_enable_local_echo(client.Con)
}

/*
func telnet_suppress_ga(con *net.TCPConn) {
	_, err := con.Write([]byte{NET_IAC, NET_WILL, NET_GA})
	if err != nil {
		con.Close()
	}
	resp := make([]byte, 3)
	con.Read(resp)
	if resp[0] != NET_IAC || resp[1] != NET_DO || resp[2] != NET_GA {
		panic(fmt.Sprintf("%v", resp))
	}
}
func telnet_unsuppress_ga(con *net.TCPConn) {
	_, err := con.Write([]byte{NET_IAC, NET_WONT, NET_GA})
	if err != nil {
		con.Close()
	}
	resp := make([]byte, 3)
	con.Read(resp)
	if resp[0] != NET_IAC || resp[1] != NET_DONT || resp[2] != NET_GA {
		panic(fmt.Sprintf("%v", resp))
	}
}
*/
/*
func telnet_disable_local_echo(con *net.TCPConn) {
	_, err := con.Write([]byte{NET_IAC, NET_WILL, NET_ECHO})
	if err != nil {
		con.Close()
	}
	resp := make([]byte, 3)
	con.Read(resp)
	if resp[0] != NET_IAC || resp[1] != NET_DO || resp[2] != NET_ECHO {
		panic(fmt.Sprintf("%v", resp))
	}
}

func telnet_enable_local_echo(con *net.TCPConn) {
	_, err := con.Write([]byte{NET_IAC, NET_WONT, NET_ECHO})
	if err != nil {
		con.Close()
	}
	resp := make([]byte, 3)
	con.Read(resp)
	if resp[0] != NET_IAC || resp[1] != NET_DONT || resp[2] != NET_ECHO {
		panic(fmt.Sprintf("%v", resp))
	}
}*/
