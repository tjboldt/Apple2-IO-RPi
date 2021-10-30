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

var out_write gpio.PinIO
var out_read gpio.PinIO
var out_reserved2 gpio.PinIO
var out_reserved1 gpio.PinIO
var out_bit7 gpio.PinIO
var out_bit6 gpio.PinIO
var out_bit5 gpio.PinIO
var out_bit4 gpio.PinIO
var out_bit3 gpio.PinIO
var out_bit2 gpio.PinIO
var out_bit1 gpio.PinIO
var out_bit0 gpio.PinIO
var in_write gpio.PinIO
var in_read gpio.PinIO
var in_reserved2 gpio.PinIO
var in_reserved1 gpio.PinIO
var in_bit7 gpio.PinIO
var in_bit6 gpio.PinIO
var in_bit5 gpio.PinIO
var in_bit4 gpio.PinIO
var in_bit3 gpio.PinIO
var in_bit2 gpio.PinIO
var in_bit1 gpio.PinIO
var in_bit0 gpio.PinIO

type A2Gpio struct {
}

func (a2 A2Gpio) Init() {
	host.Init()

	out_write = gpioreg.ByName("GPIO24")
	out_read = gpioreg.ByName("GPIO25")
	out_reserved2 = gpioreg.ByName("GPIO7") //note GPIO7 and CPIO8 require extra effort to use
	out_reserved1 = gpioreg.ByName("GPIO8")
	out_bit7 = gpioreg.ByName("GPIO5")
	out_bit6 = gpioreg.ByName("GPIO11")
	out_bit5 = gpioreg.ByName("GPIO9")
	out_bit4 = gpioreg.ByName("GPIO10")
	out_bit3 = gpioreg.ByName("GPIO22")
	out_bit2 = gpioreg.ByName("GPIO27")
	out_bit1 = gpioreg.ByName("GPIO17")
	out_bit0 = gpioreg.ByName("GPIO4")
	in_write = gpioreg.ByName("GPIO23")
	in_read = gpioreg.ByName("GPIO18")
	in_reserved2 = gpioreg.ByName("GPIO14")
	in_reserved1 = gpioreg.ByName("GPIO15")
	in_bit7 = gpioreg.ByName("GPIO12")
	in_bit6 = gpioreg.ByName("GPIO16")
	in_bit5 = gpioreg.ByName("GPIO20")
	in_bit4 = gpioreg.ByName("GPIO21")
	in_bit3 = gpioreg.ByName("GPIO26")
	in_bit2 = gpioreg.ByName("GPIO19")
	in_bit1 = gpioreg.ByName("GPIO13")
	in_bit0 = gpioreg.ByName("GPIO6")

	in_write.In(gpio.PullDown, gpio.BothEdges)
	in_read.In(gpio.PullDown, gpio.BothEdges)
	out_reserved1.Out(gpio.High)
	out_reserved2.Out(gpio.High)
	out_read.Out(gpio.High)
	out_write.Out(gpio.High)
	out_bit7.Out(gpio.Low)
	out_bit6.Out(gpio.Low)
	out_bit5.Out(gpio.Low)
	out_bit4.Out(gpio.Low)
	out_bit3.Out(gpio.Low)
	out_bit2.Out(gpio.Low)
	out_bit1.Out(gpio.Low)
	out_bit0.Out(gpio.Low)

	edgeTimeout = time.Second
}

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
	return string(inBytes.Bytes()), nil
}

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

func (a2 A2Gpio) ReadByte() (byte, error) {
	// let the Apple II know we are ready to read
	out_read.Out(gpio.Low)

	// wait for the Apple II to write
	for in_write.Read() == gpio.High {
		if !in_write.WaitForEdge(edgeTimeout) {
			out_read.Out(gpio.High)
			return 0, errors.New("Timed out reading byte -- write stuck high\n")
		}
	}

	// get a nibble of data
	var data byte
	data = 0
	bit7 := in_bit7.Read()
	bit6 := in_bit6.Read()
	bit5 := in_bit5.Read()
	bit4 := in_bit4.Read()
	bit3 := in_bit3.Read()
	bit2 := in_bit2.Read()
	bit1 := in_bit1.Read()
	bit0 := in_bit0.Read()

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
		data += 1
	}

	// let the Apple II know we are done reading
	//fmt.Printf("let the Apple II know we are done reading\n")
	out_read.Out(gpio.High)

	// wait for the Apple II to finish writing
	//fmt.Printf("wait for the Apple II to finish writing\n")
	for in_write.Read() == gpio.Low {
		if !in_write.WaitForEdge(edgeTimeout) {
			return 0, errors.New("Timed out reading byte -- write stuck low")
		}
	}

	return data, nil
}

func (a2 A2Gpio) WriteByte(data byte) error {
	// check if the Apple II wants to send a byte to us first
	if in_write.Read() == gpio.Low {
		out_write.Out(gpio.High)
		return errors.New("Can't write byte while byte is incoming")
	}

	// wait for the Apple II to be ready to read
	for in_read.Read() == gpio.High {
		if !in_read.WaitForEdge(edgeTimeout) {
			out_write.Out(gpio.High)
			return errors.New("Timed out writing byte -- read stuck high")
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
	out_bit7.Out(bit7)

	if ((data & 64) >> 6) == 1 {
		bit6 = gpio.High
	}
	out_bit6.Out(bit6)

	if ((data & 32) >> 5) == 1 {
		bit5 = gpio.High
	}
	out_bit5.Out(bit5)

	if ((data & 16) >> 4) == 1 {
		bit4 = gpio.High
	}
	out_bit4.Out(bit4)

	if ((data & 8) >> 3) == 1 {
		bit3 = gpio.High
	}
	out_bit3.Out(bit3)

	if ((data & 4) >> 2) == 1 {
		bit2 = gpio.High
	}
	out_bit2.Out(bit2)

	if ((data & 2) >> 1) == 1 {
		bit1 = gpio.High
	}
	out_bit1.Out(bit1)

	if (data & 1) == 1 {
		bit0 = gpio.High
	}
	out_bit0.Out(bit0)

	// let Apple II know we're writing
	out_write.Out(gpio.Low)

	// wait for the Apple II to finsih reading
	//fmt.Printf("wait for the Apple II to finsih reading\n")
	for in_read.Read() == gpio.Low {
		if !in_read.WaitForEdge(edgeTimeout) {
			out_write.Out(gpio.High)
			return errors.New("Timed out writing byte -- read stuck low")
		}
	}

	// let the Apple II know we are done writing
	out_write.Out(gpio.High)
	return nil
}

func (a2 A2Gpio) WriteBlock(buffer []byte) error {
	for i := 0; i < 512; i++ {
		err := a2.WriteByte(buffer[i])
		if err != nil {
			return err
		}
	}

	return nil
}

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
