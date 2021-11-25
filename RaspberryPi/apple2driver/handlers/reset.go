// Copyright Terence J. Boldt (c)2020-2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is contains the handler for displaying the menu of choices on
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
