package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/eczarny/multic/config"
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
	fmt.Printf("Loaded directory group configuration: %s\n\n", config.GetPath())
	for directoryGroupName, directoryGroup := range config.DirectoryGroups() {
		printDirectoryGroup(directoryGroupName, directoryGroup)
	}
}

func printDirectoryGroup(directoryGroupName string, directoryGroup []string) {
	fmt.Printf("Directory group: %s\n", directoryGroupName)
	for _, directory := range directoryGroup {
		fmt.Printf("    %s\n", directory)
	}
	fmt.Println()
}

func loadConfig(path string) *config.Config {
	return config.NewConfig(expandPath(path))
}

func expandPath(path string) string {
	user, _ := user.Current()
	if user != nil && path[:2] == "~/" {
		path = strings.Replace(path, "~", user.HomeDir, 1)
	}
	return path
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
		config := loadConfig(c.String("configuration"))
		if c.IsSet("list") {
			printDirectoryGroups(config)
		} else if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
		} else {
			directoryGroupName := c.String("group")
			fmt.Printf("Running command: %s\n\n", c.Args())
			printDirectoryGroup(directoryGroupName, config.GetDirectoryGroup(directoryGroupName))
		}
	}
	app.Run(os.Args)
}
