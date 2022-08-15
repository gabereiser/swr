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
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	ANSI_ESC          = "\x1b["
	ANSI_CLEAR        = ANSI_ESC + "2J"
	ANSI_RESET        = ANSI_ESC + "0m"
	ANSI_CURSOR_UP    = ANSI_ESC + "1A"
	ANSI_CURSOR_DOWN  = ANSI_ESC + "1B"
	ANSI_CURSOR_LEFT  = ANSI_ESC + "1C"
	ANSI_CURSOR_RIGHT = ANSI_ESC + "1D"
	ANSI_BOLD         = ANSI_ESC + "1;"
	ANSI_ITALIC       = ANSI_ESC + "3m"
	ANSI_UNDERLINE    = ANSI_ESC + "4m"
	ANSI_BLINK        = ANSI_ESC + "5m"
)

const (
	ANSI_FG_BLACK   = "30"
	ANSI_FG_RED     = "31"
	ANSI_FG_GREEN   = "32"
	ANSI_FG_YELLOW  = "33"
	ANSI_FG_BLUE    = "34"
	ANSI_FG_MAGENTA = "35"
	ANSI_FG_CYAN    = "36"
	ANSI_FG_WHITE   = "37"
	ANSI_BG_BLACK   = "40"
	ANSI_BG_RED     = "41"
	ANSI_BG_GREEN   = "42"
	ANSI_BG_YELLOW  = "43"
	ANSI_BG_BLUE    = "44"
	ANSI_BG_MAGENTA = "45"
	ANSI_BG_CYAN    = "46"
	ANSI_BG_WHITE   = "47"
)

const (
	EMOJI_HEART_FULL   = "♥"
	EMOJI_HEART_EMPTY  = "♡"
	EMOJI_SWORDS       = "⚔️"
	EMOJI_SKULL        = "💀"
	EMOJI_HIT          = "💥"
	EMOJI_MAGIC        = "✨"
	EMOJI_TOMBSTONE    = "🪦"
	EMOJI_ALERT        = "🚨"
	EMOJI_COMM         = "⚡️"
	EMOJI_MONEY        = "💳"
	EMOJI_BANK         = "🏦"
	EMOJI_UNIVERSITY   = "🏫"
	EMOJI_HOSPITAL     = "🏥"
	EMOJI_MARKET       = "🏛"
	EMOJI_FACTORY      = "🏭"
	EMOJI_DRUGS        = "💊"
	EMOJI_RESEARCH     = "🧬"
	EMOJI_SCIENCE      = "🔬"
	EMOJI_SPEECH       = "💬"
	EMOJI_SHOUT        = "🗯"
	EMOJI_ANNOUNCEMENT = "📢"
	EMOJI_PACKAGE      = "📦"
	EMOJI_PARTY        = "🎉"
	EMOJI_HAPPY        = "🙂"
	EMOJI_ANGRY        = "😡"
	EMOJI_NEUTRAL      = "😐"
)

/*
┌─┬┐  ╔═╦╗  ╓─╥╖  ╒═╤╕
│ ││  ║ ║║  ║ ║║  │ ││
├─┼┤  ╠═╬╣  ╟─╫╢  ╞═╪╡
└─┴┘  ╚═╩╝  ╙─╨╜  ╘═╧╛
┌───────────────────┐
│  ╔═══╗ Some Text  │▒
│  ╚═╦═╝ in the box │▒
╞═╤══╩══╤═══════════╡▒
│ ├──┬──┤           │▒
│ └──┴──┘           │▒
└───────────────────┘▒

	▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒

	░ ▒ ▓ █
*/
const (
	ANSI_DBOX_TOP_LEFT     = "╔"
	ANSI_DBOX_HORIZONTAL   = "═"
	ANSI_DBOX_TOP_RIGHT    = "╗"
	ANSI_DBOX_VERTICAL     = "║"
	ANSI_DBOX_BOTTOM_LEFT  = "╚"
	ANSI_DBOX_BOTTOM_RIGHT = "╝"
	ANSI_BOX_TOP_LEFT      = "┌"
	ANSI_BOX_HORIZONTAL    = "─"
	ANSI_BOX_TOP_RIGHT     = "┐"
	ANSI_BOX_VERTICAL      = "│"
	ANSI_BOX_BOTTOM_LEFT   = "└"
	ANSI_BOX_BOTTOM_RIGHT  = "┘"
)

