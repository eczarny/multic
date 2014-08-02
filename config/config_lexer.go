package config

import (
	"unicode"

	. "github.com/eczarny/lexer"
)

const (
	TokenText TokenType = iota
	TokenVariable
	TokenAssignment
	TokenComma
	TokenEOF
)

func NewConfigLexer(input string) *Lexer {
	return NewLexer(input, initialState)
}

func initialState(l *Lexer) StateFunction {
	r := l.IgnoreUpTo(func(r rune) bool {
		return variable(r) || assignment(r) || comma(r) || nonWhitespace(r)
	})
	switch {
	case variable(r):
		return variableState
	case assignment(r):
		return assignmentState
	case comma(r):
		return commaState
	case nonWhitespace(r) && r != EOF:
		return textState
	}
	l.Emit(TokenEOF)
	return nil
}

func textState(l *Lexer) StateFunction {
	l.NextUpTo(func(r rune) bool {
		return variable(r) || assignment(r) || comma(r) || whitespace(r)
	})
	l.Emit(TokenText)
	return initialState
}

func variableState(l *Lexer) StateFunction {
	l.Ignore()
	l.NextUpTo(nonAlphanumeric)
	l.Emit(TokenVariable)
	return initialState
}

func assignmentState(l *Lexer) StateFunction {
	l.Ignore()
	l.Emit(TokenAssignment)
	return initialState
}

func commaState(l *Lexer) StateFunction {
	l.Ignore()
	l.Emit(TokenComma)
	return initialState
}

func alphanumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func nonAlphanumeric(r rune) bool {
	return !alphanumeric(r)
}

func whitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func nonWhitespace(r rune) bool {
	return !whitespace(r)
}

func assignment(r rune) bool {
	return r == '='
}

func variable(r rune) bool {
	return r == '$'
}

func comma(r rune) bool {
	return r == ','
}
