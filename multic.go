package main

import (
	"io"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/codegangsta/cli"
	"github.com/eczarny/multic/config"
	"github.com/wsxiaoys/terminal"
)

var appHelpTemplate = `{{.Name}}

   {{.Usage}}

Usage:

   {{.Name}} [options] [command] [command arguments...]

Options:

   {{range .Flags}}{{.}}
   {{end}}
`

func printDirectoryGroups(config *config.Config) {
	config.WalkDirectoryGroups(func(n string, g []string) {
		printDirectoryGroup(n, g)
	})
}

func printDirectoryGroup(directoryGroupName string, directoryGroup []string) {
	o := terminal.Stdout
	o.Colorf("@{.}Directory group: ").Reset()
	o.Color(".y").Print(directoryGroupName).Reset().Nl()
	for _, d := range directoryGroup {
		o.Print(d).Nl()
	}
}

func terminalSize() (int, int, error) {
	var d [4]uint16
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(0),
		uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&d)), 0, 0, 0)
	if err != 0 {
		return -1, -1, err
	}
	return int(d[1]), int(d[0]), nil
}

func printDirectorySeparator(directory string) {
	o := terminal.Stdout
	w, _, _ := terminalSize()
	o.Color(".").Print("\u2514 ").Color(".g").Print(directory).Color(".").Print(" ")
	for i := 1; i < w-len(directory)-2; i++ {
		o.Print("\u2500")
	}
	o.Reset().Nl()
}

func run(dir string, stdin io.Reader, stdout, stderr io.Writer, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func directoryExists(path string) bool {
	s, err := os.Stat(path)
	return err == nil && s.IsDir()
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

func runCommandInDirectories(config *config.Config, directoryGroupName string, args cli.Args) {
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

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = appHelpTemplate
	app.Name = "multic"
	app.Version = "0.0.2"
	app.Usage = "Run shell commands in multiple directories."
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "list configured directory groups",
		},
		cli.StringFlag{
			Name:  "group, g",
			Value: "default",
			Usage: "specify a directory group",
		},
		cli.StringFlag{
			Name:  "configuration, c",
			Value: "~/.multic/config",
			Usage: "specify a configuration file",
		},
	}
	app.Action = func(c *cli.Context) {
		args := c.Args()
		config := config.LoadConfig(c.String("configuration"))
		if c.IsSet("list") {
			printDirectoryGroups(config)
		} else if len(args) == 0 {
			cli.ShowAppHelp(c)
		} else {
			runCommandInDirectories(config, c.String("group"), args)
		}
	}
	app.Run(os.Args)
}
