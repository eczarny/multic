package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/eczarny/multic/config"
)

var appHelpTemplate = `{{.Name}}

   {{.Usage}}

Usage:

   {{.Name}} [options] [arguments...]

Options:

   {{range .Flags}}{{.}}
   {{end}}
`

func printDirectoryGroups(directoryGroups map[string][]string) {
	for directoryGroupName, directoryGroup := range directoryGroups {
		fmt.Printf("Directory group: %s\n", directoryGroupName)
		for _, directory := range directoryGroup {
			fmt.Printf("    %s\n", directory)
		}
		fmt.Println()
	}
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
		config := config.NewConfig("/Users/eczarny/.mc/config")
		if c.IsSet("list") {
			printDirectoryGroups(config.DirectoryGroups())
		} else {
			cli.ShowAppHelp(c)
		}
	}
	app.Run(os.Args)
}
