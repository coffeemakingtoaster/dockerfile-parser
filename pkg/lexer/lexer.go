// Parser package
package lexer

import (
	"fmt"
	"strings"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/util"
)

// Lexer
type Lexer struct {
	lines        []string
	currentLine  int
	currentIndex int
}

// Create new lexer based on the a file
// Errors if path is not a file or does not exist
func NewFromFile(path string) (Lexer, error) {
	lines, err := util.ReadFileLines(path)
	if err != nil {
		return Lexer{}, err
	}
	return Lexer{mergeLines(lines), 0, 0}, nil
}

// Create new lexer based on the input provided
func NewFromInput(input []string) Lexer {
	return Lexer{mergeLines(input), 0, 0}
}

// Lex lines provided when initializing lexer
// Technically this is both a lexer and a tokenizer in one
// Returns tokens
func (l *Lexer) Lex() ([]token.Token, error) {
	tokens := []token.Token{}
	for l.currentLine < len(l.lines) {
		instruction := l.getCurrentInstruction()
		switch instruction {
		case token.EOF:
			break
		case token.ILLEGAL:
			return tokens, fmt.Errorf("Illegal instruction encountered (line: %d) see above for details", l.currentLine)
		default:
			t := l.buildToken(instruction)
			tokens = append(tokens, t)
		}
		l.currentLine += 1
		l.currentIndex = 0
	}
	return tokens, nil
}

// Advance index to end of instruction and return token kind
func (l *Lexer) getCurrentInstruction() int {
	if l.currentLine == len(l.lines) {
		return token.EOF
	}
	currentLine := l.lines[l.currentLine]
	if len(currentLine) == 0 {
		return token.EMPTY_LINE
	}
	if currentLine[0] == '#' {
		l.currentIndex++
		return token.COMMENT
	}
	l.advanceWord()
	cmd := currentLine[:l.currentIndex]
	instruction, ok := token.TokenLookupTable[strings.ToUpper(cmd)]
	if ok {
		return instruction
	}
	fmt.Printf("Illegal statement: %s", currentLine)
	return token.ILLEGAL
}

// Advance index to end of content and start of comment
// If no comment exist -> advance to end of content
func (l *Lexer) advanceToStartOfComment() {
	stack := util.Stack[rune]{}
	for l.currentIndex < len(l.lines[l.currentLine]) {
		if l.expectCurrentCharacter('"') || l.expectCurrentCharacter('\'') || l.expectCurrentCharacter('`') || (l.expectCurrentCharacter('$') && l.expectNextCharacter('{')) {
			// Handle case of ${a#b} -> This does not count as a comment
			if l.expectCurrentCharacter('$') {
				l.currentIndex++
			}
			if stack.TopEquals(l.getCurrentCharacter()) {
				stack.Pop()
			} else {
				stack.Push(l.getCurrentCharacter())
			}
		}
		if l.expectCurrentCharacter('#') && stack.Size() == 0 {
			// Remove comment symbol
			l.currentIndex = min(l.currentIndex+1, len(l.lines[l.currentLine]))
			return
		}
		l.currentIndex++
	}
}
