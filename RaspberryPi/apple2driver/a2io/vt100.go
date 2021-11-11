// Copyright Terence J. Boldt (c)2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is contains VT100 terminal emulation

package a2io

import (
	"fmt"
	"time"
)

var escapeSequence string
var htab, vtab, savedHtab, savedVtab int
var applicationMode bool
var windowTop int
var windowBottom = 22

func sendCharacter(comm A2Io, b byte) {
	if b == 0x1b {
		escapeSequence = "^["
		return
	}
	if len(escapeSequence) > 0 {
		escapeSequence += string(b)
		// save cursor
		if escapeSequence == "^[7" {
			savedHtab = htab
			savedVtab = vtab
			escapeSequence = ""
		}
		// restore cursor
		if escapeSequence == "^[8" {
			htab = savedHtab
			vtab = savedVtab
			comm.WriteByte('H')
			comm.WriteByte(byte(htab))
			comm.WriteByte('V')
			comm.WriteByte(byte(vtab))
			escapeSequence = ""
		}
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') {
			switch b {
			// Set cursor location
			case 'H', 'f':
				var ignore string
				fmt.Sscanf(escapeSequence, "^[[%d;%d%s", &vtab, &htab, &ignore)
				htab--
				vtab--
				comm.WriteByte('H')
				comm.WriteByte(byte(htab))
				comm.WriteByte('V')
				comm.WriteByte(byte(vtab))
				escapeSequence = ""
			case 'r':
				fmt.Sscanf(escapeSequence, "^[[%d;%dr", &windowTop, &windowBottom)
				windowTop--
				//windowBottom--
				comm.WriteByte('T')
				comm.WriteByte(byte(windowTop))
				comm.WriteByte('B')
				comm.WriteByte(byte(windowBottom))
				escapeSequence = ""
			case 'C':
				var right int
				fmt.Sscanf(escapeSequence, "^[[%dC", &right)
				htab -= right
				for i := 0; i < right; i++ {
					comm.WriteByte(0x08)
				}
				escapeSequence = ""
			}
			switch escapeSequence {
			// Set/clear application mode for cursor
			case "^[[?1h":
				applicationMode = true
				comm.WriteByte(0x0c) // ^L clears the screen
				escapeSequence = ""
			case "^[[?1l":
				applicationMode = false
				comm.WriteByte('T')
				comm.WriteByte(0x00)
				comm.WriteByte('B')
				comm.WriteByte(0x18)
				comm.WriteByte(0x0c) // ^L clears the screen
				escapeSequence = ""
			// Tab to home position
			case "^[[H", "^[[;H", "^[[f", "^[[;f":
				htab = 0
				vtab = 0
				comm.WriteByte(0x19) // ^Y moves to home position
				escapeSequence = ""
			// Clear screen
			case "^[[2J", "^[[c":
				htab = 0
				vtab = 0
				comm.WriteByte(0x0c) // ^L clears the screen
				escapeSequence = ""
			// Move down one line
			case "^[E":
				comm.WriteByte(0x0A) // ^J moves cursor down
				escapeSequence = ""
			case "^[D":
				comm.WriteByte(0x17) // ^W scrolls up
				escapeSequence = ""
			// Clear line to the right
			case "^[[K", "^[[0K":
				comm.WriteByte(0x1d) // ^] clears to end of line
				escapeSequence = ""
			case "^[M":
				comm.WriteByte(0x16) // ^V scrolls down
				escapeSequence = ""
			// Clear screen below cursor
			case "^[[J":
				comm.WriteByte(0x0b) // ^K clears to end of screen
				escapeSequence = ""
			case "^[[7m":
				comm.WriteByte(0x0f) // ^O inverse video
				escapeSequence = ""
			case "^[[m", "^[[0m":
				comm.WriteByte(0x0e) // ^N normal video
				escapeSequence = ""
			}

			if len(escapeSequence) > 0 {
				fmt.Printf("\nUnhandled escape sequence: %s\n", escapeSequence)
			}
			escapeSequence = ""
			return
		}
		return
	}
	fmt.Print(string(b))
	htabIncrement := 0
	switch b {
	// convert LF to CR for Apple II compatiblity
	case 10:
		b = 13
		vtab++
		htab = 0
		htabIncrement = 0
	case 13:
		htab = 0
		htabIncrement = 0
		return
	// convert TAB to spaces
	case 9:
		b = ' '
		b |= 0x80
		err := comm.WriteByte(b)
		if err != nil {
			// try again because could have been cancelled by input
			time.Sleep(time.Millisecond * 10)
			comm.WriteByte(b)
		}
		htabIncrement = 2
	default:
		htabIncrement = 1
	}
	if !applicationMode || htab < 78 {
		updateTabs(comm)
		b |= 0x80
		htab += htabIncrement
		err := comm.WriteByte(b)
		if err != nil {
			// try again because could have been cancelled by input
			time.Sleep(time.Millisecond * 10)
			comm.WriteByte(b)
		}
	}
}

func updateTabs(comm A2Io) {
	if htab >= 80 {
		htab = 0
		vtab++
	}
	if vtab > windowBottom {
		vtab = windowBottom
	}
}
