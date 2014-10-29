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
	config.WalkDirectoryGroups(func(name string, dirs []string) {
		printDirectoryGroup(name, dirs)
	})
}

func printDirectoryGroup(name string, dirs []string) {
	stdout := terminal.Stdout
	stdout.Colorf("@{.}Directory group: ").Reset()
	stdout.Color(".y").Print(name).Reset().Nl()
	for _, dir := range dirs {
		stdout.Print(dir).Nl()
	}
}

func printDirectorySeparator(dir string) {
	stdout := terminal.Stdout
	w, _, _ := terminalSize()
	stdout.Color(".").Print("\u2514 ").Color(".g").Print(dir).Color(".").Print(" ")
	for i := 1; i < w-len(dir)-2; i++ {
		stdout.Print("\u2500")
	}
	stdout.Reset().Nl()
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

func run(config *config.Config, groups []string, args cli.Args) {
	if len(groups) == 0 {
		groups = []string{"default"}
	}
	for _, group := range groups {
		dirs, err := config.GetDirectoryGroup(group)
		stderr := terminal.Stderr
		if err == nil {
			for _, dir := range dirs {
				err = Run(dir, args.First(), args.Tail())
				if err != nil {
					stderr.Color("r").Print(err).Reset().Nl()
				}
				printDirectorySeparator(dir)
			}
		} else {
			stderr.Color("r").Print(err).Reset().Nl()
		}
	}
}

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = appHelpTemplate
	app.Name = "multic"
	app.Version = "0.0.4"
	app.Usage = "Run shell commands in multiple directories."
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "list configured directory groups",
		},
		cli.StringSliceFlag{
			Name:  "group, g",
			Value: &cli.StringSlice{},
			Usage: "specify a directory group",
		},
		cli.StringFlag{
			Name:  "configuration, c",
			Value: "~/.multic/config",
			Usage: "specify a configuration file",
		},
	}
	app.Action = func(ctx *cli.Context) {
		args := ctx.Args()
		config := config.LoadConfig(ctx.String("configuration"))
		if ctx.IsSet("list") {
			printDirectoryGroups(config)
		} else if len(args) == 0 {
			cli.ShowAppHelp(ctx)
		} else {
			run(config, ctx.StringSlice("group"), args)
		}
	}
	app.Run(os.Args)
}
