// Copyright Terence J. Boldt (c)2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is contains VT100 terminal emulation

package a2io

import (
	"fmt"
)

var escapeSequence string
var operatingSystemSequence bool
var windowTop int
var windowBottom = 23

func sendCharacter(comm A2Io, b byte) {
	if b == 0x1b {
		escapeSequence = "^["
		return
	}
	if b == 0x0d {
		comm.WriteByte('H')
		comm.WriteByte(0)
		return
	}
	if b > 0x0d && b < 0x20 || b > 0x80 && b < 0xa0 {
		return
	}
	if b >= 0xa0 {
		comm.WriteByte('+' | 0x80)
		return
	}
	if len(escapeSequence) > 0 {
		// Parse the escape codes that don't end with a letter
		escapeSequence += string(b)
		if escapeSequence == "^[]" {
			operatingSystemSequence = true
			return
		}
		if operatingSystemSequence {
			if b == 0x07 {
				operatingSystemSequence = false
				escapeSequence = ""
			}
			return
		}
		// save cursor
		if escapeSequence == "^[7" {
			escapeSequence = ""
		}
		// restore cursor
		if escapeSequence == "^[8" {
			escapeSequence = ""
		}
		if escapeSequence == "^[>" {
			// this sequence is undocumented and shows up exiting "top" command
			escapeSequence = ""
		}
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') {
			// Parse simple escape codes
			switch escapeSequence {
			// Set/clear application mode for cursor
			case "^[[?1h":
				//applicationMode = true
				escapeSequence = ""
				return
			case "^[[?1l":
				//applicationMode = false
				comm.WriteByte('T')
				comm.WriteByte(0x00)
				comm.WriteByte('B')
				comm.WriteByte(0x18)
				escapeSequence = ""
				return
			// Tab to home position
			case "^[[f", "^[[;f":
				comm.WriteByte(0x19) // ^Y moves to home position
				escapeSequence = ""
				return
			// Clear screen
			case "^[[2J", "^[[c":
				comm.WriteByte(0x0c) // ^L clears the screen
				escapeSequence = ""
				return
			case "^[D":
				comm.WriteByte(0x17) // ^W scrolls up
				escapeSequence = ""
				return
			case "^[[K", "^[[0K":
				comm.WriteByte(0x1d) // ^] clears to end of line
				escapeSequence = ""
				return
			case "^[[2K":
				comm.WriteByte(0x1a) // ^Z clears line
				escapeSequence = ""
				return
			case "^[M":
				comm.WriteByte(0x16) // ^V scrolls down
				escapeSequence = ""
				return
			case "^[[J":
				comm.WriteByte(0x0b) // ^K clears to end of screen
				escapeSequence = ""
				return
			case "^[[7m":
				comm.WriteByte(0x0f) // ^O inverse video
				escapeSequence = ""
				return
			case "^[[m", "^[[0m", "^[[0;7m", "^[[0;1m":
				comm.WriteByte(0x0e) // ^N normal video
				escapeSequence = ""
				return
			}

			// Parse escape codes that need further parsing
			switch b {
			// Set cursor location
			case 'H', 'f':
				if escapeSequence == "^[[H" || escapeSequence == "^[[;H" {
					comm.WriteByte(0x19) // ^Y moves to home position
					escapeSequence = ""
				} else {
					var ignore string
					var htab, vtab int
					fmt.Sscanf(escapeSequence, "^[[%d;%d%s", &vtab, &htab, &ignore)
					htab--
					vtab--
					if htab < 0 { // this occastionally gets called with 0 that becomes -1
						htab = 0
					}
					if vtab > 23 { // top command sets vtab 25 on exit even in 24 line mode
						vtab = 23
					}
					comm.WriteByte('H')
					comm.WriteByte(byte(htab))
					comm.WriteByte('V')
					comm.WriteByte(byte(vtab))
					escapeSequence = ""
				}
			case 'r':
				fmt.Sscanf(escapeSequence, "^[[%d;%dr", &windowTop, &windowBottom)
				windowTop--
				comm.WriteByte('T')
				comm.WriteByte(byte(windowTop))
				comm.WriteByte('B')
				comm.WriteByte(byte(windowBottom))
				escapeSequence = ""
			case 'A':
				if escapeSequence == "^[[A" || escapeSequence == "^[A" {
					comm.WriteByte('U')
				} else {
					var up int
					fmt.Sscanf(escapeSequence, "^[[%dA", &up)
					for i := 0; i < up; i++ {
						comm.WriteByte('U')
					}
				}
				escapeSequence = ""
			case 'B':
				if escapeSequence == "^[(B" || escapeSequence == "^[)B" || escapeSequence == "^[B" {
					escapeSequence = ""
				} else if escapeSequence == "^[[B" {
					comm.WriteByte(0x0a)
				} else {
					var down int
					fmt.Sscanf(escapeSequence, "^[[%dB", &down)
					for i := 0; i < down; i++ {
						comm.WriteByte(0x0a)
					}
				}
				escapeSequence = ""
			case 'C':
				if escapeSequence == "^[[C" || escapeSequence == "^[C" {
					comm.WriteByte(0x1c)
				} else {
					var right int
					fmt.Sscanf(escapeSequence, "^[[%dC", &right)
					for i := 0; i < right; i++ {
						comm.WriteByte(0x1c)
					}
				}
				escapeSequence = ""
			case 'D':
				if escapeSequence == "^[[D" {
					comm.WriteByte(0x08)
				} else {
					var left int
					fmt.Sscanf(escapeSequence, "^[[%dD", &left)
					for i := 0; i < left; i++ {
						comm.WriteByte(0x08)
					}
				}
				escapeSequence = ""
			}
			if len(escapeSequence) > 0 {
				fmt.Printf("Unhandled escape sequence: %s\n", escapeSequence)
			}
			escapeSequence = ""
			return
		}
		return
	}
	switch b {
	// convert LF to CR for Apple II compatiblity
	case 10:
		b = 13
	case 13:
		return
	// convert TAB to spaces
	case 9:
		b = ' '
		b |= 0x80
		comm.WriteByte(b)
	}
	b |= 0x80
	comm.WriteByte(b)
}

func readCharacter(comm A2Io) (string, error) {
	b, err := comm.ReadByte()
	var s = string(b)
	if err == nil {
		switch b {
		case 0x0b: // up
			s = "\033[A"
		case 0x0a: // down
			s = "\033[B"
		case 0x15: // right
			s = "\033[C"
		case 0x08: // left
			s = "\033[D"
		case 0x0d: // return
			s = string(byte(0x0a))
		}
	}
	return s, err
}
