package handlers

import (
	"strings"
	"testing"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
)

func TestMenuCommand(t *testing.T) {
	mockIoData := a2io.MockIoData{
		ErrorToThrow: nil,
		BytesWritten: make([]byte, 1000),
	}

	mockIo := a2io.MockIo{Data: &mockIoData}

	SetCommunication(&mockIo)

	MenuCommand()

	want := "Apple2-IO-RPi"
	got := string(mockIoData.BytesWritten)

	if strings.Index(got, want) != 0 {
		t.Errorf("MenuCommand() sent = %s; want startsWith %s", got, want)
	}
}
