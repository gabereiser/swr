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
	"os"
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

type TCPClient struct {
	Id      string
	Con     *net.TCPConn
	fd      *os.File
	Closed  bool
	Idle    int
	Editing bool
	EditPtr *string
	Queue   []string
}

func (c *TCPClient) Send(str string) {
	str = Color().Colorize(str)
	if c.Editing {
		c.Queue = append(c.Queue, str)
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

func (c *TCPClient) Sendf(format string, any ...interface{}) {
	c.Send(fmt.Sprintf(format, any...))
}
func (c *TCPClient) ReadRaw(b []byte) (int, error) {
	return c.Con.Read(b)
}
func (c *TCPClient) Read() string {
	b := make([]byte, 1)
	buf := ""
	for {
		if c.Closed {
			break
		}
		if c.Editing {
			break
		}
		i, err := c.Con.Read(b)
		if err != nil {
			c.Close()
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

func (c *TCPClient) Close() {
	c.Closed = true
	c.fd.Close()
	c.Con.Close()
}

func (c *TCPClient) IsClosed() bool {
	return c.Closed
}

func (c *TCPClient) BufferEditor(str *string) {
	if !c.Editing {
		c.EditPtr = str
		c.Editing = true
	}
}

func (c *TCPClient) Raw(buffer []byte) {
	_, e := c.Con.Write(buffer)
	if e == io.EOF {
		c.Close()
		return
	}
	if e != nil {
		ErrorCheck(e)
	}
}

func (c *TCPClient) GetId() string {
	return c.Id
}

func (c *TCPClient) SetEditing(editing bool) {
	c.Editing = editing
}

func (c *TCPClient) IsEditing() bool {
	return c.Editing
}

func (c *TCPClient) IdleInc() {
	c.Idle++
}
func (c *TCPClient) GetIdle() int {
	return c.Idle
}
func (c *TCPClient) SendQueue() {
	for _, s := range c.Queue {
		c.Con.Write([]byte(s))
	}
}
func (c *TCPClient) ClearQueue() {
	c.Queue = make([]string, 0)
}

type Client interface {
	IsClosed() bool
	Raw(buffer []byte)
	Send(str string)
	Sendf(format string, any ...interface{})
	Read() string
	ReadRaw(b []byte) (int, error)
	BufferEditor(buf *string)
	Close()
	GetId() string
	SetEditing(editing bool)
	IsEditing() bool
	IdleInc()
	GetIdle() int
	SendQueue()
	ClearQueue()
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
		updateMinerDifficulty()
		time.Sleep(1 * time.Second)
	}
	log.Printf("Server Pump has exited!\n")
}
func acceptClient(con *net.TCPConn) {
	fd, _ := con.File()
	db := DB()
	client := new(TCPClient)
	client.Id = hex.EncodeToString([]byte(con.RemoteAddr().String()))
	client.Con = con
	client.fd = fd
	client.Closed = false
	client.Idle = 0
	db.AddClient(client)
	auth_do_welcome(client)
	if client.Closed {
		return
	}
	entity := db.GetEntityForClient(client)
	if entity == nil {
		con.Close()
		db.RemoveClient(client)
		return
	}
	for {
		if !ServerRunning {
			break
		}
		if client.Closed {
			break
		} else {
			input := client.Read()
			if len(input) > 0 {
				ServerQueue <- MudClientCommand{
					Entity:  entity,
					Command: input,
				}
				client.Idle = 0
			}
		}
	}
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
			client.IdleInc()
			minuteSeconds := 60 * 60
			if client.GetIdle() == minuteSeconds-60 {
				client.Sendf("\r\n}YConnection Idle Warning!!&d &wYou have been idle for %d minutes. You're connection will close in 1 minute.&d\r\n", client.GetIdle()/60)
			}
			if client.GetIdle() == minuteSeconds-30 {
				client.Sendf("\r\n}YConnection Idle Warning!!&d &wYou have been idle for %d minutes. You're connection will close in 30 seconds.&d\r\n", client.GetIdle()/60)
			}
			if client.GetIdle() == minuteSeconds-15 {
				client.Sendf("\r\n}YConnection Idle Warning!!&d &wYou have been idle for %d minutes. You're connection will close in 15 seconds.&d\r\n", client.GetIdle()/60)
			}
			if client.GetIdle() > minuteSeconds {
				client.Send("\r\n&xClosing idle connection...&d\r\n")
				client.Close()
			}
		}

	}
}

//lint:ignore U1000 useful code
func telnet_suppress_ga(con Client) {
	con.Raw([]byte{NET_IAC, NET_WILL, NET_GA})
	/*resp := make([]byte, 3)
	n, e := con.ReadRaw(resp)
	if e != nil {
		con.Close()
	}
	if n != 3 {
		panic(fmt.Sprintf("%v", resp))
	}
	if resp[0] != NET_IAC || resp[1] != NET_DONT || resp[2] != NET_GA {
		panic(fmt.Sprintf("%v", resp))
	} else {
		log.Printf("SUPGA: %v", resp)
	}*/
}

//lint:ignore U1000 useful code
func telnet_unsuppress_ga(con Client) {
	con.Raw([]byte{NET_IAC, NET_WONT, NET_GA})
	/*resp := make([]byte, 3)
	con.ReadRaw(resp)
	if resp[0] != NET_IAC || resp[1] != NET_DO || resp[2] != NET_GA {
		panic(fmt.Sprintf("%v", resp))
	} else {
		log.Printf("SUPGA: %v", resp)
	}*/
}

//lint:ignore U1000 useful code
func telnet_disable_local_echo(con Client) {
	con.Raw([]byte{NET_IAC, NET_WILL, NET_ECHO})
	resp := make([]byte, 3)
	con.ReadRaw(resp)
	if resp[0] != NET_IAC || resp[1] != NET_DO || resp[2] != NET_ECHO {
		log.Print("Error: client didn't respond to disable local echo properly!")
	}
}

//lint:ignore U1000 useful code
func telnet_enable_local_echo(con Client) {
	con.Raw([]byte{NET_IAC, NET_WONT, NET_ECHO})
	resp := make([]byte, 3)
	con.ReadRaw(resp)
	if resp[0] != NET_IAC || resp[1] != NET_DONT || resp[2] != NET_ECHO {
		log.Print("Error: client didn't respond to enable local echo properly!")
	}
}
