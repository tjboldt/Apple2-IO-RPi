// Copyright Terence J. Boldt (c)2020-2023
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is used for communicating with the Apple II data bus via the
// GPIO ports on the Raspberry Pi

package a2io

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

var edgeTimeout time.Duration

var outWrite rpio.Pin
var outRead rpio.Pin
var outReserved2 rpio.Pin
var outReserved1 rpio.Pin
var outBit7 rpio.Pin
var outBit6 rpio.Pin
var outBit5 rpio.Pin
var outBit4 rpio.Pin
var outBit3 rpio.Pin
var outBit2 rpio.Pin
var outBit1 rpio.Pin
var outBit0 rpio.Pin
var inWrite rpio.Pin
var inRead rpio.Pin
var inReserved2 rpio.Pin
var inReserved1 rpio.Pin
var inBit7 rpio.Pin
var inBit6 rpio.Pin
var inBit5 rpio.Pin
var inBit4 rpio.Pin
var inBit3 rpio.Pin
var inBit2 rpio.Pin
var inBit1 rpio.Pin
var inBit0 rpio.Pin

// A2Gpio is the live implementation of A2Io interface
type A2Gpio struct {
}

// Init initializes the GPIO ports on the Raspberry Pi
func (a2 A2Gpio) Init() {
	err := rpio.Open()
	if err != nil {
		panic("failed to open gpio")
	}

	outWrite = rpio.Pin(24)
	outRead = rpio.Pin(25)
	outReserved2 = rpio.Pin(7) //note GPIO7 and CPIO8 require extra effort to use
	outReserved1 = rpio.Pin(8)
	outBit7 = rpio.Pin(5)
	outBit6 = rpio.Pin(11)
	outBit5 = rpio.Pin(9)
	outBit4 = rpio.Pin(10)
	outBit3 = rpio.Pin(22)
	outBit2 = rpio.Pin(27)
	outBit1 = rpio.Pin(17)
	outBit0 = rpio.Pin(4)
	inWrite = rpio.Pin(23)
	inRead = rpio.Pin(18)
	inReserved2 = rpio.Pin(14)
	inReserved1 = rpio.Pin(15)
	inBit7 = rpio.Pin(12)
	inBit6 = rpio.Pin(16)
	inBit5 = rpio.Pin(20)
	inBit4 = rpio.Pin(21)
	inBit3 = rpio.Pin(26)
	inBit2 = rpio.Pin(19)
	inBit1 = rpio.Pin(13)
	inBit0 = rpio.Pin(6)

	inWrite.PullDown()
	inRead.PullDown()

	outWrite.Output()
	outRead.Output()
	outReserved2.Output()
	outReserved1.Output()
	outBit7.Output()
	outBit6.Output()
	outBit5.Output()
	outBit4.Output()
	outBit3.Output()
	outBit2.Output()
	outBit1.Output()
	outBit0.Output()

	inWrite.Input()
	inRead.Input()
	inReserved2.Input()
	inReserved1.Input()
	inBit7.Input()
	inBit6.Input()
	inBit5.Input()
	inBit4.Input()
	inBit3.Input()
	inBit2.Input()
	inBit1.Input()
	inBit0.Input()

	outReserved1.High()
	outReserved2.High()
	outRead.High()
	outWrite.High()
	outBit7.Low()
	outBit6.Low()
	outBit5.Low()
	outBit4.Low()
	outBit3.Low()
	outBit2.Low()
	outBit1.Low()
	outBit0.Low()

	edgeTimeout = time.Second
}

