package handlers

import (
	"fmt"
	"os"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
)

func WriteBlockCommand(file *os.File) {
	blockLow, _ := a2io.ReadByte()
	blockHigh, _ := a2io.ReadByte()

	buffer := make([]byte, 512)
	var block int64
	block = int64(blockHigh)*256 + int64(blockLow)

	fmt.Printf("Write block %d\n", block)

	a2io.ReadBlock(buffer)
	file.WriteAt(buffer, int64(block)*512)
	file.Sync()
	fmt.Printf("Write block completed\n")
}
