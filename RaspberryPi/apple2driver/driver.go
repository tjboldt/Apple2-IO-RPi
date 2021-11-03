// Copyright Terence J. Boldt (c)2020-2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is contains the main driver code for the Raspberry Pi side of
// the Apple2-IO-RPi hardware. Commands are sent from the Apple II and
// responses are sent back from the Raspberry Pi.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/handlers"
)

const readBlockCommand = 1
const writeBlockCommand = 2
const getTimeCommand = 3
const changeDriveCommand = 4
const execCommand = 5
const loadFileCommand = 6
const saveFileCommand = 7
const menuCommand = 8

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
			case readBlockCommand:
				handlers.ReadBlockCommand(drive1, drive2)
			case writeBlockCommand:
				handlers.WriteBlockCommand(drive1, drive2)
			case getTimeCommand:
				handlers.GetTimeCommand()
			case execCommand:
				handlers.ExecCommand()
			case loadFileCommand:
				handlers.LoadFileCommand()
			case menuCommand:
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

	var drive1 *os.File
	var drive2 *os.File
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
