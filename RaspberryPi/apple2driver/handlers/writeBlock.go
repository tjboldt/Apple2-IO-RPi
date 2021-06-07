package handlers

import (
	"fmt"
	"os"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

func WriteBlockCommand(file *os.File) {
	blockLow, _ := a2io.ReadByte()
	blockHigh, _ := a2io.ReadByte()

	buffer := make([]byte, 512)
	var block int
	block = int(blockHigh)*256 + int(blockLow)

	fmt.Printf("Write block %d\n", block)

	a2io.ReadBlock(buffer)
	prodos.WriteBlock(file, block, buffer)
	fmt.Printf("Write block completed\n")
}
