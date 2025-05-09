package lexer

import (
	"fmt"
	"strings"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
)

type Lexer struct {
	lines       []string
	currentLine int
}

func New(input []string) Lexer {
	return Lexer{mergeLines(input), 0}
}

func (l *Lexer) Lex() []token.Token {
	tokens := []token.Token{}
	for l.currentLine < len(l.lines) {
		instruction, content := l.getCurrentInstruction()
		switch instruction {
		case token.EOF:
			break
		case token.ILLEGAL:
			fmt.Println("Illegal instruction encountered")
			//panic("Illegal instruction")
		default:
			t := buildToken(instruction, content)
			tokens = append(tokens, t)
		}
		l.currentLine += 1
	}
	return tokens
}

func (l Lexer) panic(msg string) {
	panic(fmt.Sprintf("Error: %s at line '%s'", msg, l.lines[l.currentLine]))
}

// TODO: Migrate to a system of a pointer within the line rather than passing the remaining content around
func (l *Lexer) getCurrentInstruction() (int, string) {
	if l.currentLine == len(l.lines) {
		return token.EOF, ""
	}
	currentLine := l.lines[l.currentLine]
	if len(currentLine) == 0 {
		l.currentLine += 1
		return l.getCurrentInstruction()
	}
	if currentLine[0] == '#' {
		l.currentLine += 1
		return l.getCurrentInstruction()
	}
	cmd, content := splitFirstWord(currentLine)
	instruction, ok := token.TokenLookupTable[strings.ToUpper(cmd)]
	if ok {
		return instruction, content
	}
	return token.ILLEGAL, ""
}
