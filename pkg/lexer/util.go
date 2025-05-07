package lexer

import (
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
)

func splitFirstWord(input string) (string, string) {
	for i := range input {
		if input[i] == ' ' {
			return input[0:i], input[i:]
		}
	}
	return input, ""
}

func getParam(input string) (string, string, string, bool) {
	// TODO: clean this up
	key := ""
	value := ""
	ok := false
	endIndex := len(input) - 1
	currentWordStartIndex := 0
	seenCharacters := false
	for i := range input {
		if i+1 == len(input)-1 {
			break
		}
		if input[i] == ' ' && !ok {
			if !seenCharacters {
				continue
			}
			break
		}
		seenCharacters = true
		if input[i] == '-' && input[i+1] == '-' {
			ok = true
			i += 2
			currentWordStartIndex = i
			continue
		}
		if ok {
			if len(key) == 0 {
				if input[i] == '=' {
					key = input[currentWordStartIndex:i]
					currentWordStartIndex = i + 1
					continue
				}
				// Support params with no passed value
				if input[i] == ' ' {
					endIndex = i
					value = "true"
					key = input[currentWordStartIndex:i]
					break
				}

			} else {
				if input[i] == ' ' {
					endIndex = i
					value = input[currentWordStartIndex:i]
					break
				}
			}
		}
	}
	return key, value, input[endIndex:], ok
}

func buildToken(kind int, content string) token.Token {
	params := make(map[string]string)
	for {
		var key string
		var value string
		var ok bool
		key, value, content, ok = getParam(content)
		if !ok {
			break
		}
		params[key] = value
	}
	return token.Token{
		Kind:    kind,
		Params:  params,
		Content: content,
	}
}
