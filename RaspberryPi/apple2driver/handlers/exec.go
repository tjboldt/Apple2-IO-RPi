package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"bufio"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
)

var forceLowercase = false

func ExecCommand() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		workingDirectory = "/home"
		a2io.WriteString("Failed to get current working directory, setting to /home\r")
	}

	fmt.Printf("Reading command to execute...\n")
	linuxCommand, err := a2io.ReadString()
	if forceLowercase {
		linuxCommand = strings.ToLower(linuxCommand)
	}
	fmt.Printf("Command to run: %s\n", linuxCommand)
	if strings.HasPrefix(linuxCommand, "cd ") {
		workingDirectory = strings.Replace(linuxCommand, "cd ", "", 1)
		err = os.Chdir(workingDirectory)
		if err != nil {
			a2io.WriteString("Failed to set working directory\r")
			return
		}
		a2io.WriteString("Working directory set\r")
		return
	}
	if linuxCommand == "a2help" {
		a2io.WriteString("\r" +
			"This is a pseudo shell. Each command is executed as a process. The cd command\r" +
			"is intercepted and sets the working directory for the next command. Running\r" +
			"commands that do not exit will hang. For example, do not use ping without a\r" +
			"way to limit output like -c 1.\r" +
			"\r" +
			"Built-in commands:\r" +
			"a2help - display this message\r" +
			"a2wifi - set up wifi\r" +
			"A2LOWER - force lowercase for II+\r" +
			"\r")
		return
	}
	if linuxCommand == "A2LOWER" {
		forceLowercase = true
		a2io.WriteString("All commands will be converted to lowercase\r")
		return
	}
	if linuxCommand == "a2wifi" {
		a2io.WriteString("\r" +
			"Usage: a2wifi list\r" +
			"       a2wifi select SSID PASSWORD\r" +
			"\r")
		return
	}
	if linuxCommand == "a2wifi list" {
		linuxCommand = "sudo iwlist wlan0 scanning | grep ESSID | sed s/.*ESSID://g | sed s/\\\"//g"
	}
	if strings.HasPrefix(linuxCommand, "a2wifi select") {
		params := strings.Fields(linuxCommand)
		if len(params) != 4 {
			a2io.WriteString("\rIncorrect number of parameters. Usage: a2wifi select SSID PASSWORD\r\r")
			return
		}
		ssid := params[2]
		psk := params[3]
		linuxCommand = "printf \"country=ca\\nupdate_config=1\\nctrl_interface=/var/run/wpa_supplicant\\n\\nnetwork={\\n  scan_ssid=1\\n  ssid=\\\"%s\\\"\n  psk=\\\"%s\\\"\\n}\\n\" " +
			ssid + " " +
			psk + " " +
			" > /tmp/wpa_supplicant.conf; " +
			"sudo mv /tmp/wpa_supplicant.conf /etc/wpa_supplicant/; " +
			"sudo wpa_cli -i wlan0 reconfigure"
	}
	cmd := exec.Command("bash", "-c", linuxCommand)
	cmd.Dir = workingDirectory
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to set stdout\n")
		a2io.WriteString("Failed to set stdout\r")
		return
	}
	fmt.Printf("Command output:\n")
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start command\n")
		a2io.WriteString("Failed to start command\r")
		return
	}

	reader := bufio.NewReader(stdout)

	for err == nil {
		var b byte
		b, err = reader.ReadByte()
		if err == nil {
			fmt.Print(string(b))
			if b == 10 { // convert LF to CR for Apple II compatiblity
				b = 13
			}
			b |= 128
			err = a2io.WriteByte(b)
			if err != nil {
				fmt.Printf("\nFailed to write byte\n")
				cmd.Process.Kill()
				return
			}
		}
	}

	cmd.Wait()
	a2io.WriteByte(0)
}
