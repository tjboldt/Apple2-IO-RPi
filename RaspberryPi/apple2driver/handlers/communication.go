package handlers

import (
	"github.com/tjboldt/Apple2-IO-RPi/RaspberryPi/apple2driver/a2io"
)

var comm a2io.A2Io

func SetCommunication(commIO a2io.A2Io) {
	comm = commIO
}
