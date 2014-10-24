package config_test

import (
	"github.com/eczarny/multic/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigParser", func() {
	It("should parse a directory group with a single directory", func() {
		c := config.ParseLines([]string{"PROJECTS=~/Projects"})
		Expect(c["PROJECTS"]).Should(ConsistOf("~/Projects"))
	})

	It("should parse multiple directory groups with a single directory", func() {
		c := config.ParseLines([]string{"PROJECTS=~/Projects", "DOCUMENTS=~/Documents"})
		Expect(c["PROJECTS"]).Should(ConsistOf("~/Projects"))
		Expect(c["DOCUMENTS"]).Should(ConsistOf("~/Documents"))
	})

	It("should parse multiple directory groups with a single directory, ignoring empty lines", func() {
		c := config.ParseLines([]string{"PROJECTS=~/Projects", "", " ", "DOCUMENTS=~/Documents", "	", "		"})
		Expect(c["PROJECTS"]).Should(ConsistOf("~/Projects"))
		Expect(c["DOCUMENTS"]).Should(ConsistOf("~/Documents"))
	})

	It("should parse a directory group with multiple directories", func() {
		l := []string{"projects=~/Projects/Spectacle,~/Projects/Go/src/github.com/eczarny/multic,~/Projects/Go/src/github.com/eczarny/lexer"}
		c := config.ParseLines(l)
		Expect(c["projects"]).Should(ConsistOf("~/Projects/Spectacle", "~/Projects/Go/src/github.com/eczarny/multic", "~/Projects/Go/src/github.com/eczarny/lexer"))
	})

	It("should parse multiple directory groups with multiple directories", func() {
		l := []string{
			"projects=~/Projects/Spectacle,~/Projects/Go/src/github.com/eczarny/lexer,~/Projects/Go/src/github.com/eczarny/multic",
			"go_projects=~/Projects/Go/src/github.com/eczarny/lexer,~/Projects/Go/src/github.com/eczarny/multic",
		}
		c := config.ParseLines(l)
		Expect(c["projects"]).Should(ConsistOf("~/Projects/Spectacle", "~/Projects/Go/src/github.com/eczarny/lexer", "~/Projects/Go/src/github.com/eczarny/multic"))
		Expect(c["go_projects"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny/lexer", "~/Projects/Go/src/github.com/eczarny/multic"))
	})

	It("should parse a directory group with a variable referencing a single directory", func() {
		l := []string{
			"PROJECTS=~/Projects",
			"GO_SRC=$PROJECTS/Go/src/github.com/eczarny",
		}
		c := config.ParseLines(l)
		Expect(c["PROJECTS"]).Should(ConsistOf("~/Projects"))
		Expect(c["GO_SRC"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny"))
	})

	It("should parse directory groups with multiple directories and variables referencing a single directory", func() {
		l := []string{
			"PROJECTS=~/Projects",
			"GO_SRC=$PROJECTS/Go/src/github.com/eczarny",
			"lexer=$GO_SRC/lexer",
			"multic=$GO_SRC/multic",
			"go_projects=$lexer,$multic",
		}
		c := config.ParseLines(l)
		Expect(c["PROJECTS"]).Should(ConsistOf("~/Projects"))
		Expect(c["GO_SRC"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny"))
		Expect(c["lexer"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny/lexer"))
		Expect(c["multic"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny/multic"))
		Expect(c["go_projects"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny/lexer", "~/Projects/Go/src/github.com/eczarny/multic"))
	})

	It("should parse a directory group with variables referencing multiple directories", func() {
		l := []string{
			"PROJECTS=~/Projects",
			"GO_SRC=$PROJECTS/Go/src/github.com/eczarny",
			"lexer=$GO_SRC/lexer",
			"multic=$GO_SRC/multic",
			"go_projects=$lexer,$multic",
			"default=$go_projects",
		}
		c := config.ParseLines(l)
		Expect(c["PROJECTS"]).Should(ConsistOf("~/Projects"))
		Expect(c["GO_SRC"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny"))
		Expect(c["lexer"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny/lexer"))
		Expect(c["multic"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny/multic"))
		Expect(c["go_projects"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny/lexer", "~/Projects/Go/src/github.com/eczarny/multic"))
		Expect(c["default"]).Should(ConsistOf("~/Projects/Go/src/github.com/eczarny/lexer", "~/Projects/Go/src/github.com/eczarny/multic"))
	})
})
