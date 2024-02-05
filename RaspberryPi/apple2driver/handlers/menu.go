// Copyright Terence J. Boldt (c)2020-2024
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file contains the handler for displaying the menu of choices on
// the Apple II

package handlers

import (
	"fmt"
)

// MenuCommand handles the request to show menu options on the Apple II
func MenuCommand() {
	fmt.Printf("Sending menu...\n")
	comm.WriteString("Apple2-IO-RPi\r" +
		"(c)2020-2024 Terence J. Boldt\r" +
		"\r" +
		"Select an option:\r" +
		"\r" +
		"1. Boot\r" +
		"2. Command Line\r" +
		"3. Load File\r")
}
