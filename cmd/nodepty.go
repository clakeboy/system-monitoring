package main

import (
	"github.com/creack/pty"
	"os"
	"os/exec"
)

func GetPty() (*os.File, error) {
	// Create arbitrary command.
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return nil, err
	}

	return ptmx, nil
}
