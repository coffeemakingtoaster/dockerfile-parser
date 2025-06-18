package lexer

import (
	"strings"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/util"
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
		if k, _ := util.ParseAssign(l.lines[l.currentLine][l.currentIndex:]); len(k) != 0 {
			kind = token.PARSER_DIRECTIVE
		}
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
	if kind == token.RUN || kind == token.COPY {
		if l.containsHeredoc() {
			return l.buildHereDocToken(kind, params)
		}
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

func (l *Lexer) buildHereDocToken(kind int, params map[string]string) token.Token {
	// Get identifier and go until end
	heredocContent := []string{}
	encounteredRedirection := strings.HasPrefix(l.lines[l.currentLine][l.currentIndex:], "<<-")

	l.currentIndex += 2         // advance <<
	if encounteredRedirection { //advance -
		l.currentIndex++
	}

	heredocContent = append(heredocContent, l.lines[l.currentLine][l.currentIndex:])
	heredocStartIndex := l.currentIndex
	hasSeenCharacter := false
	// Parse delim
	for l.currentIndex < len(l.lines[l.currentLine]) {
		if l.expectCurrentCharacter(' ') {
			if hasSeenCharacter {
				break
			}
			l.currentIndex++
			heredocStartIndex++
			continue
		}
		hasSeenCharacter = true
		l.currentIndex++
	}
	delim := l.lines[l.currentLine][heredocStartIndex:l.currentIndex]
	delim = strings.ReplaceAll(strings.ReplaceAll(delim, "'", ""), "\"", "") // Delim definition may contain quotes...this is the easiest way to handle them for now

	l.currentLine++

	for l.currentLine < len(l.lines) {
		heredocContent = append(heredocContent, l.lines[l.currentLine])
		if strings.HasPrefix(l.lines[l.currentLine], delim) {
			break
		}
		l.currentLine++
	}

	return token.Token{
		Kind:               kind,
		Params:             params,
		MultiLineContent:   heredocContent,
		HereDocRedirection: encounteredRedirection,
	}
}

func (l *Lexer) containsHeredoc() bool {
	startIndex := l.currentIndex
	for l.currentIndex < len(l.lines[l.currentLine]) {
		if l.expectCurrentCharacter('<') && l.expectNextCharacter('<') {
			return true
		}
		if !l.expectCurrentCharacter(' ') {
			break
		}
		l.currentIndex++
	}
	l.currentIndex = startIndex
	return false
}

func mergeLines(input []string) []string {
	target := []string{}
	buffer := ""

	for i := range input {
		in := strings.TrimSpace(input[i])
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
