package handlers

import (
	"fmt"
	"time"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
)

func GetTimeCommand() {
	fmt.Printf("Sending date/time...\n")
	/*  49041 ($BF91)     49040 ($BF90)

	        7 6 5 4 3 2 1 0   7 6 5 4 3 2 1 0
	       +-+-+-+-+-+-+-+-+ +-+-+-+-+-+-+-+-+
	DATE:  |    year     |  month  |   day   |
	       +-+-+-+-+-+-+-+-+ +-+-+-+-+-+-+-+-+

	        7 6 5 4 3 2 1 0   7 6 5 4 3 2 1 0
	       +-+-+-+-+-+-+-+-+ +-+-+-+-+-+-+-+-+
	TIME:  |    hour       | |    minute     |
	       +-+-+-+-+-+-+-+-+ +-+-+-+-+-+-+-+-+

	       49043 ($BF93)     49042 ($BF92)
	*/
	now := time.Now()

	year := now.Year() % 100
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()

	bf91 := (byte(year) << 1) + (byte(month) >> 3)
	bf90 := ((byte(month) & 15) << 5) + byte(day)
	bf93 := byte(hour)
	bf92 := byte(minute)

	a2io.WriteByte(bf90)
	a2io.WriteByte(bf91)
	a2io.WriteByte(bf92)
	a2io.WriteByte(bf93)
	fmt.Printf("Send time complete\n")
}
