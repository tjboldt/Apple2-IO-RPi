// Copyright Terence J. Boldt (c)2020-2022
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file contains the main driver code for the Raspberry Pi side of
// the Apple2-IO-RPi hardware. Commands are sent from the Apple II and
// responses are sent back from the Raspberry Pi.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.bug.st/serial"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/handlers"
	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/info"
)

const resetCommand = 0
const readBlockCommand = 1
const writeBlockCommand = 2
const getTimeCommand = 3
const changeDriveCommand = 4 // not currently used
const execCommand = 5
const loadFileCommand = 6
const saveFileCommand = 7 // not implemented yet
const menuCommand = 8
const shellCommand = 9

func main() {
	drive1, drive2, cdc := getFlags()

	fmt.Printf("Starting Apple II RPi v%s...\n", info.Version)

	var comm a2io.A2Io
	if cdc {
		comm = a2io.CDCio{}
	} else {
		comm = a2io.A2Gpio{}
	}

	handlers.SetCommunication(comm)
	comm.Init()

	lastCommandTime := time.Now()

	// In case Apple II is waiting, send 0 byte to start
	comm.WriteByte(0)

	for {
		command, err := comm.ReadByte()
		if err == nil {
			lastCommandTime = time.Now()
			switch command {
			case resetCommand:
				handlers.ResetCommand()
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
			case shellCommand:
				handlers.ShellCommand()
			}
			// the A2Io interface should be extended in one way or another
			// to encapsulate this, e.g. by a ReadByte variant / parameter 
		} else if cdc {
			var portErr *serial.PortError
			if errors.As(err, &portErr) && portErr.Code() == serial.PortClosed {
				comm.Init()
			}
			// temporary workaround for busy wait loop heating up the RPi
		} else if time.Since(lastCommandTime) > time.Millisecond*100 {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func getFlags() (*os.File, *os.File, bool) {
	var drive1Name string
	var drive2Name string
	var cdc bool

	execName, _ := os.Executable()
	path := filepath.Dir(execName)
	path = filepath.Join(path, "..")
	path, _ = filepath.EvalSymlinks(path)
	defaultFileName := filepath.Join(path, "Apple2-IO-RPi.hdv")

	flag.StringVar(&drive1Name, "d1", "", "A ProDOS format drive image for drive 1")
	flag.StringVar(&drive2Name, "d2", defaultFileName, "A ProDOS format drive image for drive 2 and will be used for drive 1 if drive 1 empty")
	flag.BoolVar(&cdc, "cdc", false, "Use ACM CDC serial device")
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

	return drive1, drive2, cdc
}
