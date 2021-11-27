// Copyright Terence J. Boldt (c)2020-2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is contains the handler for executing Linux shell

package handlers

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// ShellCommand handles requests for the Apple II executing a Linux shell
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

	go ptyIn(ptmx, outputComplete, inputComplete, userCancelled)
	go ptyOut(ptmx, outputComplete, userCancelled)

	for {
		select {
		case <-outputComplete:
			ptmx.Close()
			comm.WriteByte(0)
			return
		case <-userCancelled:
			fmt.Printf("User cancelled, killing process\n")
			ptmx.Close()
			cmd.Process.Kill()
			comm.WriteByte(0)
			return
		case <-inputComplete:
			cmd.Wait()
			comm.WriteByte(0)
			return
		}
	}
}

func ptyOut(stdout io.ReadCloser, outputComplete chan bool, userCancelled chan bool) {
	for {
		select {
		case <-userCancelled:
			fmt.Printf("User Cancelled stdout\n")
			stdout.Close()
			return
		default:
			bb := make([]byte, 1)
			n, err := stdout.Read(bb)
			if err != nil {
				stdout.Close()
				outputComplete <- true
				fmt.Printf("stdout closed\n")
				return
			}
			if n > 0 {
				b := bb[0]
				comm.SendCharacter(b)
			}
		}
	}
}

func ptyIn(stdin io.WriteCloser, done chan bool, inputComplete chan bool, userCancelled chan bool) {
	for {
		select {
		case <-done:
			stdin.Close()
			inputComplete <- true
			fmt.Printf("stdin closed\n")
			return
		default:
			s, err := comm.ReadCharacter()
			if err == nil {
				if s == string(byte(0x00)) {
					stdin.Close()
					userCancelled <- true
					fmt.Printf("\nUser cancelled stdin\n")
					return
				}
				io.WriteString(stdin, string(s))
			}
		}
	}
}
