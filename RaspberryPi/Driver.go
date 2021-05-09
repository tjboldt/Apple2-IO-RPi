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
	"strings"
	"time"
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

const ReadBlockCommand = 1
const WriteBlockCommand = 2
const GetTimeCommand = 3
const ChangeDriveCommand = 4
const ExecCommand = 5
const LoadFileCommand = 6
const SaveFileCommand = 7

var debug bool = false

var workingDirectory string = "/home"

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
			case LoadFileCommand:
				handleLoadFileCommand()
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
	err := readBlock(buffer)
	if err == nil {
		fmt.Printf("Read block completed\n")
	} else {
		fmt.Printf("Failed to read block\n")
	}
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
	fmt.Printf("Write block completed\n")
}

func handleExecCommand() {
	fmt.Printf("Reading command to execute...\n")
	linuxCommand, err := readString()
	fmt.Printf("Command to run: %s\n", linuxCommand)
	if strings.HasPrefix(linuxCommand, "cd /") {
		workingDirectory = strings.Replace(linuxCommand, "cd ", "", 1)
		writeString("Working directory set")
		return
	}
	if strings.HasPrefix(linuxCommand, "cd ") {
		workingDirectory = workingDirectory + "/" + strings.Replace(linuxCommand, "cd ", "", 1)
		writeString("Working directory set")
		return
	}
        if linuxCommand == "a2help" {
		writeString("This is a pseudo shell. Each command is executed as a process. The cd command is\n" +
                            "intercepted and sets the working directory for the next command. Running\n" +
                            "commands that do not exit will hang. For example, do not use ping without a -c 1\n" +
                            "or something else to limit its output.\n")
	}
	cmd := exec.Command("bash", "-c", linuxCommand)
	cmd.Dir = workingDirectory
	cmdOut, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to execute command\n")
		writeString("Failed to execute command")
		return
	}
	fmt.Printf("Command output: %s\n", cmdOut)
	apple2string := strings.Replace(string(cmdOut), "\n", "\r", -1)
	err = writeString(apple2string)
	if err != nil {
		fmt.Printf("Failed to send command output\n")
		return
	}
}

func handleGetTimeCommand() {
	fmt.Printf("Sending date/time...\n")
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
	fmt.Printf("Send time complete\n")
}

func handleLoadFileCommand() {
	fileName, _ := readString()

	file, err := os.OpenFile(fileName, os.O_RDWR, 0755)
        if err != nil {
                fmt.Printf("ERROR: %s\n", err.Error())
                writeByte(0)
		writeByte(0)
		return 
        }

	fileInfo, _ := file.Stat()
	fileSize := int(fileInfo.Size())

	fmt.Printf("FileSize: %d\n", fileSize)

	fileSizeHigh := byte(fileSize >> 8)
	fileSizeLow := byte(fileSize & 255)

	writeByte(fileSizeLow)
	writeByte(fileSizeHigh)

        buffer := make([]byte, fileSize)

        fmt.Printf("Read file %s SizeHigh: %d SizeLow: %d\n", fileName, fileSizeHigh, fileSizeLow)

        file.Read(buffer)

	for i := 0; i < fileSize; i++ {
		err := writeByte(buffer[i])
		if err != nil {
			return
		}
	}
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

	edgeTimeout = time.Second * 5
}

func dumpBlock(buffer []byte) {
	for i := 0; i < 512; i++ {
		fmt.Printf("%02X ", buffer[i])
	}
}

func readString() (string, error) {
	var inBytes bytes.Buffer
	for {
		inByte, err := readByte()
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

func writeString(outString string) error {
	for _, character := range outString {
		err := writeByte(byte(character) | 128)
		if err != nil {
			fmt.Printf("Failed to write string\n")
			return err
		}
	}
	writeByte(0)
	return nil
}

func readByte() (byte, error) {
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
				fmt.Printf("Timed out reading byte -- write stuck high\n")
			}
			return 0, errors.New("Timed out reading byte -- write stuck high\n")
		}
	}

	// get a nibble of data
	if debug {
		fmt.Printf("get a byte of data\n")
	}
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
			if debug {
				fmt.Printf("Timed out reading byte -- write stuck low\n")
			}
			return 0, errors.New("Timed out reading byte -- write stuck low")
		}
	}

	return data, nil
}

func writeByte(data byte) error {
	// wait for the Apple II to be ready to read
	if debug {
		fmt.Printf("wait for the Apple II to be ready to read\n")
	}
	for in_read.Read() == gpio.High {
		if !in_read.WaitForEdge(edgeTimeout) {
			if debug {
				fmt.Printf("Timed out writing byte -- read stuck high\n")
			}
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
	if debug {
		fmt.Printf("let Apple II know we're writing\n")
	}
	out_write.Out(gpio.Low)

	// wait for the Apple II to finsih reading
	//fmt.Printf("wait for the Apple II to finsih reading\n")
	for in_read.Read() == gpio.Low {
		if !in_read.WaitForEdge(edgeTimeout) {
			if debug {
				fmt.Printf("Timed out writing byte -- read stuck low\n")
			}
			return errors.New("Timed out writing byte -- read stuck low")
		}
	}

	// let the Apple II know we are done writing
	if debug {
		fmt.Printf("let the Apple II know we are done writing\n")
	}
	out_write.Out(gpio.High)
	return nil
}
