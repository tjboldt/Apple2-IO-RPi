package main

import (
	"fmt"
	"os"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/handlers"
)

const ReadBlockCommand = 1
const WriteBlockCommand = 2
const GetTimeCommand = 3
const ChangeDriveCommand = 4
const ExecCommand = 5
const LoadFileCommand = 6
const SaveFileCommand = 7
const MenuCommand = 8

func main() {
	a2io.InitGpio()

	fmt.Printf("Starting Apple II RPi...\n")

	fileName := os.Args[1]

	file, err := os.OpenFile(fileName, os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		os.Exit(1)
	}

	for {
		command, err := a2io.ReadByte()
		if err == nil {
			switch command {
			case ReadBlockCommand:
				handlers.ReadBlockCommand(file)
			case WriteBlockCommand:
				handlers.WriteBlockCommand(file)
			case GetTimeCommand:
				handlers.GetTimeCommand()
			case ExecCommand:
				handlers.ExecCommand()
			case LoadFileCommand:
				handlers.LoadFileCommand()
			case MenuCommand:
				handlers.MenuCommand()
			}
		}
	}
}