/*
Colorize
The following are the tags used for adding color in your text.

Foreground text tag: &
Tokens for foreground text are:

&x - Black				&r - Dark Red			&g - Dark Green
&y - Dark Yellow     	&b - Dark Blue			&p - Purple
&c - Cyan				&w - Grey
&R - Red				&G - Green				&Y - Yellow
&B - Blue				&P - Pink				&C - Light Cyan
&W - White

&u or &U - Underline the text.
&i or &I - Italicize the text.
&D - Resets to custom color for whatever is being displayed.
&d - Resets to terminal default color.

Blinking foreground text tag: }
Tokens for blinking text are:

}x - Black           	}r - Dark Red  		}g - Dark Green
}y - Dark Yellow		}b - Dark Blue 		}p - Purple
}c - Cyan            	}w - Grey
}R - Red             	}G - Green     		}Y - Yellow
}B - Blue            	}P - Pink      		}C - Light Blue
}W - White

When putting color in something, please try to remember to close your
colors with a &D tag so that anyone viewing it won't have to deal with
color bleeding all over the place. The same holds for italic or underlined
text as well.

The &d tag should only be used when absolutely necessary.

Background color tag: ^
Tokens for background color are:

^x - Black         ^r - Red           ^g - Green
^y - Yellow        ^b - Blue          ^p - Purple
^c - Cyan          ^w - Grey

If setting both foreground and background colors. The foreground must
be used before the background color. Also, the last color setting in your
prompt will wash over into the text you type. So, if you want a set
of colors for your typed text, include that at the end of your prompt set.

Example (assuming current h.p.'s of 43 and mana of 23):

Prompt &Y^b<%h/&x^r%m>&w^x = <43/23>

	{A}     {B}       {C}

A) Yellow with blue background.
B) Black with dark red background.
C) Light Grey with black background.
*/
type Colorize struct {
}

func Color() *Colorize {
	return &Colorize{}
}

func (c *Colorize) ClearScreen() string {
	return ANSI_CLEAR
}

func (c *Colorize) Reset() string {
	return ANSI_RESET
}

func (c *Colorize) IsUpArrow(input string) bool {
	return input == ANSI_CURSOR_UP
}

func (c *Colorize) IsDownArrow(input string) bool {
	return input == ANSI_CURSOR_DOWN
}

func (c *Colorize) IsLeftArrow(input string) bool {
	return input == ANSI_CURSOR_LEFT
}

func (c *Colorize) IsRightArrow(input string) bool {
	return input == ANSI_CURSOR_RIGHT
}

