package main

import (
	"os"
	"os/exec"
)

type Callback func(string, error)

func Run(dirs []string, name string, args []string, callback Callback) {
	for _, dir := range dirs {
		callback(dir, run(dir, name, args))
	}
}

func run(dir string, name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func directoryExists(path string) bool {
	s, err := os.Stat(path)
	return err == nil && s.IsDir()
}
