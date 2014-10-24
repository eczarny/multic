package config

import (
	"bytes"
	"fmt"

	"github.com/eczarny/lexer"
)

type ConfigParser struct {
	config       map[string][]string
	lexer        *lexer.Lexer
	currentToken lexer.Token
}

func ParseLines(lines []string) map[string][]string {
	p := new(ConfigParser)
	p.config = make(map[string][]string)
	for _, l := range lines {
		p.lexer = NewConfigLexer(l)
		p.parse()
	}
	return p.config
}

func (p *ConfigParser) parse() {
	switch p.nextToken().Type {
	case TokenText:
		p.parseDirectoryGroup(p.tokenValue())
	case TokenEOF:
	default:
		panic("Unexpected input.")
	}
}

func (p *ConfigParser) parseDirectoryGroup(directoryGroupName string) {
	p.expectToken(TokenAssignment, "Directory group assignment expected.")
	d := p.parseDirectories(directoryGroupName)
	if len(d) > 0 {
		p.config[directoryGroupName] = d
	} else {
		panic(fmt.Sprintf("Directory group %s cannot be empty.", directoryGroupName))
	}
}

func (p *ConfigParser) parseDirectories(directoryGroupName string) []string {
	r := make([]string, 0)
Loop:
	for {
		switch p.currentToken.Type {
		case TokenText, TokenVariable:
			for _, d := range p.parseDirectoryOrVariable(directoryGroupName) {
				r = append(r, d)
			}
		case TokenAssignment:
		default:
			break Loop
		}
		if p.currentToken.Type != TokenEOF {
			p.nextToken()
		}
	}
	return r
}

func (p *ConfigParser) parseDirectoryOrVariable(directoryGroupName string) []string {
	var r []string
	switch p.currentToken.Type {
	case TokenText:
		r = []string{p.parseDirectory(directoryGroupName)}
	case TokenVariable:
		r = p.parseDirectoryVariable(directoryGroupName)
	}
	return r
}

func (p *ConfigParser) parseDirectory(directoryGroupName string) string {
	var r bytes.Buffer
Loop:
	for {
		switch p.currentToken.Type {
		case TokenText:
			r.WriteString(p.tokenValue())
		case TokenVariable:
			r.WriteString(p.lookupSingleDirectoryGroup(directoryGroupName, p.tokenValue()))
		case TokenEOF, TokenComma:
			break Loop
		}
		p.nextToken()
	}
	return r.String()
}

func (p *ConfigParser) parseDirectoryVariable(directoryGroupName string) []string {
	var r []string
Loop:
	for {
		switch p.currentToken.Type {
		case TokenText:
			r = append(r[:len(r)-1], r[len(r)-1] + p.tokenValue())
		case TokenVariable:
			r = p.lookupDirectoryGroup(directoryGroupName, p.tokenValue())
		case TokenEOF, TokenComma:
			break Loop
		}
		p.nextToken()
	}
	return r
}

func (p *ConfigParser) lookupDirectoryGroup(directoryGroupName, referencedDirectoryGroupName string) []string {
	c := p.config[referencedDirectoryGroupName]
	var r []string
	if len(c) == 0 {
		panic(fmt.Sprintf("%s references an invalid directory group: %s", directoryGroupName, referencedDirectoryGroupName))
	}
	r = make([]string, len(c))
	copy(r, c)
	return r
}

func (p *ConfigParser) lookupSingleDirectoryGroup(directoryGroupName, referencedDirectoryGroupName string) string {
	r := p.lookupDirectoryGroup(directoryGroupName, referencedDirectoryGroupName)
	if len(r) != 1 {
		panic(fmt.Sprintf("%s references an ambiguous directory group: %s", directoryGroupName, referencedDirectoryGroupName))
	}
	return r[0]
}

func (p *ConfigParser) tokenValue() string {
	return p.currentToken.Value.(string)
}

func (p *ConfigParser) nextToken() lexer.Token {
	p.currentToken = p.lexer.NextToken()
	return p.currentToken
}

func (p *ConfigParser) acceptToken(tokenType lexer.TokenType) bool {
	return p.nextToken().Type == tokenType
}

func (p *ConfigParser) expectToken(tokenType lexer.TokenType, v interface{}) string {
	if !p.acceptToken(tokenType) {
		panic(v)
	}
	return p.tokenValue()
}
