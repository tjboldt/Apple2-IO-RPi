// Copyright Terence J. Boldt (c)2020-2022
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file contains the handler for retrieving the ProDOS timestamp
// based on the current time

package handlers

import (
	"fmt"
	"time"

	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

// GetTimeCommand handles the request to get ProDOS time bytes
func GetTimeCommand() {
	fmt.Printf("Sending date/time...")
	prodosTime := prodos.DateTimeToProDOS(time.Now())

	for i := 0; i < len(prodosTime); i++ {
		comm.WriteByte(prodosTime[i])
	}

	fmt.Printf("%02X %02X %02X %02X\n", prodosTime[0], prodosTime[1], prodosTime[2], prodosTime[3])
}
