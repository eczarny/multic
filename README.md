# multic

[![Build Status](https://travis-ci.org/eczarny/multic.svg?branch=master)](https://travis-ci.org/eczarny/multic)

A utility that runs shell commands in multiple directories.

## Installing from Homebrew

multic is not yet distributed as a binary via [Homebrew][1]. Check back shortly!

## Installing from source

To install multic from source:

	$ go get github.com/eczarny/multic

Assuming that [Go][2] has been installed the `multic` binary should reside under `$GOPATH/bin/multic`. Add this to the `$PATH`.

# Usage

multic will run what ever command it is given in the configured group of directories, or directory groups.

By default multic will look for a `config` file under the `~/.multic` directory. The path to an alternative configuration file can be provided by specifying the `--configuration` or `-c` option.

Unless otherwise specified via the `--group` or `-g` option multic will use the _default_ directory group. Details on the configuration file format used by `multic` are provided in a subsequent section of the README file.

Refer to multic's help for additional options.

## Examples

multic will execute shell commands under each of the directories in a particular directory group. For each example assume the following configuration:

	PROJECTS=~/Projects
	GO_SRC=$PROJECTS/Go/src/github.com/eczarny
	lexer=$GO_SRC/lexer
	multic=$GO_SRC/multic
	go_projects=$lexer,$multic
	default=$go_projects

To list the contents of each directory in the default directory group:

	$ multic ls -la

multic will run the `ls -la` command in both `~/Projects/Go/src/github.com/eczarny/lexer` and `~/Projects/Go/src/github.com/eczarny/multic` directories (the `lexer` and `multic` directory groups, respectively).

Running a command in a particular directory group is just a bit more involved:

	$ multic -g multic ginkgo -r=true

# Configuring directory groups

Directory groups are declared in multic's configuration file. Each line declares a key and a list of values. Keys can then be referenced in subsequent lines as variables.

Take the multic configuration used in the examples above. The `PROJECTS` directory group declares one directory `~/Projects`. The `GO_SRC` directory group refers to `PROJECTS`, resulting in a directory group with one directory.

If directory groups consisting of multiple directories are referenced by another directory group all are added to that directory group. The `default` directory group is an example of this; `$go_projects` references a directory group including both `$lexer` and `$multic`. Upon evaluation multic will use `~/Projects/Go/src/github.com/eczarny/lexer` and `~/Projects/Go/src/github.com/eczarny/multic` when running commands in the `default` directory group.

To review the configured directory groups:

	$ multic -l

Specifying a `default` directory group is a good idea. If a default is not configured multic will require an explicit directory group every time it is used.

[1]: http://brew.sh/
[2]: http://golang.org/
