package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	drive1, drive2 := getDriveFiles()

	fmt.Printf("Starting Apple II RPi...\n")

	comm := a2io.A2Gpio{}

	handlers.SetCommunication(comm)
	comm.Init()

	for {
		command, err := comm.ReadByte()
		if err == nil {
			switch command {
			case ReadBlockCommand:
				handlers.ReadBlockCommand(drive1, drive2)
			case WriteBlockCommand:
				handlers.WriteBlockCommand(drive1, drive2)
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

func getDriveFiles() (*os.File, *os.File) {
	var drive1Name string
	var drive2Name string

	execName, _ := os.Executable()
	path := filepath.Dir(execName)
	path = filepath.Join(path, "..")
	path, _ = filepath.EvalSymlinks(path)
	defaultFileName := filepath.Join(path, "Apple2-IO-RPi.hdv")

	flag.StringVar(&drive1Name, "d1", "", "A ProDOS format drive image for drive 1")
	flag.StringVar(&drive2Name, "d2", defaultFileName, "A ProDOS format drive image for drive 2 and will be used for drive 1 if drive 1 empty")
	flag.Parse()

	var drive1 *os.File = nil
	var drive2 *os.File = nil
	var err error

	if len(drive1Name) > 0 {
		fmt.Printf("Opening Drive 1 as: %s\n", drive1Name)
		drive1, err = os.OpenFile(drive1Name, os.O_RDWR, 0755)
		if err != nil {
			fmt.Printf("ERROR: %s", err.Error())
			os.Exit(1)
		}
	}

	if len(drive2Name) > 0 {
		if drive1 == nil {
			fmt.Printf("Opening Drive 1 as: %s because Drive 1 was empty\n", drive2Name)
			drive1, err = os.OpenFile(drive2Name, os.O_RDWR, 0755)
			if err != nil {
				fmt.Printf("ERROR: %s", err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Printf("Opening Drive 2 as: %s\n", drive2Name)
			drive2, err = os.OpenFile(drive2Name, os.O_RDWR, 0755)
			if err != nil {
				fmt.Printf("ERROR: %s", err.Error())
				os.Exit(1)
			}
		}
	}

	if drive1 == nil {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return drive1, drive2
}
