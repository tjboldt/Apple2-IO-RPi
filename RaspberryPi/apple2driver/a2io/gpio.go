// Copyright Terence J. Boldt (c)2020-2021
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

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

var edgeTimeout time.Duration

var outWrite gpio.PinIO
var outRead gpio.PinIO
var outReserved2 gpio.PinIO
var outReserved1 gpio.PinIO
var outBit7 gpio.PinIO
var outBit6 gpio.PinIO
var outBit5 gpio.PinIO
var outBit4 gpio.PinIO
var outBit3 gpio.PinIO
var outBit2 gpio.PinIO
var outBit1 gpio.PinIO
var outBit0 gpio.PinIO
var inWrite gpio.PinIO
var inRead gpio.PinIO
var inReserved2 gpio.PinIO
var inReserved1 gpio.PinIO
var inBit7 gpio.PinIO
var inBit6 gpio.PinIO
var inBit5 gpio.PinIO
var inBit4 gpio.PinIO
var inBit3 gpio.PinIO
var inBit2 gpio.PinIO
var inBit1 gpio.PinIO
var inBit0 gpio.PinIO

// A2Gpio is the live implementation of A2Io interface
type A2Gpio struct {
}

// Init initializes the GPIO ports on the Raspberry Pi
func (a2 A2Gpio) Init() {
	host.Init()

	outWrite = gpioreg.ByName("GPIO24")
	outRead = gpioreg.ByName("GPIO25")
	outReserved2 = gpioreg.ByName("GPIO7") //note GPIO7 and CPIO8 require extra effort to use
	outReserved1 = gpioreg.ByName("GPIO8")
	outBit7 = gpioreg.ByName("GPIO5")
	outBit6 = gpioreg.ByName("GPIO11")
	outBit5 = gpioreg.ByName("GPIO9")
	outBit4 = gpioreg.ByName("GPIO10")
	outBit3 = gpioreg.ByName("GPIO22")
	outBit2 = gpioreg.ByName("GPIO27")
	outBit1 = gpioreg.ByName("GPIO17")
	outBit0 = gpioreg.ByName("GPIO4")
	inWrite = gpioreg.ByName("GPIO23")
	inRead = gpioreg.ByName("GPIO18")
	inReserved2 = gpioreg.ByName("GPIO14")
	inReserved1 = gpioreg.ByName("GPIO15")
	inBit7 = gpioreg.ByName("GPIO12")
	inBit6 = gpioreg.ByName("GPIO16")
	inBit5 = gpioreg.ByName("GPIO20")
	inBit4 = gpioreg.ByName("GPIO21")
	inBit3 = gpioreg.ByName("GPIO26")
	inBit2 = gpioreg.ByName("GPIO19")
	inBit1 = gpioreg.ByName("GPIO13")
	inBit0 = gpioreg.ByName("GPIO6")

	inWrite.In(gpio.PullDown, gpio.BothEdges)
	inRead.In(gpio.PullDown, gpio.BothEdges)
	outReserved1.Out(gpio.High)
	outReserved2.Out(gpio.High)
	outRead.Out(gpio.High)
	outWrite.Out(gpio.High)
	outBit7.Out(gpio.Low)
	outBit6.Out(gpio.Low)
	outBit5.Out(gpio.Low)
	outBit4.Out(gpio.Low)
	outBit3.Out(gpio.Low)
	outBit2.Out(gpio.Low)
	outBit1.Out(gpio.Low)
	outBit0.Out(gpio.Low)

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
func (a2 A2Gpio) ReadByte() (byte, error) {
	// let the Apple II know we are ready to read
	outRead.Out(gpio.Low)

	// wait for the Apple II to write
	startTime := time.Now()
	for inWrite.Read() == gpio.High {
		if time.Since(startTime) > edgeTimeout {
			outRead.Out(gpio.High)
			return 0, errors.New("timed out reading byte -- write stuck high")
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

	if bit7 == gpio.High {
		data += 128
	}
	if bit6 == gpio.High {
		data += 64
	}
	if bit5 == gpio.High {
		data += 32
	}
	if bit4 == gpio.High {
		data += 16
	}
	if bit3 == gpio.High {
		data += 8
	}
	if bit2 == gpio.High {
		data += 4
	}
	if bit1 == gpio.High {
		data += 2
	}
	if bit0 == gpio.High {
		data++
	}

	// let the Apple II know we are done reading
	//fmt.Printf("let the Apple II know we are done reading\n")
	outRead.Out(gpio.High)

	// wait for the Apple II to finish writing
	//fmt.Printf("wait for the Apple II to finish writing\n")
	startTime = time.Now()
	for inWrite.Read() == gpio.Low {
		if time.Since(startTime) > edgeTimeout {
			return 0, errors.New("timed out reading byte -- write stuck low")
		}
	}

	return data, nil
}

// WriteByte writes a byte to the Apple II via Raspberry Pi's GPIO ports
func (a2 A2Gpio) WriteByte(data byte) error {
	// check if the Apple II wants to send a byte to us first
	if inWrite.Read() == gpio.Low {
		outWrite.Out(gpio.High)
		return errors.New("can't write byte while byte is incoming")
	}

	// wait for the Apple II to be ready to read
	startTime := time.Now()
	for inRead.Read() == gpio.High {
		if time.Since(startTime) > edgeTimeout {
			outWrite.Out(gpio.High)
			return errors.New("timed out writing byte -- read stuck high")
		}
	}

	bit7 := gpio.Low
	bit6 := gpio.Low
	bit5 := gpio.Low
	bit4 := gpio.Low
	bit3 := gpio.Low
	bit2 := gpio.Low
	bit1 := gpio.Low
	bit0 := gpio.Low

	if ((data & 128) >> 7) == 1 {
		bit7 = gpio.High
	}
	outBit7.Out(bit7)

	if ((data & 64) >> 6) == 1 {
		bit6 = gpio.High
	}
	outBit6.Out(bit6)

	if ((data & 32) >> 5) == 1 {
		bit5 = gpio.High
	}
	outBit5.Out(bit5)

	if ((data & 16) >> 4) == 1 {
		bit4 = gpio.High
	}
	outBit4.Out(bit4)

	if ((data & 8) >> 3) == 1 {
		bit3 = gpio.High
	}
	outBit3.Out(bit3)

	if ((data & 4) >> 2) == 1 {
		bit2 = gpio.High
	}
	outBit2.Out(bit2)

	if ((data & 2) >> 1) == 1 {
		bit1 = gpio.High
	}
	outBit1.Out(bit1)

	if (data & 1) == 1 {
		bit0 = gpio.High
	}
	outBit0.Out(bit0)

	// let Apple II know we're writing
	outWrite.Out(gpio.Low)

	// wait for the Apple II to finsih reading
	//fmt.Printf("wait for the Apple II to finsih reading\n")
	startTime = time.Now()
	for inRead.Read() == gpio.Low {
		if time.Since(startTime) > edgeTimeout {
			outWrite.Out(gpio.High)
			return errors.New("timed out writing byte -- read stuck low")
		}
	}

	// let the Apple II know we are done writing
	outWrite.Out(gpio.High)
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
