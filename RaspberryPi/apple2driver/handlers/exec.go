package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
)

func ExecCommand() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		workingDirectory = "/home"
		a2io.WriteString("Failed to get current working directory, setting to /home\r")
	}

	fmt.Printf("Reading command to execute...\n")
	linuxCommand, err := a2io.ReadString()
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
			"\r")
		return
	}
	cmd := exec.Command("bash", "-c", linuxCommand)
	cmd.Dir = workingDirectory
	cmdOut, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to execute command\n")
		a2io.WriteString("Failed to execute command\r")
		return
	}
	fmt.Printf("Command output: %s\n", cmdOut)
	apple2string := strings.Replace(string(cmdOut), "\n", "\r", -1)
	err = a2io.WriteString(apple2string)
	if err != nil {
		fmt.Printf("Failed to send command output\n")
		return
	}
}
