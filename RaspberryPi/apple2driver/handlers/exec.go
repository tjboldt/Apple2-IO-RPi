// Copyright Terence J. Boldt (c)2020-2023
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file contains the handler for executing Linux and internal
// commands

package handlers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/drive"
	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/info"
	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

var forceLowercase = false
var execTimeoutSeconds = int(10)

// ExecCommand handles requests for the Apple II executing Linux commands
func ExecCommand(drive1 *prodos.ReaderWriterAt, drive2 *prodos.ReaderWriterAt) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		workingDirectory = "/home"
		comm.WriteString("Failed to get current working directory, setting to /home\r")
	}

	fmt.Printf("Reading command to execute...\n")
	linuxCommand, err := comm.ReadString()
	if forceLowercase {
		linuxCommand = strings.ToLower(linuxCommand)
	}
	linuxCommand = strings.Trim(linuxCommand, " ")
	if linuxCommand == "" {
		linuxCommand = "a2help"
	}
	fmt.Printf("Command to run: %s\n", linuxCommand)
	if strings.HasPrefix(linuxCommand, "cd ") {
		workingDirectory = strings.Replace(linuxCommand, "cd ", "", 1)
		err = os.Chdir(workingDirectory)
		if err != nil {
			comm.WriteString("Failed to set working directory\r")
			return
		}
		comm.WriteString("Working directory set\r")
		return
	}
	if linuxCommand == "a2version" {
		a2version()
		return
	}
	if linuxCommand == "a2help" {
		a2help()
		return
	}
	if linuxCommand == "a2lower" {
		a2lower(false)
		return
	}
	if linuxCommand == "A2LOWER" {
		a2lower(true)
		return
	}
	if linuxCommand == "a2wifi" {
		a2wifi()
		return
	}
	if strings.HasPrefix(linuxCommand, "a2timeout") {
		a2timeout(linuxCommand)
		return
	}
	if strings.HasPrefix(linuxCommand, "a2drive") {
		a2drive(linuxCommand, drive1, drive2)
		return
	}
	if linuxCommand == "a2prompt" {
		prompt := fmt.Sprintf("A2IO:%s ", workingDirectory)
		comm.WriteString(prompt)
		return
	}
	if linuxCommand == "a2wifi list" {
		linuxCommand = a2wifiList()
	}
	if strings.HasPrefix(linuxCommand, "a2wifi select") {
		linuxCommand, err = a2wifiSelect(linuxCommand)
	}
	if err == nil {
		execCommand(linuxCommand, workingDirectory)
	}
}

func execCommand(linuxCommand string, workingDirectory string) {
	// force the command to combine stderr(2) into stdout(1)
	linuxCommand += " 2>&1"
	cmd := exec.Command("bash", "-c", linuxCommand)
	cmd.Dir = workingDirectory
	cmd.Env = append(os.Environ(),
		"TERM=vt100",
		"LINES=24",
		"COLUMNS=80",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to set stdout\n")
		comm.WriteString("Failed to set stdout\r")
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start command\n")
		comm.WriteString("Failed to start command\r")
		return
	}

	timeout := time.After(time.Duration(execTimeoutSeconds) * time.Second)

	for {
		select {
		case <-timeout:
			comm.WriteString("\rCancelled by apple2driver\r")
			cmd.Process.Kill()
			return
		default:
			bb := make([]byte, 1)
			n, stdOutErr := stdout.Read(bb)
			if stdOutErr == nil {
				if n > 0 {
					b := bb[0]
					comm.SendCharacter(b)
				}
			} else {
				comm.WriteByte(0)
				cmd.Wait()
				return
			}
		}
	}
}

func a2version() {
	comm.WriteString("\rDriver version: " + info.Version + "\r")
}

func a2help() {
	comm.WriteString("\rDriver version: " + info.Version + "\r" +
		"\r" +
		"Example from ] prompt:\r" +
		"]RPI ls /home/pi\r" +
		"\r" +
		"Example from Applesoft BASIC:\r" +
		"]10 PRINT CHR$(4)\"RPI ping apple.com\"\r" +
		"]RUN\r" +
		"\r" +
		"Driver commands called with RPI:\r" +
		"a2help - display this message\r" +
		"a2version - display driver version\r" +
		"a2wifi - set up wifi\r" +
		"a2timeout - seconds to timeout commands\r" +
		"A2LOWER - force lowercase for II+\r" +
		"a2lower - disable force lowercase\r" +
		"a2drive - change drive images\r" +
		"\r")
}