func (c *Colorize) Colorize(input string) string {
	// &=FG code
	// ^=BG code
	// }=Blink codes

	color_func := func(str string) string {
		switch str {
		case "&x":
			return ANSI_ESC + "0;" + ANSI_FG_BLACK + "m"
		case "&r":
			return ANSI_ESC + "0;" + ANSI_FG_RED + "m"
		case "&g":
			return ANSI_ESC + "0;" + ANSI_FG_GREEN + "m"
		case "&y":
			return ANSI_ESC + "0;" + ANSI_FG_YELLOW + "m"
		case "&b":
			return ANSI_ESC + "0;" + ANSI_FG_BLUE + "m"
		case "&p":
			return ANSI_ESC + "0;" + ANSI_FG_MAGENTA + "m"
		case "&c":
			return ANSI_ESC + "0;" + ANSI_FG_CYAN + "m"
		case "&w":
			return ANSI_ESC + "0;" + ANSI_FG_WHITE + "m"
		case "&X":
			return ANSI_BOLD + ANSI_FG_BLACK + "m"
		case "&R":
			return ANSI_BOLD + ANSI_FG_RED + "m"
		case "&G":
			return ANSI_BOLD + ANSI_FG_GREEN + "m"
		case "&Y":
			return ANSI_BOLD + ANSI_FG_YELLOW + "m"
		case "&B":
			return ANSI_BOLD + ANSI_FG_BLUE + "m"
		case "&P":
			return ANSI_BOLD + ANSI_FG_MAGENTA + "m"
		case "&C":
			return ANSI_BOLD + ANSI_FG_CYAN + "m"
		case "&W":
			return ANSI_BOLD + ANSI_FG_WHITE + "m"
		case "^x":
			return ANSI_ESC + ANSI_BG_BLACK + "m"
		case "^r":
			return ANSI_ESC + ANSI_BG_RED + "m"
		case "^g":
			return ANSI_ESC + ANSI_BG_GREEN + "m"
		case "^y":
			return ANSI_ESC + ANSI_BG_YELLOW + "m"
		case "^b":
			return ANSI_ESC + ANSI_BG_BLUE + "m"
		case "^p":
			return ANSI_ESC + ANSI_BG_MAGENTA + "m"
		case "^c":
			return ANSI_ESC + ANSI_BG_CYAN + "m"
		case "^w":
			return ANSI_ESC + ANSI_BG_WHITE + "m"
		case "^X":
			return ANSI_BOLD + ANSI_BG_BLACK + "m"
		case "^R":
			return ANSI_BOLD + ANSI_BG_RED + "m"
		case "^G":
			return ANSI_BOLD + ANSI_BG_GREEN + "m"
		case "^Y":
			return ANSI_BOLD + ANSI_BG_YELLOW + "m"
		case "^B":
			return ANSI_BOLD + ANSI_BG_BLUE + "m"
		case "^P":
			return ANSI_BOLD + ANSI_BG_MAGENTA + "m"
		case "^C":
			return ANSI_BOLD + ANSI_BG_CYAN + "m"
		case "^W":
			return ANSI_BOLD + ANSI_BG_WHITE + "m"
		case "}x":
			return ANSI_BLINK + ANSI_ESC + ANSI_FG_BLACK + "m"
		case "}r":
			return ANSI_BLINK + ANSI_ESC + ANSI_FG_RED + "m"
		case "}g":
			return ANSI_BLINK + ANSI_ESC + ANSI_FG_GREEN + "m"
		case "}y":
			return ANSI_BLINK + ANSI_ESC + ANSI_FG_YELLOW + "m"
		case "}b":
			return ANSI_BLINK + ANSI_ESC + ANSI_FG_BLUE + "m"
		case "}p":
			return ANSI_BLINK + ANSI_ESC + ANSI_FG_MAGENTA + "m"
		case "}c":
			return ANSI_BLINK + ANSI_ESC + ANSI_FG_CYAN + "m"
		case "}w":
			return ANSI_BLINK + ANSI_ESC + ANSI_FG_WHITE + "m"
		case "}X":
			return ANSI_BLINK + ANSI_BOLD + ANSI_FG_BLACK + "m"
		case "}R":
			return ANSI_BLINK + ANSI_BOLD + ANSI_FG_RED + "m"
		case "}G":
			return ANSI_BLINK + ANSI_BOLD + ANSI_FG_GREEN + "m"
		case "}Y":
			return ANSI_BLINK + ANSI_BOLD + ANSI_FG_YELLOW + "m"
		case "}B":
			return ANSI_BLINK + ANSI_BOLD + ANSI_FG_BLUE + "m"
		case "}P":
			return ANSI_BLINK + ANSI_BOLD + ANSI_FG_MAGENTA + "m"
		case "}C":
			return ANSI_BLINK + ANSI_BOLD + ANSI_FG_CYAN + "m"
		case "}W":
			return ANSI_BLINK + ANSI_BOLD + ANSI_FG_WHITE + "m"
		case "&U", "&u":
			return ANSI_UNDERLINE
		case "&I", "&i":
			return ANSI_ITALIC
		case "&&":
			return "&"
		case "^^":
			return "^"
		case "}}":
			return "}"
		case "&d", "&D":
			return ANSI_RESET
		default:
			return str
		}
	}

	r, err := regexp.Compile(`&[a-zA-z&]`) // foreground
	ErrorCheck(err)

	input = r.ReplaceAllStringFunc(input, color_func)

	r, err = regexp.Compile(`\^[a-zA-z^]`) // background
	ErrorCheck(err)

	input = r.ReplaceAllStringFunc(input, color_func)

	r, err = regexp.Compile(`}[a-zA-z}]`) // foreground
	ErrorCheck(err)

	input = r.ReplaceAllStringFunc(input, color_func)
	return input
}
func (c *Colorize) Decolorize(input string) string {
	// &=FG code
	// ^=BG code
	// }=Blink code
	// ]=Underline code

	color_func := func(str string) string {
		switch str {
		case "&d", "&D", "&x", "&r", "&g", "&y", "&b", "&p", "&c", "&w", "&X", "&R", "&G", "&Y", "&B", "&P", "&C", "&W", "&U", "&u", "&I", "&i", "^x", "^r", "^g", "^y", "^b", "^p", "^c", "^w", "^X", "^R", "^G", "^Y", "^B", "^P", "^C", "^W", "}x", "}r", "}g", "}y", "}b", "}p", "}c", "}w", "}X", "}R", "}G", "}Y", "}B", "}P", "}C", "}W":
			return ""
		case "&&":
			return "&"
		case "^^":
			return "^"
		case "}}":
			return "}"
		default:
			return str
		}
	}

	r, err := regexp.Compile(`&[a-zA-z&]`) // foreground
	ErrorCheck(err)

	input = r.ReplaceAllStringFunc(input, color_func)

	r, err = regexp.Compile(`\^[a-zA-z^]`) // background
	ErrorCheck(err)

	input = r.ReplaceAllStringFunc(input, color_func)

	r, err = regexp.Compile(`}[a-zA-z}]`) // blinky
	ErrorCheck(err)

	input = r.ReplaceAllStringFunc(input, color_func)
	return input
}

