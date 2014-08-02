package config_test

import (
	"github.com/eczarny/lexer"
	"github.com/eczarny/multic/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigLexer", func() {
	assertTokenType := func(token lexer.Token, tokenType lexer.TokenType) {
		Expect(token.Type).To(Equal(tokenType))
	}

	assertToken := func(token lexer.Token, tokenType lexer.TokenType, tokenValue interface{}) {
		Expect(token).To(Equal(lexer.Token{tokenType, tokenValue}))
	}

	assertEOF := func(token lexer.Token) {
		assertTokenType(token, config.TokenEOF)
	}

	It("should emit an EOF if the input is empty", func() {
		l := config.NewConfigLexer("")
		assertEOF(l.NextToken())
	})

	It("should emit a text token", func() {
		l := config.NewConfigLexer("PROJECTS")
		assertToken(l.NextToken(), config.TokenText, "PROJECTS")
		assertEOF(l.NextToken())
	})

	It("should emit a variable token", func() {
		l := config.NewConfigLexer("$PROJECTS")
		assertToken(l.NextToken(), config.TokenVariable, "PROJECTS")
		assertEOF(l.NextToken())
	})

	It("should emit an assignment token", func() {
		l := config.NewConfigLexer("=")
		assertTokenType(l.NextToken(), config.TokenAssignment)
		assertEOF(l.NextToken())
	})

	It("should emit a comma token", func() {
		l := config.NewConfigLexer(",")
		assertTokenType(l.NextToken(), config.TokenComma)
		assertEOF(l.NextToken())
	})

	It("should emit a text and assignment token (e.g. PROJECTS=)", func() {
		l := config.NewConfigLexer("PROJECTS=")
		assertToken(l.NextToken(), config.TokenText, "PROJECTS")
		assertTokenType(l.NextToken(), config.TokenAssignment)
		assertEOF(l.NextToken())
	})

	It("should emit text and assignment tokens (e.g. PROJECTS=~/Projects)", func() {
		l := config.NewConfigLexer("PROJECTS=~/Projects")
		assertToken(l.NextToken(), config.TokenText, "PROJECTS")
		assertTokenType(l.NextToken(), config.TokenAssignment)
		assertToken(l.NextToken(), config.TokenText, "~/Projects")
		assertEOF(l.NextToken())
	})

	It("should emit text and assignment tokens if the input contains spaces", func() {
		l := config.NewConfigLexer("  PROJECTS =    ~/Projects")
		assertToken(l.NextToken(), config.TokenText, "PROJECTS")
		assertTokenType(l.NextToken(), config.TokenAssignment)
		assertToken(l.NextToken(), config.TokenText, "~/Projects")
		assertEOF(l.NextToken())
	})

	It("should emit text and assignment tokens if the input contains tabs", func() {
		l := config.NewConfigLexer("	PROJECTS		=		~/Projects")
		assertToken(l.NextToken(), config.TokenText, "PROJECTS")
		assertTokenType(l.NextToken(), config.TokenAssignment)
		assertToken(l.NextToken(), config.TokenText, "~/Projects")
		assertEOF(l.NextToken())
	})

	It("should emit text, assignment, and variable tokens (e.g. GO_SRC=$PROJECTS/Go/src/github.com/eczarny)", func() {
		l := config.NewConfigLexer("GO_SRC=$PROJECTS/Go/src/github.com/eczarny")
		assertToken(l.NextToken(), config.TokenText, "GO_SRC")
		assertTokenType(l.NextToken(), config.TokenAssignment)
		assertToken(l.NextToken(), config.TokenVariable, "PROJECTS")
		assertToken(l.NextToken(), config.TokenText, "/Go/src/github.com/eczarny")
		assertEOF(l.NextToken())
	})

	It("should emit text, assignment, variable, and comma tokens (e.g. GO_PROJECTS=$GO_SRC/lexer,$GO_SRC/multic)", func() {
		l := config.NewConfigLexer("GO_PROJECTS=$GO_SRC/github.com/eczarny/lexer,$GO_SRC/github.com/eczarny/multic")
		assertToken(l.NextToken(), config.TokenText, "GO_PROJECTS")
		assertTokenType(l.NextToken(), config.TokenAssignment)
		assertToken(l.NextToken(), config.TokenVariable, "GO_SRC")
		assertToken(l.NextToken(), config.TokenText, "/github.com/eczarny/lexer")
		assertTokenType(l.NextToken(), config.TokenComma)
		assertToken(l.NextToken(), config.TokenVariable, "GO_SRC")
		assertToken(l.NextToken(), config.TokenText, "/github.com/eczarny/multic")
		assertEOF(l.NextToken())
	})

	It("should emit text, assignment, variable, and comma tokens (e.g. GO_PROJECTS=$lexer,$multic)", func() {
		l := config.NewConfigLexer("GO_PROJECTS=$lexer,$multic")
		assertToken(l.NextToken(), config.TokenText, "GO_PROJECTS")
		assertTokenType(l.NextToken(), config.TokenAssignment)
		assertToken(l.NextToken(), config.TokenVariable, "lexer")
		assertTokenType(l.NextToken(), config.TokenComma)
		assertToken(l.NextToken(), config.TokenVariable, "multic")
		assertEOF(l.NextToken())
	})
})
