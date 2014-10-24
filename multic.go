package main

import (
	"os"
	"os/user"
	"strings"

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
	o.Nl()
}

func loadConfig(path string) *config.Config {
	return config.NewConfig(expandPath(path))
}

func expandPath(path string) string {
	u, _ := user.Current()
	if u != nil && path[:2] == "~/" {
		path = strings.Replace(path, "~", u.HomeDir, 1)
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
		o := terminal.Stdout
		config := loadConfig(c.String("configuration"))
		if c.IsSet("list") {
			printDirectoryGroups(config)
		} else if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
		} else {
			n := c.String("group")
			g, err := config.GetDirectoryGroup(n)
			if err == nil {
				o.Colorf("@{g}Running command: ").Color("_").Print(c.Args()).Reset().Nl().Nl()
				printDirectoryGroup(n, g)
			} else {
				o.Color("r").Print(err).Reset().Nl()
			}
		}
	}
	app.Run(os.Args)
}
