// Copyright Terence J. Boldt (c)2023-2024
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// This file is used for handling ProDOS image generation

package drive

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tjboldt/ProDOS-Utilities/prodos"
)

// GetDriveImageDirectory gets the default directory for driveimage
func GetDriveImageDirectory() (string, error) {
	exec, err := os.Executable()
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return "", err
	}
	driveImageDirectory := filepath.Join(filepath.Dir(filepath.Dir(exec)), "driveimage")

	return driveImageDirectory, nil
}

// GenerateDriveFromDirectory regenerates a ProDOS drive from a host directory
func GenerateDriveFromDirectory(volumeName string, directory string) (prodos.ReaderWriterAt, error) {
	drive := prodos.NewMemoryFile(0x2000000)
	fmt.Printf("Generating Drive in memory from: %s\n", directory)
	prodos.CreateVolume(drive, volumeName, 65535)
	err := prodos.AddFilesFromHostDirectory(drive, directory, "/"+volumeName+"/", true)
	return drive, err
}
