// Copyright Terence J. Boldt (c)2020-2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is contains the handler for executing Linux and internal
// commands

package handlers

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// ExecCommand handles requests for the Apple II executing Linux commands
func ShellCommand() {
	cmd := exec.Command("bash", "-i")
	cmd.Env = append(os.Environ(),
		"TERM=vt100",
		"LINES=24",
		"COLUMNS=80",
	)

	var ws pty.Winsize
	ws.Cols = 80
	ws.Rows = 24
	ws.X = 0
	ws.Y = 0

	ptmx, _ := pty.StartWithSize(cmd, &ws)
	defer func() { _ = ptmx.Close() }()

	outputComplete := make(chan bool)
	inputComplete := make(chan bool)
	userCancelled := make(chan bool)

	go getStdin(ptmx, outputComplete, inputComplete, userCancelled)
	go getStdout(ptmx, outputComplete, userCancelled)

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
