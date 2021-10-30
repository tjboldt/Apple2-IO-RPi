package handlers

import (
	"fmt"
	"time"

	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

func GetTimeCommand() {
	fmt.Printf("Sending date/time...\n")
	prodosTime := prodos.DateTimeToProDOS(time.Now())

	for i := 0; i < len(prodosTime); i++ {
		comm.WriteByte(prodosTime[i])
	}

	fmt.Printf("Send time complete\n")
}