const (
	ANSI_TITLE_ALIGNMENT_LEFT = iota
	ANSI_TITLE_ALIGNMENT_CENTER
	ANSI_TITLE_ALIGNMENT_RIGHT
)

const (
	ANSI_TITLE_STYLE_SYSTEM = iota
	ANSI_TITLE_STYLE_NORMAL
	ANSI_TITLE_STYLE_BLOCK
	ANSI_TITLE_STYLE_ELEGANT
	ANSI_TITLE_STYLE_HACKED
	ANSI_TITLE_STYLE_IMPERIAL
	ANSI_TITLE_STYLE_REBEL
	ANSI_TITLE_STYLE_SENATE
)

// Takes a string and makes a Title based on style [ANSI_TITLE_STYLE_] and alignment [ANSI_TITLE_ALIGNMENT_]
func MakeTitle(title string, style int, alignment int) string {
	t := ""
	cap_left := ""
	cap_right := ""
	switch style {
	case ANSI_TITLE_STYLE_NORMAL:
		t = strings.Repeat("-=", 38) + "-"
		cap_left = "("
		cap_right = ")"
	case ANSI_TITLE_STYLE_BLOCK:
		t = strings.Repeat("==", 38) + "="
		cap_left = "["
		cap_right = "]"
	case ANSI_TITLE_STYLE_ELEGANT:
		t = strings.Repeat("-~", 38) + "-"
		cap_left = "{"
		cap_right = "}}"
	case ANSI_TITLE_STYLE_HACKED:
		t = strings.Repeat("-/\\#", 19) + "-"
		cap_left = "<"
		cap_right = ">"
	case ANSI_TITLE_STYLE_IMPERIAL:
		t = strings.Repeat("::", 39)
		cap_left = ":"
		cap_right = ":"
	case ANSI_TITLE_STYLE_REBEL:
		t = strings.Repeat("::", 39)
		cap_left = ":"
		cap_right = ":"
	case ANSI_TITLE_STYLE_SENATE:
		t = strings.Repeat("-=", 38) + "-"
		cap_left = "["
		cap_right = "]"
	default:
		t = "+----------------------------------------------------------------------------+"
		cap_left = "["
		cap_right = "]"
	}
	offset := 0
	title_length := len(title)

	switch alignment {
	case ANSI_TITLE_ALIGNMENT_CENTER:
		offset = (len(t) / 2) - ((title_length / 2) + 2)
	case ANSI_TITLE_ALIGNMENT_RIGHT:
		offset = (len(t) - (title_length + 2))
	default:
		offset = 2
	}
	text_color := "&W"
	title_color := "&g"
	if ANSI_TITLE_STYLE_IMPERIAL == style {
		title_color = "&B"
	}
	if ANSI_TITLE_STYLE_REBEL == style {
		title_color = "&R"
	}
	if ANSI_TITLE_STYLE_HACKED == style {
		title_color = "&P"
	}
	if ANSI_TITLE_STYLE_SENATE == style {
		title_color = "&C"
	}
	ret := fmt.Sprintf("%s%s%s %s%s %s%s%s&d\r\n", title_color, t[:offset], cap_left, text_color, title, title_color, cap_right, t[(offset+title_length):])

	return Color().Colorize(ret)
}

