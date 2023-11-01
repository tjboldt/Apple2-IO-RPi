// Copyright Terence J. Boldt (c)2021-2022
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is only used during development for interactive testing of the
// driver code. Simply replace `comm := a2io.A2Gpio{}` with
// `comm := a2io.UserIo{}` in the driver.go file to allow local testing

package a2io

import (
	"errors"
	"fmt"
	"strings"
)

// UserIo implements A2Io for the purpose of debugging locally
type UserIo struct {
}

// Init is only here to complete A2Io interface
func (userIo UserIo) Init() {

}

// WriteByte simulates writing to the Apple II but uses stdout instead
func (userIo UserIo) WriteByte(data byte) error {
	fmt.Printf("WriteByte: %02X\n", data)
	return nil
}

// WriteString simulates writing to the Apple II but uses stdout instead
func (userIo UserIo) WriteString(outString string) error {
	fmt.Printf("WriteString: %s\n", strings.ReplaceAll(outString, "\r", "\n"))
	return nil
}

// WriteBlock simulates writing to the Apple II but uses stdout instead
func (userIo UserIo) WriteBlock(buffer []byte) error {
	fmt.Printf("WriteBlock:\n")
	return userIo.WriteBuffer(buffer)
}

// WriteBuffer simulates writing to the Apple II but uses stdout instead
func (userIo UserIo) WriteBuffer(buffer []byte) error {
	fmt.Printf("WriteBuffer:\n")
	for i, b := range buffer {
		if i%16 == 0 {
			fmt.Printf("\n%04X:", i)
		}
		fmt.Printf(" %02X", b)
	}
	fmt.Printf("\n")
	return nil
}

// ReadByte simulates reading to the Apple II but uses stdin instead
func (userIo UserIo) ReadByte(noDelay ...bool) (byte, error) {
	fmt.Printf("ReadByte: ")
	var b byte
	fmt.Scanf("%x", &b)
	fmt.Printf("\n")
	return b, nil
}

// ReadString simulates reading to the Apple II but uses stdin instead
func (userIo UserIo) ReadString() (string, error) {
	fmt.Printf("ReadString: ")
	var s string
	fmt.Scanf("%s", &s)
	fmt.Printf("\n")
	return s, nil
}

// ReadBlock should simulate reading to the Apple II but is not yet supported
func (userIo UserIo) ReadBlock(buffer []byte) error {
	fmt.Printf("ReadBlock: (Not supported)")
	return errors.New("ReadBlock not supported")
}

// SendCharacter is a pass-through to vt100 implementation
func (userIo UserIo) SendCharacter(character byte) {
	sendCharacter(userIo, character)
}

// ReadCharacter is a pass-through to vt100 implementation
func (userIo UserIo) ReadCharacter() (string, error) {
	return readCharacter(userIo)
}
