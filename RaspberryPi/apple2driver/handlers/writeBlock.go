package handlers

import (
	"fmt"
	"os"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
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

	if driveUnit >= 8 {
		file = drive2
	}

	buffer := make([]byte, 512)
	var block int64
	block = int64(blockHigh)*256 + int64(blockLow)

	fmt.Printf("Write block %d\n", block)

	a2io.ReadBlock(buffer)
	file.WriteAt(buffer, int64(block)*512)
	file.Sync()
	fmt.Printf("Write block completed\n")
}