func MakeProgressBar(value int, max int, size int) string {
	percent := float64(value) / float64(max)
	size_percent := float64(size) * percent
	cap := int(math.Floor(size_percent))
	remainder := 0
	if value > 0 {
		remainder = max % value
	}

	ret := ""
	for i := 0; i < size; i++ {
		if i < cap {
			ret += "█"
		} else if i == cap {
			if remainder%2 == 1 {
				ret += "░"
			} else {
				ret += "▒"
			}
		} else {
			ret += "▪"
		}
	}
	ret += ""
	return ret
}
func MakeTunerBar(freq string, size int) string {
	max := 400.0
	value, err := strconv.ParseFloat(freq, 64)
	value -= 100.0 // bands the value to 0-400  (freq is 100-500mhz)
	ErrorCheck(err)
	if err != nil {
		return ""
	}
	percent := float64(value) / float64(max)
	size_percent := float64(size) * percent
	cap := int(math.Floor(size_percent))
	if cap >= size {
		cap = size - 1
	}
	ret := ""
	for i := 0; i < size; i++ {
		if i == cap {
			ret += "█"
		} else {
			switch i % 4 {
			case 0:
				ret += "&x.&d"
			case 1:
				ret += "&x*&d"
			case 2:
				ret += "&x'&d"
			case 3:
				ret += "&x*&d"
			}
		}
	}
	ret += ""
	return ret
}

func StitchParagraphs(paragraph1 string, paragraph2 string) string {
	p1_parts := strings.Split(paragraph1, "\r\n")

	p2_parts := strings.Split(paragraph2, "\r\n")
	p1_len := len(p1_parts)
	p2_len := len(p2_parts)

	tallest := p1_len
	if p2_len > p1_len {
		tallest = p2_len
	}

	buf := ""
	for row := 0; row < tallest; row++ {
		a_side := ""
		b_side := ""
		if row < p1_len {
			a_side = p1_parts[row]
		}
		if row < p2_len {
			b_side = p2_parts[row]
		}
		buf += sprintf("%-*s %-*s\r\n", 70, a_side, 10, b_side)
	}
	return buf
}

func tstring(str string, length int) string {
	if len(str) < length {
		return str
	} else {
		return str[:length]
	}
}
