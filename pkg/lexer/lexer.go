// Parser package
package lexer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/util"
)

// Lexer
type Lexer struct {
	lines       []string
	currentLine int
}

// Create new lexer based on the a file
// Errors if path is not a file or does not exist
func NewFromFile(path string) (Lexer, error) {
	lines, err := util.ReadFileLines(path)
	if err != nil {
		return Lexer{}, err
	}
	return Lexer{mergeLines(lines), 0}, nil
}

// Create new lexer based on the input provided
func NewFromInput(input []string) Lexer {
	return Lexer{mergeLines(input), 0}
}

// Lex lines provided when initializing lexer
// Returns tokens
func (l *Lexer) Lex() ([]token.Token, error) {
	tokens := []token.Token{}
	for l.currentLine < len(l.lines) {
		instruction, content := l.getCurrentInstruction()
		switch instruction {
		case token.EOF:
			break
		case token.ILLEGAL:
			return tokens, errors.New(fmt.Sprintf("Illegal instruction encountered (line: %d)", l.currentLine))
		default:
			t := buildToken(instruction, content)
			tokens = append(tokens, t)
		}
		l.currentLine += 1
	}
	return tokens, nil
}

func (l *Lexer) Reset() {
	l.currentLine = 0
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
		strippedLine := strings.Trim(currentLine[1:], " ")
		return token.COMMENT, strippedLine
	}
	cmd, content := splitFirstWord(currentLine)
	instruction, ok := token.TokenLookupTable[strings.ToUpper(cmd)]
	if ok {
		return instruction, content
	}
	return token.ILLEGAL, ""
}
