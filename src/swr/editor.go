package swr

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func do_editor(entity Entity, args ...string) {
	result := run_editor(entity, strings.Join(args, " "))
	entity.Send("\r\nResult: %s\r\n", result)
}

//lint:ignore U1000 useful code
func run_editor(entity Entity, buffer string) string {
	client := entity.(*PlayerProfile).Client
	//telnet_disable_local_echo(client.Con)
	telnet_suppress_ga(client)
	value := editor(entity, buffer)
	if value == "" {
		log.Printf("Editor has no contents, not setting buffer")
	}
	p := entity.(*PlayerProfile)
	c := p.Client
	c.SetEditing(false)
	c.SendQueue()
	c.ClearQueue()
	telnet_unsuppress_ga(client)
	return value
	//telnet_enable_local_echo(client.Con)
}
func editor(entity Entity, contents string) string {
	player := entity.(*PlayerProfile)
	client := player.Client.(*TCPClient)
	filename := sprintf("/tmp/%s", strings.ToLower(strings.ReplaceAll(entity.GetCharData().Name, " ", "")))
	e := ioutil.WriteFile(filename, []byte(contents), 0755)
	ErrorCheck(e)
	cmd := exec.Command("vim", "-Z", filename)
	cmd.Env = append(cmd.Env, "TERM=xterm256-color")
	cmd.Stdin = client.fd
	cmd.Stdout = client.Con
	cmd.Stderr = client.Con
	err := cmd.Run()
	ErrorCheck(err)
	buf, _ := ioutil.ReadFile(filename)
	return string(buf)

}

/*
func editor(entity Entity, args ...string) string {
	player := entity.(*PlayerProfile)
	client := player.Client
	client.SetEditing(true)
	con := client

	con.Raw([]byte("\x1b[2J\x1b[H"))
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
		n, e := con.ReadRaw(b[:])
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
				break
			} else if k == 27 {
				if ctrl {
					con.Raw([]byte("\x1b[u"))
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

func print_editor(con Client, buffer []byte) {
	con.Raw([]byte("\x1b[2J\x1b[H"))
	str := string(buffer)
	lines := strings.Split(str, "\n")
	for i := 0; i < 9; i++ {
		con.Raw([]byte(sprintf("\x1b[1;32m%1d\x1b[m ", i)))
		if i < len(lines) {
			con.Raw([]byte(lines[i]))
		}
		con.Raw([]byte("\r\n"))
	}
	con.Raw([]byte("-------------------------------------------------------------------------\r\n"))
	con.Raw([]byte("Quit: \\q  |  Write: \\w \r\n"))
}
*/
