// Copyright Terence J. Boldt (c)2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is used for mocking the GPIO code so that unit tests can be run
// on any machine and not just a Raspberry Pi
package a2io

import (
	"errors"
	"strings"
)

type MockIoData struct {
	BytesToRead  []byte
	BytesWritten []byte
	byteRead     int
	byteWritten  int
	ErrorToThrow error
}

type MockIo struct {
	Data *MockIoData
}

func (mockIo MockIo) Init() {

}

func (mockIo MockIo) WriteByte(data byte) error {
	mockIo.Data.BytesWritten[mockIo.Data.byteWritten] = data
	mockIo.Data.byteWritten++
	return mockIo.Data.ErrorToThrow
}

func (mockIo MockIo) WriteString(outString string) error {
	for i, b := range outString {
		mockIo.Data.BytesWritten[i+mockIo.Data.byteWritten] = byte(b)
	}
	mockIo.Data.byteWritten += len(outString)
	return mockIo.Data.ErrorToThrow
}

func (mockIo MockIo) WriteBlock(buffer []byte) error {
	for i, b := range buffer {
		mockIo.Data.BytesWritten[i+mockIo.Data.byteWritten] = b
	}
	mockIo.Data.byteWritten += len(buffer)
	return mockIo.Data.ErrorToThrow
}

func (mockIo MockIo) WriteBuffer(buffer []byte) error {
	for i, b := range buffer {
		mockIo.Data.BytesWritten[i+mockIo.Data.byteWritten] = b
	}
	mockIo.Data.byteWritten += len(buffer)
	return mockIo.Data.ErrorToThrow
}

func (mockIo MockIo) ReadByte() (byte, error) {
	b := mockIo.Data.BytesToRead[mockIo.Data.byteRead]
	mockIo.Data.byteRead++
	return b, mockIo.Data.ErrorToThrow
}

func (mockIo MockIo) ReadString() (string, error) {
	builder := strings.Builder{}
	for {
		if mockIo.Data.byteRead > len(mockIo.Data.BytesToRead) {
			return "", errors.New("Read more data than available")
		}
		builder.WriteByte(mockIo.Data.BytesToRead[mockIo.Data.byteRead])
		mockIo.Data.byteRead++
		if mockIo.Data.BytesToRead[mockIo.Data.byteRead] == 0 {
			mockIo.Data.byteRead++
			break
		}
	}
	return builder.String(), mockIo.Data.ErrorToThrow
}

func (mockIo MockIo) ReadBlock(buffer []byte) error {
	if mockIo.Data.byteRead+512 > len(mockIo.Data.BytesToRead) {
		return errors.New("Read more data than available")
	}
	for i := 0; i < 512; i++ {
		buffer[i] = mockIo.Data.BytesToRead[mockIo.Data.byteRead]
		mockIo.Data.byteRead++
	}
	return mockIo.Data.ErrorToThrow
}
