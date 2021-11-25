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
var operatingSystemSequence bool
var htab, vtab, savedHtab, savedVtab int
var applicationMode bool
var windowTop int
var windowBottom = 23

func sendCharacter(comm A2Io, b byte) {
	if b == 0x1b {
		escapeSequence = "^["
		return
	}
	if b == 13 {
		fmt.Printf("CR\n")
		comm.WriteByte('H')
		comm.WriteByte(0)
		return
	}
	if b > 13 && b < 32 {
		fmt.Printf("Control code: %02X\n", b)
		return
	}
	if len(escapeSequence) == 0 {
		fmt.Printf("%c", b)
	}
	if len(escapeSequence) > 0 {
		escapeSequence += string(b)
		if escapeSequence == "^[]" {
			fmt.Printf("Start operating system sequence\n")
			operatingSystemSequence = true
			return
		}
		if operatingSystemSequence {
			if b == 0x07 {
				fmt.Printf("Operating system sequence: %s\n", escapeSequence)
				operatingSystemSequence = false
				escapeSequence = ""
			}
			return
		}
		// save cursor
		if escapeSequence == "^[7" {
			savedHtab = htab
			savedVtab = vtab
			fmt.Printf("Save Cursor (%d, %d): %s\n", htab, vtab, escapeSequence)
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
			fmt.Printf("Restore Cursor (%d, %d): %s\n", htab, vtab, escapeSequence)
			escapeSequence = ""
		}
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') {
			switch b {
			// Set cursor location
			case 'H', 'f':
				if escapeSequence == "^[[H" || escapeSequence ==  "^[[;H" {
					htab = 0
					vtab = 0
					comm.WriteByte(0x19) // ^Y moves to home position
					fmt.Printf("Home: %s\n", escapeSequence)
					escapeSequence = ""
				} else {
					var ignore string
					fmt.Sscanf(escapeSequence, "^[[%d;%d%s", &vtab, &htab, &ignore)
					htab--
					vtab--
					if htab < 0 {
						htab = 0
					}
					comm.WriteByte('H')
					comm.WriteByte(byte(htab))
					comm.WriteByte('V')
					comm.WriteByte(byte(vtab))
					fmt.Printf("Set Cursor (%d, %d): %s\n", htab, vtab, escapeSequence)
					escapeSequence = ""
				}
			case 'r':
				fmt.Sscanf(escapeSequence, "^[[%d;%dr", &windowTop, &windowBottom)
				windowTop--
				//windowBottom--
				comm.WriteByte('T')
				comm.WriteByte(byte(windowTop))
				comm.WriteByte('B')
				comm.WriteByte(byte(windowBottom))
				fmt.Printf("Set Window (%d, %d): %s\n", windowTop, windowBottom, escapeSequence)
				escapeSequence = ""
			case 'A':
				if escapeSequence == "^[[A" || escapeSequence == "^[A" {
					vtab--
					comm.WriteByte('U')
					fmt.Printf("Up: %s\n", escapeSequence)
				} else {
					var up int
					fmt.Sscanf(escapeSequence, "^[[%dA", &up)
					vtab -= up
					for i := 0; i < up; i++ {
						comm.WriteByte('U')
					}
					fmt.Printf("Up (%d): %s\n", up, escapeSequence)
				}
				escapeSequence = ""
			case 'B':
				if escapeSequence == "^[(B" || escapeSequence == "^[)B" || escapeSequence == "^[B" {
					escapeSequence = ""
				} else if escapeSequence == "^[[B" {
					vtab++
					comm.WriteByte(0x0a)
					fmt.Printf("Down: %s\n", escapeSequence)
				} else {
					var down int
					fmt.Sscanf(escapeSequence, "^[[%dB", &down)
					vtab += down
					for i := 0; i < down; i++ {
						comm.WriteByte(0x0a)
					}
					fmt.Printf("Down (%d): %s\n", down, escapeSequence)
				}
				escapeSequence = ""
			case 'C':
				if escapeSequence == "^[[C" || escapeSequence == "^[C" {
					htab++
					comm.WriteByte(0x1c)
					fmt.Printf("Right: %s\n", escapeSequence)
				} else {
					var right int
					fmt.Sscanf(escapeSequence, "^[[%dC", &right)
					htab += right
					for i := 0; i < right; i++ {
						comm.WriteByte(0x1c)
					}
					fmt.Printf("Right (%d): %s\n", right, escapeSequence)
				}
				escapeSequence = ""
			case 'D':
				if escapeSequence == "^[[D" || escapeSequence == "^[D" {
					htab--
					comm.WriteByte(0x08)
					fmt.Printf("Left: %s\n", escapeSequence)
				} else {
					var left int
					fmt.Sscanf(escapeSequence, "^[[%dD", &left)
					htab -= left
					for i := 0; i < left; i++ {
						comm.WriteByte(0x08)
					}
					fmt.Printf("Left (%d): %s\n", left, escapeSequence)
				}
				escapeSequence = ""
			}
			switch escapeSequence {
			// Set/clear application mode for cursor
			case "^[[?1h":
				applicationMode = true
				comm.WriteByte(0x0c) // ^L clears the screen
				fmt.Printf("Start application mode: %s\n", escapeSequence)
				escapeSequence = ""
			case "^[[?1l":
				applicationMode = false
				comm.WriteByte('T')
				comm.WriteByte(0x00)
				comm.WriteByte('B')
				comm.WriteByte(0x18)
				comm.WriteByte(0x0c) // ^L clears the screen
				fmt.Printf("End application mode: %s\n", escapeSequence)
				escapeSequence = ""
			// Tab to home position
			case "^[[f", "^[[;f":
				htab = 0
				vtab = 0
				comm.WriteByte(0x19) // ^Y moves to home position
				fmt.Printf("Home: %s\n", escapeSequence)
				escapeSequence = ""
			// Clear screen
			case "^[[2J", "^[[c":
				htab = 0
				vtab = 0
				comm.WriteByte(0x0c) // ^L clears the screen
				fmt.Printf("Clear screen: %s\n", escapeSequence)
				escapeSequence = ""
			case "^[E":
				comm.WriteByte(0x0A) // ^J moves cursor down
				fmt.Printf("Move down: %s\n", escapeSequence)
				escapeSequence = ""
			case "^[D":
				comm.WriteByte(0x17) // ^W scrolls up
				fmt.Printf("Scroll up: %s\n", escapeSequence)
				escapeSequence = ""
			case "^[[K", "^[[0K":
				comm.WriteByte(0x1d) // ^] clears to end of line
				fmt.Printf("Clear line right: %s\n", escapeSequence)
				escapeSequence = ""
			case "^[[2K":
				comm.WriteByte(0x1a) // ^Z clears line
				fmt.Printf("Clear line: %s\n", escapeSequence)
				escapeSequence = ""
			case "^[M":
				comm.WriteByte(0x16) // ^V scrolls down
				fmt.Printf("Scroll down: %s\n", escapeSequence)
				escapeSequence = ""
			case "^[[J":
				comm.WriteByte(0x0b) // ^K clears to end of screen
				fmt.Printf("Clear below cursor: %s\n", escapeSequence)
				escapeSequence = ""
			case "^[[7m":
				comm.WriteByte(0x0f) // ^O inverse video
				fmt.Printf("Inverse: %s\n", escapeSequence)
				escapeSequence = ""
			case "^[[m", "^[[0m":
				//comm.WriteByte(0x0e) // ^N normal video
				fmt.Printf("Normal: %s\n", escapeSequence)
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
	//fmt.Print(string(b))
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
	if !applicationMode || htab < 79 {
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

func updateTabs(comm A2Io) {
	if htab >= 80 {
		htab = 0
		vtab++
	}
	if vtab > windowBottom {
		vtab = windowBottom
	}
}
