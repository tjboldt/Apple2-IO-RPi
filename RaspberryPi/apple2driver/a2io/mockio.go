// Copyright Terence J. Boldt (c)2021-2022
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is used for mocking the GPIO code so that unit tests can be run
// on any machine and not just a Raspberry Pi

package a2io

import (
	"errors"
	"strings"
)

// MockIoData is used to verify tests were successful
type MockIoData struct {
	BytesToRead        []byte
	BytesWritten       []byte
	NumberBytesRead    int
	NumberBytesWritten int
	ErrorToThrow       error
}

// MockIo implements A2Io to allow unit tests to run without needing GPIO functioning
type MockIo struct {
	Data *MockIoData
}

// Init is only here to complete A2Io interface
func (mockIo MockIo) Init() {

}

// WriteByte mocks writing a byte to the Apple II
func (mockIo MockIo) WriteByte(data byte) error {
	mockIo.Data.BytesWritten[mockIo.Data.NumberBytesWritten] = data
	mockIo.Data.NumberBytesWritten++
	return mockIo.Data.ErrorToThrow
}

// WriteString mocks writing a string to the Apple II
func (mockIo MockIo) WriteString(outString string) error {
	for i, b := range outString {
		mockIo.Data.BytesWritten[i+mockIo.Data.NumberBytesWritten] = byte(b)
	}
	mockIo.Data.NumberBytesWritten += len(outString)
	return mockIo.Data.ErrorToThrow
}

// WriteBlock mocks writing a block to the Apple II
func (mockIo MockIo) WriteBlock(buffer []byte) error {
	for i, b := range buffer {
		mockIo.Data.BytesWritten[i+mockIo.Data.NumberBytesWritten] = b
	}
	mockIo.Data.NumberBytesWritten += len(buffer)
	return mockIo.Data.ErrorToThrow
}

// WriteBuffer mocks writing a buffer to the Apple II
func (mockIo MockIo) WriteBuffer(buffer []byte) error {
	for i, b := range buffer {
		mockIo.Data.BytesWritten[i+mockIo.Data.NumberBytesWritten] = b
	}
	mockIo.Data.NumberBytesWritten += len(buffer)
	return mockIo.Data.ErrorToThrow
}

// ReadByte mocks reading a byte from the Apple II
func (mockIo MockIo) ReadByte() (byte, error) {
	b := mockIo.Data.BytesToRead[mockIo.Data.NumberBytesRead]
	mockIo.Data.NumberBytesRead++
	return b, mockIo.Data.ErrorToThrow
}

// ReadString mocks reading a null terminated string from the Apple II
func (mockIo MockIo) ReadString() (string, error) {
	builder := strings.Builder{}
	for {
		if mockIo.Data.NumberBytesRead > len(mockIo.Data.BytesToRead) {
			return "", errors.New("read more data than available")
		}
		builder.WriteByte(mockIo.Data.BytesToRead[mockIo.Data.NumberBytesRead])
		mockIo.Data.NumberBytesRead++
		if mockIo.Data.BytesToRead[mockIo.Data.NumberBytesRead] == 0 {
			mockIo.Data.NumberBytesRead++
			break
		}
	}
	return builder.String(), mockIo.Data.ErrorToThrow
}

// ReadBlock mocks reading a 512 byte block from the Apple II
func (mockIo MockIo) ReadBlock(buffer []byte) error {
	if mockIo.Data.NumberBytesRead+512 > len(mockIo.Data.BytesToRead) {
		return errors.New("read more data than available")
	}
	for i := 0; i < 512; i++ {
		buffer[i] = mockIo.Data.BytesToRead[mockIo.Data.NumberBytesRead]
		mockIo.Data.NumberBytesRead++
	}
	return mockIo.Data.ErrorToThrow
}

// SendCharacter is a pass-through to vt100 implementation
func (mockIo MockIo) SendCharacter(character byte) {
	sendCharacter(mockIo, character)
}

// ReadCharacter is a pass-through to vt100 implementation
func (mockIo MockIo) ReadCharacter() (string, error) {
	return readCharacter(mockIo)
}
