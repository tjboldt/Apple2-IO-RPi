// Copyright Terence J. Boldt (c)2020-2022
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file contains the main driver code for the Raspberry Pi side of
// the Apple2-IO-RPi hardware. Commands are sent from the Apple II and
// responses are sent back from the Raspberry Pi.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/handlers"
	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/info"
	"github.com/tjboldt/ProDOS-Utilities/prodos"
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
	drive1, drive2 := getDriveFiles()

	fmt.Printf("Starting Apple II RPi v%s...\n", info.Version)

	comm := a2io.A2Gpio{}

	handlers.SetCommunication(comm)
	comm.Init()

	lastCommandTime := time.Now()

	// In case Apple II is waiting, send 0 byte to start
	comm.WriteByte(0)

	cwd, _ := os.Getwd()

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
				newCwd, _ := os.Getwd()
				if newCwd != cwd {
					cwd = newCwd
					drive1, err = generateDrive1FromCwd()
				}
			case loadFileCommand:
				handlers.LoadFileCommand()
			case menuCommand:
				handlers.MenuCommand()
			case shellCommand:
				handlers.ShellCommand()
			}
			// temporary workaround for busy wait loop heating up the RPi
		} else if time.Since(lastCommandTime) > time.Millisecond*100 {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func getDriveFiles() (prodos.ReaderWriterAt, prodos.ReaderWriterAt) {
	var drive1Name string
	var drive2Name string

	execName, _ := os.Executable()
	path := filepath.Dir(execName)
	path = filepath.Join(path, "..")
	path, _ = filepath.EvalSymlinks(path)

	flag.StringVar(&drive1Name, "d1", "", "A ProDOS format drive image for drive 1")
	flag.StringVar(&drive2Name, "d2", "", "A ProDOS format drive image for drive 2 and will be used for drive 1 if drive 1 empty")
	flag.Parse()

	var drive1 prodos.ReaderWriterAt
	var drive2 prodos.ReaderWriterAt
	var err error

	if len(drive1Name) > 0 {
		fmt.Printf("Opening Drive 1 as: %s\n", drive1Name)
		drive1, err = os.OpenFile(drive1Name, os.O_RDWR, 0755)
		logAndExitOnErr(err)
	} else {
		var exec string
		exec, err = os.Executable()
		if err != nil {
			fmt.Printf("ERROR: %s", err.Error())
			os.Exit(1)
		}
		cwd := filepath.Dir(exec)
		err = os.Chdir(filepath.Join(cwd, "../driveimage"))
		logAndExitOnErr(err)
		drive1, err = generateDrive1FromCwd()
		logAndExitOnErr(err)
	}

	if len(drive2Name) > 0 {
		fmt.Printf("Opening Drive 2 as: %s\n", drive2Name)
		drive2, err = os.OpenFile(drive2Name, os.O_RDWR, 0755)
		logAndExitOnErr(err)
	}

	return drive1, drive2
}

func logAndExitOnErr(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		os.Exit(1)
	}
}

func generateDrive1FromCwd() (prodos.ReaderWriterAt, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	drive1 := prodos.NewMemoryFile(0x2000000)
	fmt.Printf("Generating Drive 1 in memory from: %s\n", cwd)
	prodos.CreateVolume(drive1, "APPLE2.IO.RPI", 65535)
	err = prodos.AddFilesFromHostDirectory(drive1, cwd)
	return drive1, err
}
