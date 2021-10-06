package handlers

import (
	"fmt"
	"os"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
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

	buffer := make([]byte, 512)
	var block int64
	block = int64(blockHigh)*256 + int64(blockLow)

	fmt.Printf("Read block %d\n", block)

	file.ReadAt(buffer, int64(block)*512)

	err = a2io.WriteBlock(buffer)
	if err == nil {
		fmt.Printf("Read block completed\n")
	} else {
		fmt.Printf("Failed to read block\n")
	}
}
