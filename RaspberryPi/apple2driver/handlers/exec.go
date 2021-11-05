// Copyright Terence J. Boldt (c)2020-2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is contains the handler for executing Linux and internal
// commands

package handlers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var forceLowercase = false

// ExecCommand handles requests for the Apple II executing Linux commands
func ExecCommand() {
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
	if linuxCommand == "a2help" {
		a2help()
		return
	}
	if linuxCommand == "A2LOWER" {
		a2lower()
		return
	}
	if linuxCommand == "a2wifi" {
		a2wifi()
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
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to set stdout\n")
		comm.WriteString("Failed to set stdout\r")
		return
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("Failed to set stdin\n")
		comm.WriteString("Failed to set stdin\r")
		return
	}

	fmt.Printf("Command output:\n")
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start command\n")
		comm.WriteString("Failed to start command\r")
		return
	}

	outputComplete := make(chan bool)
	inputComplete := make(chan bool)
	userCancelled := make(chan bool)

	if linuxCommand == "openssl" {
		fmt.Printf("\nSending help command...\n")
		io.WriteString(stdin, "help\n")
	}

	go getStdin(stdin, outputComplete, inputComplete, userCancelled)
	go getStdout(stdout, outputComplete, userCancelled)

	for {
		select {
		case <-outputComplete:
			outputComplete <- true
			cmd.Wait()
			comm.WriteByte(0)
			return
		case <-userCancelled:
			userCancelled <- true
			comm.WriteString("^C\r")
			cmd.Process.Kill()
			return
		case <-inputComplete:
			cmd.Wait()
			comm.WriteByte(0)
			return
		}
	}
}

func getStdout(stdout io.ReadCloser, outputComplete chan bool, userCancelled chan bool) {
	for {
		select {
		case <-userCancelled:
			stdout.Close()
			return
		default:
			bb := make([]byte, 1)
			n, err := stdout.Read(bb)
			if err != nil {
				stdout.Close()
				outputComplete <- true
				return
			}
			if n > 0 {
				b := bb[0]
				sendCharacter(b)
			}
		}
	}
}

func getStdin(stdin io.WriteCloser, done chan bool, inputComplete chan bool, userCancelled chan bool) {
	for {
		select {
		case <-done:
			stdin.Close()
			inputComplete <- true
			return
		default:
			b, err := comm.ReadByte()
			if err == nil {
				if b == 3 {
					stdin.Close()
					userCancelled <- true
					return
				}
				if b == 13 {
					b = 10
				}
				fmt.Printf("%c", b)
				io.WriteString(stdin, string(b))
			}
		}
	}
}

func a2help() {
	comm.WriteString("\r" +
		"This is a pseudo shell. Each command is executed as a process. The cd command\r" +
		"is intercepted and sets the working directory for the next command. The exit\r" +
		"command will exit the shell when not running from firmware.\r" +
		"\r" +
		"Built-in commands:\r" +
		"a2help - display this message\r" +
		"a2wifi - set up wifi\r" +
		"A2LOWER - force lowercase for II+\r" +
		"\r")
}

func a2lower() {
	forceLowercase = true
	comm.WriteString("All commands will be converted to lowercase\r")
}

func a2wifi() {
	comm.WriteString("\r" +
		"Usage: a2wifi list\r" +
		"       a2wifi select SSID PASSWORD\r" +
		"\r")
}

func a2wifiList() string {
	return "sudo iwlist wlan0 scanning | grep ESSID | sed s/.*ESSID://g | sed s/\\\"//g"
}

func a2wifiSelect(linuxCommand string) (string, error) {
	params := strings.Fields(linuxCommand)
	if len(params) != 4 {
		comm.WriteString("\rIncorrect number of parameters. Usage: a2wifi select SSID PASSWORD\r\r")
		return "", errors.New("Incorrect number of parameters. Usage: a2wifi select SSID PASSWORD")
	}
	ssid := params[2]
	psk := params[3]
	linuxCommand = "printf \"country=ca\\nupdate_config=1\\nctrl_interface=/var/run/wpa_supplicant\\n\\nnetwork={\\n  scan_ssid=1\\n  ssid=\\\"%s\\\"\n  psk=\\\"%s\\\"\\n}\\n\" " +
		ssid + " " +
		psk + " " +
		" > /tmp/wpa_supplicant.conf; " +
		"sudo mv /tmp/wpa_supplicant.conf /etc/wpa_supplicant/; " +
		"sudo wpa_cli -i wlan0 reconfigure"
	return linuxCommand, nil
}

func sendCharacter(b byte) {
	fmt.Print(string(b))
	if b == 10 { // convert LF to CR for Apple II compatiblity
		b = 13
	}
	if b == 9 { // convert TAB to spaces
		b = ' '
		b += 128
		comm.WriteByte(b)
	}
	b |= 128
	comm.WriteByte(b)
}
