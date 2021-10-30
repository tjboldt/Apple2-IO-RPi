package handlers

import (
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
