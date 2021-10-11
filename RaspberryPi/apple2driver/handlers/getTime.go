package handlers

import (
	"fmt"
	"time"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

func GetTimeCommand() {
	fmt.Printf("Sending date/time...\n")
	prodosTime := prodos.DateTimeToProDOS(time.Now())

	for i := 0; i < len(prodosTime); i++ {
		a2io.WriteByte(prodosTime[i])
	}

	fmt.Printf("Send time complete\n")
}
