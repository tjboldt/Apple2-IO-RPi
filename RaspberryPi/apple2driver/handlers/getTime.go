// Copyright Terence J. Boldt (c)2020-2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is contains the handler for retrieving the ProDOS timestamp
// based on the current time

package handlers

import (
	"fmt"
	"time"

	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

// GetTimeCommand handles the request to get ProDOS time bytes
func GetTimeCommand() {
	fmt.Printf("Sending date/time...\n")
	prodosTime := prodos.DateTimeToProDOS(time.Now())

	for i := 0; i < len(prodosTime); i++ {
		comm.WriteByte(prodosTime[i])
	}

	fmt.Printf("Send time complete\n")
}
