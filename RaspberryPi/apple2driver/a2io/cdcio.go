// Copyright Terence J. Boldt (c)2020-2022
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is used for communicating with the Apple II data bus via the
// GPIO ports on the Raspberry Pi

package a2io

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

var port serial.Port

// CDCio is a live implementation of A2Io interface
type CDCio struct {
}

// Init initializes the CDC driver on the Raspberry Pi
func (a2 CDCio) Init() {
	name := ""
	for {
		portInfos, err := enumerator.GetDetailedPortsList()
		if err == nil {
			for _, portInfo := range portInfos {
				if portInfo.IsUSB && 
				        strings.ToUpper(portInfo.VID) == "2E8A" &&
					strings.ToUpper(portInfo.PID) == "000A" {
					name = portInfo.Name
					fmt.Printf("Found CDC port %s\n", name)
					break
				}
			}
			if name != "" {
				break
			}
		}
		time.Sleep(time.Millisecond)
	}

	for {
		var err error
		port, err = serial.Open(name, &serial.Mode{})
		if err == nil {
			break;
		}
		var portErr *serial.PortError
		if !errors.As(err, &portErr) || 
			portErr.Code() != serial.PortNotFound &&
			portErr.Code() != serial.PermissionDenied {
			panic(err)
		}
		time.Sleep(time.Millisecond)
	}

	err := port.SetReadTimeout(time.Second)
	if err != nil {
		panic(err)
	}
}

// ReadByte reads a byte from the Apple II via Raspberry Pi's CDC driver
func (a2 CDCio) ReadByte() (byte, error) {
	var data [1]byte
	n, err := port.Read(data[:]);
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, errors.New("timed out reading byte")
	}
	return data[0], nil
}

// WriteByte writes a byte to the Apple II via Raspberry Pi's CDC driver
func (a2 CDCio) WriteByte(data byte) error {
	_, err := port.Write([]byte{data});
	return err
}

// ReadString reads a string from the Apple II via Raspberry Pi's CDC driver
func (a2 CDCio) ReadString() (string, error) {
	var inBytes bytes.Buffer
	for {
		inByte, err := a2.ReadByte()
		if err != nil {
			return "", err
		}
		if inByte == 0 {
			break
		}
		inBytes.WriteByte(inByte)
	}
	return inBytes.String(), nil
}

// WriteString writes a string to the Apple II via Raspberry Pi's CDC driver
func (a2 CDCio) WriteString(outString string) error {
	for _, character := range outString {
		err := a2.WriteByte(byte(character) | 128)
		if err != nil {
			fmt.Printf("Failed to write string\n")
			return err
		}
	}
	a2.WriteByte(0)
	return nil
}

// WriteBlock writes 512 bytes to the Apple II via Raspberry Pi's CDC driver
func (a2 CDCio) WriteBlock(buffer []byte) error {
	_, err := port.Write(buffer);
	return err
}

// ReadBlock reads 512 bytes from the Apple II via Raspberry Pi's CDC driver
func (a2 CDCio) ReadBlock(buffer []byte) error {
	var err error
	for i := 0; i < 512; i++ {
		buffer[i], err = a2.ReadByte()
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteBuffer writes a buffer of bytes to the Apple II via Raspberry Pi's CDC driver
func (a2 CDCio) WriteBuffer(buffer []byte) error {
	_, err := port.Write(buffer);
	return err
}

// SendCharacter is a pass-through to vt100 implementation
func (a2 CDCio) SendCharacter(character byte) {
	sendCharacter(a2, character)
}

// ReadCharacter is a pass-through to vt100 implementation
func (a2 CDCio) ReadCharacter() (string, error) {
	return readCharacter(a2)
}
