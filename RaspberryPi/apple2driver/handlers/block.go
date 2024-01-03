// Copyright Terence J. Boldt (c)2020-2024
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is used for handling ProDOS image block reading and writing

package handlers

import (
	"fmt"

	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

// ReadBlockCommand handles requests to read ProDOS blocks
func ReadBlockCommand(drive1 prodos.ReaderWriterAt, drive2 prodos.ReaderWriterAt) (int, error) {
	blockLow, _ := comm.ReadByte()
	blockHigh, _ := comm.ReadByte()
	var driveUnit byte
	var err error

	driveUnit, err = comm.ReadByte()

	if err != nil {
		fmt.Printf("Failed to read block")
		return 0, err
	}

	file := drive1
	driveNumber := 1

	if driveUnit >= 128 {
		file = drive2
		driveNumber = 2
	}

	slotNumber := driveUnit & 0x7F >> 4

	block := int(blockHigh)*256 + int(blockLow)

	fmt.Printf("Read block %04X in slot %d, drive %d...", block, slotNumber, driveNumber)

	buffer, err := prodos.ReadBlock(file, block)
	if err != nil {
		fmt.Printf("failed %s\n", err)
		return 0, err
	}

	err = comm.WriteBlock(buffer)
	if err == nil {
		fmt.Printf("succeeded\n")
	} else {
		fmt.Printf("failed %s\n", err)
	}

	return block, nil
}

// WriteBlockCommand handles requests to write ProDOS blocks
func WriteBlockCommand(drive1 prodos.ReaderWriterAt, drive2 prodos.ReaderWriterAt) error {
	blockLow, _ := comm.ReadByte()
	blockHigh, _ := comm.ReadByte()

	var driveUnit byte
	var err error

	driveUnit, err = comm.ReadByte()
	if err != nil {
		fmt.Printf("Failed to write block")
		return err
	}

	file := drive1
	driveNumber := 1

	if driveUnit >= 128 {
		file = drive2
		driveNumber = 2
	}

	buffer := make([]byte, 512)

	block := int(blockHigh)*256 + int(blockLow)

	slotNumber := driveUnit & 0x7F >> 4

	fmt.Printf("Write block %04X in slot %d, drive %d...", block, slotNumber, driveNumber)

	err = comm.ReadBlock(buffer)
	if err != nil {
		fmt.Printf("failed %s\n", err)
		return err
	}
	fmt.Printf("...")

	err = prodos.WriteBlock(file, block, buffer)
	if err == nil {
		fmt.Printf("succeeded\n")
	} else {
		fmt.Printf("failed\n")
	}

	return nil
}
