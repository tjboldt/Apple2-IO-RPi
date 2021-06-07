package handlers

import (
	"fmt"
	"os"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

func ReadBlockCommand(file *os.File) {
	blockLow, _ := a2io.ReadByte()
	blockHigh, _ := a2io.ReadByte()

	var block int
	block = int(blockHigh)*256 + int(blockLow)

	fmt.Printf("Read block %d\n", block)

	buffer := prodos.ReadBlock(file, block)

	err := a2io.WriteBlock(buffer)
	if err == nil {
		fmt.Printf("Read block completed\n")
	} else {
		fmt.Printf("Failed to read block\n")
	}
}
