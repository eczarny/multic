package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Config struct {
	path            string
	directoryGroups map[string][]string
}

func NewConfig(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	l, _ := readLines(f)
	c := &Config{
		path:            path,
		directoryGroups: ParseLines(l),
	}
	return c
}

func (c *Config) GetPath() string {
	return c.path
}

func (c *Config) GetDirectoryGroup(directoryGroupName string) []string {
	d, ok := c.directoryGroups[directoryGroupName]
	if !ok {
		panic(fmt.Sprintf("Directory group %s does not exists.", directoryGroupName))
	}
	return d
}

func (c *Config) DirectoryGroups() map[string][]string {
	return c.directoryGroups
}

func readLines(reader io.Reader) ([]string, error) {
	var r []string
	s := bufio.NewScanner(reader)
	for s.Scan() {
		r = append(r, s.Text())
	}
	return r, s.Err()
}
