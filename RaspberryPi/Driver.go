package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	"time"
)

var edgeTimeout time.Duration

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

const ReadBlockCommand = 1
const WriteBlockCommand = 2
const GetTimeCommand = 3
const ChangeDriveCommand = 4
const ExecCommand = 5

var debug bool = false

func main() {
	host.Init()

	initGpio()

	if len(os.Args) == 3 && os.Args[2] == "--debug" {
		debug = true
	}

	fmt.Printf("Starting Apple II RPi...\n")

	fileName := os.Args[1]

	file, err := os.OpenFile(fileName, os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		os.Exit(1)
	}


	for {
		if debug {
			fmt.Printf("Check for command\n")
		}

		command, err := readByte()
		if err != nil {
			//fmt.Printf("Timed out waiting for command\n")
		} else {
			switch command {
			case ReadBlockCommand:
				handleReadBlockCommand(file)
			case WriteBlockCommand:
				handleWriteBlockCommand(file)
			case GetTimeCommand:
				handleGetTimeCommand()
			case ExecCommand:
				handleExecCommand()
			}
		}
	}
}

func handleReadBlockCommand(file *os.File) {
	blockLow, _ := readByte()
	blockHigh, _ := readByte()

	buffer := make([]byte, 512)
	var block int64
	block = int64(blockHigh)*256 + int64(blockLow)

	fmt.Printf("Read block %d\n", block)

	file.ReadAt(buffer, int64(block)*512)
	//dumpBlock(buffer)
	readBlock(buffer)
}

func handleWriteBlockCommand(file *os.File) {
	blockLow, _ := readByte()
	blockHigh, _ := readByte()

	buffer := make([]byte, 512)
	var block int64
	block = int64(blockHigh)*256 + int64(blockLow)

	fmt.Printf("Write block %d\n", block)

	writeBlock(buffer)
	file.WriteAt(buffer, int64(block)*512)
	file.Sync()
}

func handleExecCommand() {
	linuxCommand,err := readString()
	cmd := exec.Command("bash", "-c", linuxCommand)
	cmdOut, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to execute command\n")
	}
	writeString(cmdOut)
}

func handleGetTimeCommand() {
	/*  49041 ($BF91)     49040 ($BF90)

	        7 6 5 4 3 2 1 0   7 6 5 4 3 2 1 0
	       +-+-+-+-+-+-+-+-+ +-+-+-+-+-+-+-+-+
	DATE:  |    year     |  month  |   day   |
	       +-+-+-+-+-+-+-+-+ +-+-+-+-+-+-+-+-+

	        7 6 5 4 3 2 1 0   7 6 5 4 3 2 1 0
	       +-+-+-+-+-+-+-+-+ +-+-+-+-+-+-+-+-+
	TIME:  |    hour       | |    minute     |
	       +-+-+-+-+-+-+-+-+ +-+-+-+-+-+-+-+-+

	       49043 ($BF93)     49042 ($BF92)
	*/
	now := time.Now()

	year := now.Year() % 100
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()

	bf91 := (byte(year) << 1) + (byte(month) >> 3)
	bf90 := ((byte(month) & 15) << 5) + byte(day)
	bf93 := byte(hour)
	bf92 := byte(minute)

	writeByte(bf90)
	writeByte(bf91)
	writeByte(bf92)
	writeByte(bf93)
}

func readBlock(buffer []byte) error {
	for i := 0; i < 512; i++ {
		err := writeByte(buffer[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func writeBlock(buffer []byte) error {
	var err error
	for i := 0; i < 512; i++ {
		buffer[i], err = readByte()
		if err != nil {
			return err
		}
	}

	return nil
}

func initGpio() {
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

	in_write.In(gpio.PullDown, gpio.BothEdges)
	in_read.In(gpio.PullDown, gpio.BothEdges)

	edgeTimeout = time.Second * 5
}

func dumpBlock(buffer []byte) {
	for i := 0; i < 512; i++ {
		fmt.Printf("%02X ", buffer[i])
	}
}

func readString() (string, error) {
	inByte := byte(0)
	var inBytes bytes.Buffer
	var err error
	for inByte == 0 {
		inByte,err = readByte()
		if err != nil {
			return "", err
		}
		inBytes.WriteByte(inByte)
	}
	return string(inBytes.Bytes()), nil
}

func writeString(outBytes []byte) error {
	for outByte := range outBytes {
		err := writeByte(byte(outByte))
		if err != nil {
			return err
		}
	}
	writeByte(0)
	return nil
}

func readByte() (byte, error) {
	data, err := readNibble()
	data = data << byte(4)
	if err != nil {
		return 0, err
	}
	highNibble, err := readNibble()
	if err != nil {
		return 0, err
	}
	data += highNibble
	//fmt.Printf("R%02X ", data)
	return data, nil
}

func writeByte(data byte) error {
	//fmt.Printf("W%02X ", data)
	err := writeNibble(data >> 4)
	if err != nil {
		return err
	}
	err = writeNibble(data & 15)
	if err != nil {
		return err
	}

	return nil
}

func readNibble() (byte, error) {
	// let the Apple II know we are ready to read
	if debug {
		fmt.Printf("let the Apple II know we are ready to read\n")
	}
	out_read.Out(gpio.Low)

	// wait for the Apple II to write
	if debug {
		fmt.Printf("wait for the Apple II to write\n")
	}
	for in_write.Read() == gpio.High {
		if !in_write.WaitForEdge(edgeTimeout) {
			if debug {
				fmt.Printf("Timed out reading nibble -- write stuck high\n")
			}
			return 0, errors.New("Timed out reading nibble -- write stuck high\n")
		}
	}

	// get a nibble of data
	if debug {
		fmt.Printf("get a nibble of data\n")
	}
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
		if !in_write.WaitForEdge(edgeTimeout) {
			if debug {
				fmt.Printf("Timed out reading nibble -- write stuck low\n")
			}
			return 0, errors.New("Timed out reading nibble -- write stuck low")
		}
	}

	return nibble, nil
}

func writeNibble(data byte) error {
	// wait for the Apple II to be ready to read
	if debug {
		fmt.Printf("wait for the Apple II to be ready to read\n")
	}
	for in_read.Read() == gpio.High {
		if !in_read.WaitForEdge(edgeTimeout) {
			if debug {
				fmt.Printf("Timed out writing nibble -- read stuck high\n")
			}
			return errors.New("Timed out writing nibble -- read stuck high")
		}
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
	if debug {
		fmt.Printf("let Apple II know we're writing\n")
	}
	out_write.Out(gpio.Low)

	// wait for the Apple II to finsih reading
	//fmt.Printf("wait for the Apple II to finsih reading\n")
	for in_read.Read() == gpio.Low {
		if !in_read.WaitForEdge(edgeTimeout) {
			if debug {
				fmt.Printf("Timed out writing nibble -- read stuck low\n")
			}
			return errors.New("Timed out writing nibble -- read stuck low")
		}
	}

	// let the Apple II know we are done writing
	if debug {
		fmt.Printf("let the Apple II know we are done writing\n")
	}
	out_write.Out(gpio.High)
	return nil
}
