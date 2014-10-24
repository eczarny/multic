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
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	lines, _ := readLines(file)
	c := &Config{
		path:            path,
		directoryGroups: ParseLines(lines),
	}
	return c
}

func (c *Config) GetPath() string {
	return c.path
}

func (c *Config) GetDirectoryGroup(directoryGroupName string) []string {
	directoryGroup, ok := c.directoryGroups[directoryGroupName]
	if !ok {
		panic(fmt.Sprintf("Directory group %s does not exists.", directoryGroupName))
	}
	return directoryGroup
}

func (c *Config) DirectoryGroups() map[string][]string {
	return c.directoryGroups
}

func readLines(reader io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
