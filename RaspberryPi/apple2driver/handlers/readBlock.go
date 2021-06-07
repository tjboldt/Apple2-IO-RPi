package handlers

import (
	"fmt"
	"os"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

var oldFirmware = false

func ReadBlockCommand(drive1 *os.File, drive2 *os.File) {
	blockLow, _ := a2io.ReadByte()
	blockHigh, _ := a2io.ReadByte()
	var driveUnit byte = 0
	var err error

	if !oldFirmware {
		driveUnit, err = a2io.ReadByte()
		fmt.Printf("Drive unit: %0X\n", driveUnit)

		if err != nil {
			fmt.Printf("Drive unit not sent, assuming older firmware")
			oldFirmware = true
		}
	}

	file := drive1

	if driveUnit >= 128 {
		file = drive2
	}

	var block int
	block = int(blockHigh)*256 + int(blockLow)

	fmt.Printf("Read block %d\n", block)

	buffer := prodos.ReadBlock(file, block)

	err = a2io.WriteBlock(buffer)
	if err == nil {
		fmt.Printf("Read block completed\n")
	} else {
		fmt.Printf("Failed to read block\n")
	}
}
