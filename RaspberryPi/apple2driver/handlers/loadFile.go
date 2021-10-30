package handlers

import (
	"fmt"
	"os"
)

func LoadFileCommand() {
	fileName, _ := comm.ReadString()

	file, err := os.OpenFile(fileName, os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		comm.WriteByte(0)
		comm.WriteByte(0)
		return
	}

	fileInfo, _ := file.Stat()
	fileSize := int(fileInfo.Size())

	fmt.Printf("FileSize: %d\n", fileSize)

	fileSizeHigh := byte(fileSize >> 8)
	fileSizeLow := byte(fileSize & 255)

	err = comm.WriteByte(fileSizeLow)
	if err != nil {
		return
	}
	err = comm.WriteByte(fileSizeHigh)
	if err != nil {
		return
	}

	buffer := make([]byte, fileSize)

	fmt.Printf("Read file %s SizeHigh: %d SizeLow: %d\n", fileName, fileSizeHigh, fileSizeLow)

	file.Read(buffer)

	comm.WriteBuffer(buffer)
}