func a2timeout(linuxCommand string) {
	params := strings.Fields(linuxCommand)

	switch len(params) {
	case 1:
		comm.WriteString("\rCommand timeout: " + strconv.FormatInt(int64(execTimeoutSeconds), 10) + "\r")
	case 2:
		timeoutSeconds, err := strconv.ParseInt(params[1], 10, 32)
		if err != nil {
			comm.WriteString("\rFailed to parse timeout\r")
		} else {
			execTimeoutSeconds = int(timeoutSeconds)
			comm.WriteString("\rCommand timeout set to: " + strconv.FormatInt(int64(execTimeoutSeconds), 10) + "\r")
		}
	default:
		comm.WriteString("\rToo many parameters\n")
	}
}

func a2drive(linuxCommand string, drive1 *prodos.ReaderWriterAt, drive2 *prodos.ReaderWriterAt) {
	params := strings.Fields(linuxCommand)

	if len(params) < 3 {
		showa2DriveUsage()
		return
	}

	driveNumber, err := strconv.ParseInt(params[1], 10, 32)
	if err != nil {
		comm.WriteString("\rFailed to parse drive number\r")
		showa2DriveUsage()
		return
	}

	if params[2] == "regen" {
		directory, err := drive.GetDriveImageDirectory()
		if err != nil {
			comm.WriteString("\rFailed to parse source directory\r")
			return
		}
		if len(params) > 3 {
			directory = params[3]
		}
		switch driveNumber {
		case 1:
			*drive1, err = drive.GenerateDriveFromDirectory("APPLE2.IO.RPI", directory)
			if err != nil {
				comm.WriteString("\rFailed to regenerate drive 1\r")
				return
			}
			comm.WriteString("\rDrive 1 regenerated\r")
		case 2:
			*drive2, err = drive.GenerateDriveFromDirectory("APPLE2.IO.RPI2", directory)
			if err != nil {
				comm.WriteString("\rFailed to regenerate drive 2\r")
				return
			}
			comm.WriteString("\rDrive 2 regenerated\r")
		default:
			comm.WriteString("\rOnly drives 1 or 2 are supported\r")
			showa2DriveUsage()
			return
		}
	}

	if params[2] == "load" {

		if len(params) != 4 {
			comm.WriteString("\rMust specify a drive image to load\r")
			showa2DriveUsage()
			return
		}

		imageFileName := params[3]

		switch driveNumber {
		case 1:
			*drive1, err = os.OpenFile(imageFileName, os.O_RDWR, 0755)
			if err != nil {
				comm.WriteString("\rFailed to load drive 1\r")
				return
			}
			comm.WriteString("\rDrive 1 loaded\r")
		case 2:
			*drive2, err = os.OpenFile(imageFileName, os.O_RDWR, 0755)
			if err != nil {
				comm.WriteString("\rFailed to load drive 2\r")
				return
			}
			comm.WriteString("\rDrive 2 loaded\r")
		default:
			comm.WriteString("\rOnly drives 1 or 2 are supported\r")
			showa2DriveUsage()
			return
		}
	}
}

func showa2DriveUsage() {
	comm.WriteString("\rUsage: a2drive DRIVENUMBER [regen [PATH] | load FILENAME]\rExamples: a2drive 1 regen ~/Apple2-IO-RPi/RaspberryPi/driveimage\r          a2drive 2 load /home/pi/Games.hdv\r")
}

func a2lower(enable bool) {
	forceLowercase = enable
	comm.WriteString(fmt.Sprintf("All commands will be converted to lowercase: %t\r", forceLowercase))
}

func a2wifi() {
	comm.WriteString("\r" +
		"Usage: a2wifi list\r" +
		"       a2wifi select SSID PASSWORD REGION\r" +
		"\r")
}

func a2wifiList() string {
	return "sudo iwlist wlan0 scanning | grep ESSID | sed s/.*ESSID://g | sed s/\\\"//g"
}

func a2wifiSelect(linuxCommand string) (string, error) {
	params := strings.Fields(linuxCommand)
	if len(params) != 5 {
		comm.WriteString("\rIncorrect number of parameters. Usage: a2wifi select SSID PASSWORD REGION\r\r")
		return "", errors.New("Incorrect number of parameters. Usage: a2wifi select SSID PASSWORD REGION")
	}
	ssid := params[2]
	psk := params[3]
	region := params[4]
	linuxCommand = "printf \"country=%s\\nupdate_config=1\\nctrl_interface=/var/run/wpa_supplicant\\n\\nnetwork={\\n  scan_ssid=1\\n  ssid=\\\"%s\\\"\n  psk=\\\"%s\\\"\\n}\\n\" " +
		region + " " +
		ssid + " " +
		psk + " " +
		" > /tmp/wpa_supplicant.conf; " +
		"sudo mv /tmp/wpa_supplicant.conf /etc/wpa_supplicant/; " +
		"sudo wpa_cli -i wlan0 reconfigure"
	return linuxCommand, nil
}
