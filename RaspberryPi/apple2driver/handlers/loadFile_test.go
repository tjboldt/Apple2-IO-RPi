// Copyright Terence J. Boldt (c)2020-2023
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file contains the tests for the handler for
// loading files directly from the Raspberry Pi onto the Apple II

package handlers

import (
	"os"
	"testing"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
)

func TestLoadFileCommandReturnZeroLengthOnFileNotFound(t *testing.T) {
	mockIoData := a2io.MockIoData{
		BytesToRead:  []byte("FILE_DOES_NOT_EXIST\x00"),
		BytesWritten: make([]byte, 1000),
	}

	mockIo := a2io.MockIo{Data: &mockIoData}

	SetCommunication(&mockIo)

	LoadFileCommand()

	got := mockIoData.BytesWritten

	if got[0] != 0 && got[1] != 0 {
		t.Errorf("MenuCommand() sent = %s; want 00", got)
	}
}

func TestLoadFileCommandReturnsFile(t *testing.T) {
	mockIoData := a2io.MockIoData{
		BytesToRead:  []byte("loadFile.go\x00"),
		BytesWritten: make([]byte, 10000),
	}

	mockIo := a2io.MockIo{Data: &mockIoData}

	SetCommunication(&mockIo)

	LoadFileCommand()

	fileInfo, _ := os.Stat("loadFile.go")

	got := mockIoData.BytesWritten

	lowSize := fileInfo.Size() & 0xFF
	hiSize := (fileInfo.Size() & 0xFF00) >> 8

	if got[0] != byte(lowSize) && got[1] != byte(hiSize) {
		t.Errorf("MenuCommand() fileSize = %02X%02X; want %02X%02X", got[1], got[0], hiSize, lowSize)
	}

	if mockIoData.NumberBytesWritten-2 != int(fileInfo.Size()) {
		t.Errorf("MenuCommand() byte count = %04X; want %04X", mockIoData.NumberBytesWritten-2, fileInfo.Size())
	}
}
