// Copyright Terence J. Boldt (c)2021
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is used for setting up the communications for the handlers

package handlers

import (
	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
)

var comm a2io.A2Io

// SetCommunication configures whether to use real or mock communications
func SetCommunication(commIO a2io.A2Io) {
	comm = commIO
}
