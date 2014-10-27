package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/eczarny/multic/config"
	"github.com/wsxiaoys/terminal"
)

func RunCommandInDirectories(config *config.Config, directoryGroupName string, args cli.Args) {
	g, err := config.GetDirectoryGroup(directoryGroupName)
	if err == nil {
		for _, d := range g {
			runCommand(d, args)
			printDirectorySeparator(d)
		}
	} else {
		terminal.Stderr.Color("r").Print(err).Reset().Nl()
	}
}

func runCommand(dir string, args cli.Args) {
	if directoryExists(dir) {
		err := run(dir, os.Stdin, os.Stdout, os.Stderr, args.First(), args.Tail()...)
		if err != nil {
			terminal.Stderr.Color("r").Print(err).Reset().Nl()
		}
	} else {
		terminal.Stderr.Colorf("@{r}The directory does not exist: ").Print(dir).Nl()
	}
}

func directoryExists(path string) bool {
	s, err := os.Stat(path)
	return err == nil && s.IsDir()
}

func run(dir string, stdin io.Reader, stdout, stderr io.Writer, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
