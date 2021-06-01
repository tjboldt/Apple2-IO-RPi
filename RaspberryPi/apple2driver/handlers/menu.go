package handlers

import (
	"fmt"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
)

func MenuCommand() {
	fmt.Printf("Sending menu...\n")
	a2io.WriteString("Apple2-IO-RPi\r" +
		"(c)2020-2021 Terence J. Boldt\r" +
		"\r" +
		"Select an option:\r" +
		"\r" +
		"1. Boot\r" +
		"2. Command Line\r" +
		"3. Load File\r")
}
