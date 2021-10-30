// Copyright Terence J. Boldt (c)2021
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

type UserIo struct {
}

func (userIo UserIo) Init() {

}

func (userIo UserIo) WriteByte(data byte) error {
	fmt.Printf("WriteByte: %02X\n", data)
	return nil
}

func (userIo UserIo) WriteString(outString string) error {
	fmt.Printf("WriteString: %s\n", strings.ReplaceAll(outString, "\r", "\n"))
	return nil
}

func (userIo UserIo) WriteBlock(buffer []byte) error {
	fmt.Printf("WriteBlock:\n")
	return userIo.WriteBuffer(buffer)
}

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

func (userIo UserIo) ReadByte() (byte, error) {
	fmt.Printf("ReadByte: ")
	var b byte
	fmt.Scanf("%x", &b)
	fmt.Printf("\n")
	return b, nil
}

func (userIo UserIo) ReadString() (string, error) {
	fmt.Printf("ReadString: ")
	var s string
	fmt.Scanf("%s", &s)
	fmt.Printf("\n")
	return s, nil
}

func (userIo UserIo) ReadBlock(buffer []byte) error {
	fmt.Printf("ReadBlock: (Not supported)")
	return errors.New("ReadBlock not supported")
}
