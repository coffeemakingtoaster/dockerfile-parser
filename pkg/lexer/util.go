package lexer

import (
	"strings"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
)

func (l *Lexer) advanceWord() {
	for l.currentIndex < len(l.lines[l.currentLine]) {
		if l.lines[l.currentLine][l.currentIndex] == ' ' {
			break
		}
		l.currentIndex++
	}
}

func (l Lexer) expectNextCharacter(expected rune) bool {
	if l.currentIndex+1 >= len(l.lines[l.currentLine]) {
		return false
	}
	actual := l.lines[l.currentLine][l.currentIndex+1]
	return actual == byte(expected)
}

func (l Lexer) expectCurrentCharacter(expected rune) bool {
	if l.currentIndex >= len(l.lines[l.currentLine]) {
		return false
	}
	actual := l.lines[l.currentLine][l.currentIndex]
	return actual == byte(expected)
}

func (l Lexer) getCurrentCharacter() rune {
	return rune(l.lines[l.currentLine][l.currentIndex])
}

func (l *Lexer) advanceParam() (string, string, bool) {
	// TODO: clean this up
	key := ""
	value := ""
	ok := false
	endIndex := len(l.lines[l.currentLine]) - 1
	currentWordStartIndex := 0
	seenCharacters := false
	if l.currentIndex == endIndex {
		return key, value, ok
	}
	startIndex := l.currentIndex

	for l.currentIndex < endIndex {
		if l.expectCurrentCharacter(' ') && !ok {
			if !seenCharacters {
				l.currentIndex++
				continue
			}
			break
		}
		seenCharacters = true
		if l.expectCurrentCharacter('-') && l.expectNextCharacter('-') {
			ok = true
			l.currentIndex += 2
			currentWordStartIndex = l.currentIndex
			l.currentIndex++
			continue
		}
		if ok {
			if len(key) == 0 {
				if l.expectCurrentCharacter('=') {
					key = l.lines[l.currentLine][currentWordStartIndex:l.currentIndex]
					currentWordStartIndex = l.currentIndex + 1
					l.currentIndex++
					continue
				}
				// Support params with no passed value
				if l.expectCurrentCharacter(' ') {
					endIndex = l.currentIndex
					value = "true"
					key = l.lines[l.currentLine][currentWordStartIndex:l.currentIndex]
					break
				}

			} else {
				if l.expectCurrentCharacter(' ') {
					endIndex = l.currentIndex
					value = l.lines[l.currentLine][currentWordStartIndex:l.currentIndex]
					break
				}
			}
		}
		l.currentIndex++
	}
	if !ok {
		l.currentIndex = startIndex
	}
	return key, value, ok
}

func (l *Lexer) buildToken(kind int) token.Token {
	if kind == token.COMMENT {
		return token.Token{Kind: kind, Content: l.lines[l.currentLine][l.currentIndex:]}
	}
	params := make(map[string]string)
	for {
		key, value, ok := l.advanceParam()
		if !ok {
			break
		}
		params[key] = value
	}
	startIndex := l.currentIndex
	l.advanceToStartOfComment()
	return token.Token{
		Kind:          kind,
		Params:        params,
		Content:       l.lines[l.currentLine][startIndex:l.currentIndex],
		InlineComment: l.lines[l.currentLine][l.currentIndex:],
	}
}

func mergeLines(input []string) []string {
	target := []string{}
	buffer := ""

	for i := range input {
		in := strings.TrimSpace(input[i])
		if strings.HasPrefix(in, "#") {
			continue
		}
		buffer = buffer + in
		if strings.HasSuffix(in, "\\") {
			buffer = strings.TrimSuffix(buffer, "\\")
			continue
		}
		target = append(target, buffer)
		buffer = ""
	}
	return target
}
