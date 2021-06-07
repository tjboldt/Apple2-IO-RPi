package handlers

import (
	"fmt"
	"os"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

func WriteBlockCommand(drive1 *os.File, drive2 *os.File) {
	blockLow, _ := a2io.ReadByte()
	blockHigh, _ := a2io.ReadByte()

	var driveUnit byte = 0
	var err error

	if !oldFirmware {
		driveUnit, err = a2io.ReadByte()
		if err != nil {
			fmt.Printf("Drive unit not sent, assuming older firmware")
			oldFirmware = true
		}
	}

	file := drive1

	if driveUnit >= 128 {
		file = drive2
	}

	buffer := make([]byte, 512)
	var block int
	block = int(blockHigh)*256 + int(blockLow)

	fmt.Printf("Write block %d\n", block)

	a2io.ReadBlock(buffer)
	prodos.WriteBlock(file, block, buffer)
	fmt.Printf("Write block completed\n")
}
