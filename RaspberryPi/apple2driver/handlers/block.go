// Copyright Terence J. Boldt (c)2020-2022
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is used for handling ProDOS image block reading and writing

package handlers

import (
	"fmt"
	"os"

	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

var oldFirmware = false

// ReadBlockCommand handles requests to read ProDOS blocks
func ReadBlockCommand(drive1 *os.File, drive2 *os.File) {
	blockLow, _ := comm.ReadByte()
	blockHigh, _ := comm.ReadByte()
	var driveUnit byte
	var err error

	if !oldFirmware {
		driveUnit, err = comm.ReadByte()
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

	block := int(blockHigh)*256 + int(blockLow)

	fmt.Printf("Read block %d\n", block)

	buffer := prodos.ReadBlock(file, block)

	err = comm.WriteBlock(buffer)
	if err == nil {
		fmt.Printf("Read block completed\n")
	} else {
		fmt.Printf("Failed to read block\n")
	}
}

// WriteBlockCommand handles requests to write ProDOS blocks
func WriteBlockCommand(drive1 *os.File, drive2 *os.File) {
	blockLow, _ := comm.ReadByte()
	blockHigh, _ := comm.ReadByte()

	var driveUnit byte
	var err error

	if !oldFirmware {
		driveUnit, err = comm.ReadByte()
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

	block := int(blockHigh)*256 + int(blockLow)

	fmt.Printf("Write block %d\n", block)

	comm.ReadBlock(buffer)
	prodos.WriteBlock(file, block, buffer)
	fmt.Printf("Write block completed\n")
}
