// Copyright Terence J. Boldt (c)2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file defines the interface for communicating with the Apple II via
// the Raspberry Pi GPIO ports but can also be mocked out for tests or
// passed input to the user for interactive testing

package a2io

// A2Io provides an interface to send and receive data to the Apple II
type A2Io interface {
	Init()
	WriteByte(data byte) error
	WriteString(outString string) error
	WriteBlock(buffer []byte) error
	WriteBuffer(buffer []byte) error
	ReadByte() (byte, error)
	ReadString() (string, error)
	ReadBlock(buffer []byte) error
}
