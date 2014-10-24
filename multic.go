package main

import (
	"os"
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

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = appHelpTemplate
	app.Name = "multic"
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
			Value: "~/.mc/config",
			Usage: "specify a configuration file",
		},
	}
	app.Action = func(c *cli.Context) {
		config := config.LoadConfig(c.String("configuration"))
		if c.IsSet("list") {
			printDirectoryGroups(config)
		} else if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
		} else {
			n := c.String("group")
			g, err := config.GetDirectoryGroup(n)
			if err == nil {
				for _, d := range g {
					printDirectorySeparator(d)
				}
			} else {
				terminal.Stderr.Color("r").Print(err).Reset().Nl()
			}
		}
	}
	app.Run(os.Args)
}
