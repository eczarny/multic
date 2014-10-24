package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

type Config struct {
	path                string
	directoryGroupNames []string
	directoryGroups     map[string][]string
}

type StepFunc func(string, []string)

func LoadConfig(path string) *Config {
	f, err := os.Open(expandPath(path))
	if err != nil {
		return nil
	}
	defer f.Close()
	l, _ := readLines(f)
	k, d := ParseLines(l)
	c := &Config{
		path:                path,
		directoryGroupNames: k,
		directoryGroups:     d,
	}
	return c
}

func (c *Config) GetPath() string {
	return c.path
}

func (c *Config) GetDirectoryGroup(directoryGroupName string) ([]string, error) {
	d, ok := c.directoryGroups[directoryGroupName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Directory group '%s' does not exists.", directoryGroupName))
	}
	return d, nil
}

func (c *Config) WalkDirectoryGroups(step StepFunc) {
	for _, n := range c.directoryGroupNames {
		g, err := c.GetDirectoryGroup(n)
		if err == nil {
			step(n, g)
		}
	}
}

func (c *Config) DirectoryGroups() map[string][]string {
	return c.directoryGroups
}

func expandPath(path string) string {
	u, _ := user.Current()
	if u != nil && path[:2] == "~/" {
		path = strings.Replace(path, "~", u.HomeDir, 1)
	}
	return path
}

func readLines(reader io.Reader) ([]string, error) {
	var r []string
	s := bufio.NewScanner(reader)
	for s.Scan() {
		r = append(r, s.Text())
	}
	return r, s.Err()
}
