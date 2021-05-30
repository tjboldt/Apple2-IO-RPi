package handlers

import (
	"fmt"
	"os"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
)

func LoadFileCommand() {
	fileName, _ := a2io.ReadString()

	file, err := os.OpenFile(fileName, os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		a2io.WriteByte(0)
		a2io.WriteByte(0)
		return
	}

	fileInfo, _ := file.Stat()
	fileSize := int(fileInfo.Size())

	fmt.Printf("FileSize: %d\n", fileSize)

	fileSizeHigh := byte(fileSize >> 8)
	fileSizeLow := byte(fileSize & 255)

	a2io.WriteByte(fileSizeLow)
	a2io.WriteByte(fileSizeHigh)

	buffer := make([]byte, fileSize)

	fmt.Printf("Read file %s SizeHigh: %d SizeLow: %d\n", fileName, fileSizeHigh, fileSizeLow)

	file.Read(buffer)

	a2io.WriteBuffer(buffer)
}
