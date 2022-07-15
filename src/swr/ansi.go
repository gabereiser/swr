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
	"regexp"
)

const (
	ANSI_ESC          = "\x1b["
	ANSI_CLEAR        = ANSI_ESC + "2J"
	ANSI_RESET        = ANSI_ESC + "0m"
	ANSI_CURSOR_UP    = ANSI_ESC + "2A"
	ANSI_CURSOR_DOWN  = ANSI_ESC + "2B"
	ANSI_CURSOR_LEFT  = ANSI_ESC + "2C"
	ANSI_CURSOR_RIGHT = ANSI_ESC + "2D"
	ANSI_BOLD         = ANSI_ESC + "1;"
	ANSI_ITALIC       = ANSI_ESC + "3;"
	ANSI_UNDERLINE    = ANSI_ESC + "4;"
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

/* The following are the tags used for adding color in your text.

Foreground text tag: &
Tokens for foreground text are:

&x - Black				&r - Dark Red			&g - Dark Green
&O - Orange (brown)		&b - Dark Blue			&p - Purple
&c - Cyan				&w - Grey				&z - Dark Grey
&R - Red				&G - Green				&Y - Yellow
&B - Blue				&P - Pink				&C - Light Blue
&W - White

&u or &U - Underline the text.
&i or &I - Italicize the text.
&s or &S - Strikeover text.
&D - Resets to custom color for whatever is being displayed.
&d - Resets to terminal default color.

Blinking foreground text tag: }
Tokens for blinking text are:

}x - Black           	}r - Dark Red  		}g - Dark Green
}O - Orange (brown)  	}b - Dark Blue 		}p - Purple
}c - Cyan            	}w - Grey      		}z - Dark Grey
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
^O - Orange        ^b - Blue          ^p - Purple
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
	// }=Blink code
	// ]=Underline code

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
		case "}W":
			return ANSI_BLINK + ANSI_BOLD + ANSI_FG_WHITE + "m"
		case "&&":
			return "&"
		case "^^":
			return "^"
		case "}}":
			return "}"
		case "&d":
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
