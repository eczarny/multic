package main

import (
	"os"
	"os/exec"
)

func Run(dir string, name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
