package config

import (
	"bytes"

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
	for _, line := range lines {
		p.lexer = NewConfigLexer(line)
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
	directories := p.parseDirectories()
	if len(directories) > 0 {
		p.config[directoryGroupName] = directories
	} else {
		panic("Empty directory group invalid.")
	}
}

func (p *ConfigParser) parseDirectories() []string {
	directories := make([]string, 0)
	for {
		switch p.currentToken.Type {
		case TokenText, TokenVariable:
			directories = append(directories, p.parseDirectory())
		case TokenAssignment:
		default:
			return directories
		}
		if p.currentToken.Type != TokenEOF {
			p.nextToken()
		}
	}
}

func (p *ConfigParser) parseDirectory() string {
	var directory bytes.Buffer
	for {
		switch p.currentToken.Type {
		case TokenText:
			directory.WriteString(p.tokenValue())
		case TokenVariable:
			directory.WriteString(p.config[p.tokenValue()][0])
		case TokenEOF, TokenComma:
			return directory.String()
		}
		p.nextToken()
	}
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
