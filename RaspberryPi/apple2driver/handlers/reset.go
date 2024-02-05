// Copyright Terence J. Boldt (c)2020-2024
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file contains the handler for displaying the menu of choices on
// the Apple II

package handlers

import (
	"fmt"
)

// ResetCommand responds with zero
func ResetCommand() {
	fmt.Printf("Reset request\n")
	comm.WriteByte(0)
}