// ReadString reads a string from the Apple II via Raspberry Pi's GPIO ports
func (a2 A2Gpio) ReadString() (string, error) {
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

// WriteString writes a string to the Apple II via Raspberry Pi's GPIO ports
func (a2 A2Gpio) WriteString(outString string) error {
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

// ReadByte reads a byte from the Apple II via Raspberry Pi's GPIO ports
func (a2 A2Gpio) ReadByte(noDelay ...bool) (byte, error) {
	// let the Apple II know we are ready to read
	outRead.Low()

	// wait for the Apple II to write
	startTime := time.Now()
	lastSleepTime := time.Now()
	sleepDuration := 10
	for inWrite.Read() == 1 {
		if time.Since(startTime) > edgeTimeout {
			outRead.High()
			return 0, errors.New("timed out reading byte -- write stuck high")
		}
		if time.Since(lastSleepTime) > edgeTimeout/10 {
			if len(noDelay) == 0 || !noDelay[0] {
				sleepDuration *= 3;
			}
			time.Sleep(time.Millisecond * time.Duration(sleepDuration));
			lastSleepTime = time.Now()
		}
	}

	// get a nibble of data
	var data byte
	data = 0
	bit7 := inBit7.Read()
	bit6 := inBit6.Read()
	bit5 := inBit5.Read()
	bit4 := inBit4.Read()
	bit3 := inBit3.Read()
	bit2 := inBit2.Read()
	bit1 := inBit1.Read()
	bit0 := inBit0.Read()

	if bit7 == 1 {
		data += 128
	}
	if bit6 == 1 {
		data += 64
	}
	if bit5 == 1 {
		data += 32
	}
	if bit4 == 1 {
		data += 16
	}
	if bit3 == 1 {
		data += 8
	}
	if bit2 == 1 {
		data += 4
	}
	if bit1 == 1 {
		data += 2
	}
	if bit0 == 1 {
		data++
	}

	// let the Apple II know we are done reading
	//fmt.Printf("let the Apple II know we are done reading\n")
	outRead.High()

	// wait for the Apple II to finish writing
	//fmt.Printf("wait for the Apple II to finish writing\n")
	startTime = time.Now()
	lastSleepTime = time.Now()
	sleepDuration = 10
	for inWrite.Read() == 0 {
		if time.Since(startTime) > edgeTimeout {
			return 0, errors.New("timed out reading byte -- write stuck low")
		}
		if time.Since(lastSleepTime) > edgeTimeout/10 {
			if len(noDelay) == 0 || !noDelay[0] {
				sleepDuration *= 3;
			}
			time.Sleep(time.Millisecond * time.Duration(sleepDuration));
			lastSleepTime = time.Now()
		}
	}

	return data, nil
}

// WriteByte writes a byte to the Apple II via Raspberry Pi's GPIO ports
func (a2 A2Gpio) WriteByte(data byte) error {
	// check if the Apple II wants to send a byte to us first
	if inWrite.Read() == 0 {
		outWrite.High()
		return errors.New("can't write byte while byte is incoming")
	}

	// wait for the Apple II to be ready to read
	startTime := time.Now()
	lastSleepTime := time.Now()
	sleepDuration := 10
	for inRead.Read() == 1 {
		if time.Since(startTime) > edgeTimeout {
			outWrite.High()
			return errors.New("timed out writing byte -- read stuck high")
		}
		if time.Since(lastSleepTime) > edgeTimeout/10 {
			sleepDuration *= 3;
			time.Sleep(time.Millisecond * time.Duration(sleepDuration));
			lastSleepTime = time.Now()
		}
	}

	if ((data & 128) >> 7) == 1 {
		outBit7.High()
	} else {
		outBit7.Low()
	}

	if ((data & 64) >> 6) == 1 {
		outBit6.High()
	} else {
		outBit6.Low()
	}

	if ((data & 32) >> 5) == 1 {
		outBit5.High()
	} else {
		outBit5.Low()
	}

	if ((data & 16) >> 4) == 1 {
		outBit4.High()
	} else {
		outBit4.Low()
	}

	if ((data & 8) >> 3) == 1 {
		outBit3.High()
	} else {
		outBit3.Low()
	}

	if ((data & 4) >> 2) == 1 {
		outBit2.High()
	} else {
		outBit2.Low()
	}

	if ((data & 2) >> 1) == 1 {
		outBit1.High()
	} else {
		outBit1.Low()
	}

	if (data & 1) == 1 {
		outBit0.High()
	} else {
		outBit0.Low()
	}

	// let Apple II know we're writing
	outWrite.Low()

	// wait for the Apple II to finsih reading
	//fmt.Printf("wait for the Apple II to finsih reading\n")
	startTime = time.Now()
	lastSleepTime = time.Now()
	sleepDuration = 10
	for inRead.Read() == 0 {
		if time.Since(startTime) > edgeTimeout {
			outWrite.High()
			return errors.New("timed out writing byte -- read stuck low")
		}
		if time.Since(lastSleepTime) > edgeTimeout/10 {
			sleepDuration *= 3;
			time.Sleep(time.Millisecond * time.Duration(sleepDuration));
			lastSleepTime = time.Now()
		}
	}

	// let the Apple II know we are done writing
	outWrite.High()
	return nil
}

// WriteBlock writes 512 bytes to the Apple II via Raspberry Pi's GPIO ports
func (a2 A2Gpio) WriteBlock(buffer []byte) error {
	for i := 0; i < 512; i++ {
		err := a2.WriteByte(buffer[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadBlock reads 512 bytes from the Apple II via Raspberry Pi's GPIO ports
func (a2 A2Gpio) ReadBlock(buffer []byte) error {
	var err error
	for i := 0; i < 512; i++ {
		buffer[i], err = a2.ReadByte()
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteBuffer writes a buffer of bytes to the Apple II via Raspberry Pi's GPIO ports
func (a2 A2Gpio) WriteBuffer(buffer []byte) error {
	bufferSize := len(buffer)
	for i := 0; i < bufferSize; i++ {
		err := a2.WriteByte(buffer[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// SendCharacter is a pass-through to vt100 implementation
func (a2 A2Gpio) SendCharacter(character byte) {
	sendCharacter(a2, character)
}

// ReadCharacter is a pass-through to vt100 implementation
func (a2 A2Gpio) ReadCharacter() (string, error) {
	return readCharacter(a2)
}
