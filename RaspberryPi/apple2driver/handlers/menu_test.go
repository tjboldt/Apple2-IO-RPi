// Copyright Terence J. Boldt (c)2020-2022
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file contains the test for the handler for
// displaying the menu of choices on the Apple II

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
