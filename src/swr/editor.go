package swr

import (
	"log"
	"net"
	"strings"
)

func do_editor(entity Entity, args ...string) {
	client := entity.(*PlayerProfile).Client.(*MudClient)
	buffer := strings.Join(args, " ")
	client.BufferEditor(&buffer)
}

func editor(entity Entity, args ...string) string {
	player := entity.(*PlayerProfile)
	client := player.Client.(*MudClient)
	client.Editing = true
	con := client.Con

	con.Write([]byte("\x1b[2J\x1b[H"))
	contents := strings.Join(args, " ")
	curX := 0
	curY := 0
	buf := []byte(contents)
	ctrl := false
	esc := false
	for {
		print_editor(con, buf)
		client.Raw([]byte(sprintf("\x1b[H\x1b[%d;%dH", curY+1, curX+5)))
		var b [1]byte
		n, e := con.Read(b[:])
		ErrorCheck(e)
		if n > 0 {
			k := b[0]
			if k == 92 {
				client.Raw([]byte("\x1b[s\x1b[999;1H"))
				client.Raw([]byte("\\"))
				ctrl = true
			} else if k == 113 && ctrl {
				client.Raw([]byte("q"))
				log.Printf("Exiting edit mode...")
				client.Raw([]byte("\x1b[999;1H"))
				player.NeedPrompt = true
				break
			} else if k == 27 {
				if ctrl {
					con.Write([]byte("\x1b[u"))
					ctrl = false
				}
				esc = true
			} else if k == 37 && esc {
				curX -= 1
				if curX < 0 {
					curX = 0
				}
			} else if k == 38 && esc {
				curY -= 1
				if curY < 0 {
					curY = 0
				}
			} else if k == 39 && esc {
				curX += 1
				if curX > 80 {
					curX = 80
				}
			} else if k == 40 && esc {
				curY += 1
				if curY > 9 {
					curY = 9
				}
			} else {
				esc = false
				ctrl = false
				curX += 1
				buf = append(buf, k)
			}
		}
	}
	return ""
}

func print_editor(con net.Conn, buffer []byte) {
	con.Write([]byte("\x1b[2J\x1b[H"))
	str := string(buffer)
	lines := strings.Split(str, "\n")
	for i := 0; i < 9; i++ {
		con.Write([]byte(sprintf("\x1b[1;32m%1d\x1b[m ", i)))
		if i < len(lines) {
			con.Write([]byte(lines[i]))
		}
		con.Write([]byte("\r\n"))
	}
	con.Write([]byte("-------------------------------------------------------------------------\r\n"))
	con.Write([]byte("Quit: \\q  |  Write: \\w \r\n"))
}
