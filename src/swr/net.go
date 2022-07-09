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
	db := DB()
	client := new(MudClient)
	client.Id = hex.EncodeToString([]byte(con.RemoteAddr().String()))
	client.Reader = bufio.NewReader(con)
	client.Writer = bufio.NewWriter(con)

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
