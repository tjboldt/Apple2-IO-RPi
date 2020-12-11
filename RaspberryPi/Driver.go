package main

import (
	//"time"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	"fmt"
	"os"
)

var out_write gpio.PinIO
var out_read gpio.PinIO
var out_commandWrite gpio.PinIO
var out_commandRead gpio.PinIO
var out_bit3 gpio.PinIO
var out_bit2 gpio.PinIO
var out_bit1 gpio.PinIO
var out_bit0 gpio.PinIO
var in_write gpio.PinIO
var in_read gpio.PinIO
var in_commandWrite gpio.PinIO
var in_commandRead gpio.PinIO
var in_bit3 gpio.PinIO
var in_bit2 gpio.PinIO
var in_bit1 gpio.PinIO
var in_bit0 gpio.PinIO

func readNibble() byte {
	// let the Apple II know we are ready to read
	//fmt.Printf("let the Apple II know we are ready to read\n")
	out_read.Out(gpio.Low)

	// wait for the Apple II to write
	//fmt.Printf("wait for the Apple II to write\n")
	for in_write.Read() == gpio.High {
		in_write.WaitForEdge(-1)
	}

	// get a nibble of data
	//fmt.Printf("get a nibble of data\n")
	var nibble byte
	nibble = 0
	bit3 := in_bit3.Read()
	bit2 := in_bit2.Read()
	bit1 := in_bit1.Read()
	bit0 := in_bit0.Read()

	if bit3 == gpio.High {
		nibble += 8
	}
	if bit2 == gpio.High {
		nibble += 4
	}
	if bit1 == gpio.High {
		nibble += 2
	}
	if bit0 == gpio.High {
		nibble += 1
	}

	// let the Apple II know we are done reading
	//fmt.Printf("let the Apple II know we are done reading\n")
	out_read.Out(gpio.High)

	// wait for the Apple II to finish writing
	//fmt.Printf("wait for the Apple II to finish writing\n")
	for in_write.Read() == gpio.Low {
		in_write.WaitForEdge(-1)
	}

	return nibble
}

func readByte() byte {
	data := byte(0)
	data = readNibble() << byte(4)
	data += readNibble()
	return data
}

func writeNibble(data byte) {
	// wait for the Apple II to be ready to read
	//fmt.Printf("wait for the Apple II to be ready to read\n")
	for in_read.Read() == gpio.High {
		in_read.WaitForEdge(-1)
	}

	bit3 := gpio.Low
	bit2 := gpio.Low
	bit1 := gpio.Low
	bit0 := gpio.Low
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
	//fmt.Printf("let Apple II know we're writing\n")
	out_write.Out(gpio.Low)


	// wait for the Apple II to finsih reading
	//fmt.Printf("wait for the Apple II to finsih reading\n")
	for in_read.Read() == gpio.Low {
		in_read.WaitForEdge(-1)
	}

	// let the Apple II know we are done writing
	//fmt.Printf("let the Apple II know we are done writing\n")
	out_write.Out(gpio.High)
}

func writeByte(data byte) {
	writeNibble(data >> 4)
	writeNibble(data & 15)
}

func readBlock(buffer []byte) {
	for i := 0; i < 512; i++ {
		writeByte(buffer[i])
	}
}

func dumpBlock(buffer []byte) {
	for i := 0; i < 512; i++ {
		fmt.Printf("%02X ", buffer[i])
	}
}

func writeBlock(buffer []byte) {
	for i := 0; i < 512; i++ {
		buffer[i] = readByte()
	}
}

func main() {
	host.Init()

	out_write = gpioreg.ByName("GPIO5")
	out_read = gpioreg.ByName("GPIO11")
	out_commandWrite = gpioreg.ByName("GPIO9")
	out_commandRead = gpioreg.ByName("GPIO10")
	out_bit3 = gpioreg.ByName("GPIO22")
	out_bit2 = gpioreg.ByName("GPIO27")
	out_bit1 = gpioreg.ByName("GPIO17")
	out_bit0 = gpioreg.ByName("GPIO4")
	in_write = gpioreg.ByName("GPIO12")
	in_read = gpioreg.ByName("GPIO16")
	in_commandWrite = gpioreg.ByName("GPIO20")
	in_commandRead = gpioreg.ByName("GPIO21")
	in_bit3 = gpioreg.ByName("GPIO26")
	in_bit2 = gpioreg.ByName("GPIO19")
	in_bit1 = gpioreg.ByName("GPIO13")
	in_bit0 = gpioreg.ByName("GPIO6")

	fmt.Printf("Starting Apple II RPi...\n")
	fileName := "Total Replay v4.0-rc.1.hdv"
	file, err := os.OpenFile(fileName, os.O_RDWR, 0755)
	if err != nil {
	//log.Fatal(err)
	}
	//if err := f.Close(); err != nil {
	//log.Fatal(err)
	//}

	//for in_write.Read() == gpio.Low {
	//	in_write.WaitForEdge(-1)
	//}
	buffer := make([]byte, 512)
	//file.ReadAt(buffer, int64(0) * 512)
	//dumpBlock(buffer)

	for {
		command := readByte();
		if (command == 1) {
			blockLow := readByte();
			blockHigh := readByte();


			var block int64
			block = int64(blockHigh) * 256 + int64(blockLow)

			fmt.Printf("Read block %d\n", block)


			file.ReadAt(buffer, int64(block) * 512)
			//dumpBlock(buffer)
			readBlock(buffer)
		}
		if (command == 2) {
			blockLow := readByte();
			blockHigh := readByte();


			var block int64
			block = int64(blockHigh) * 256 + int64(blockLow)

			fmt.Printf("Write block %d\n", block)

			writeBlock(buffer)
			file.WriteAt(buffer, int64(block) * 512)
		}
	}
}

