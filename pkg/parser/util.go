package parser

import (
	"strings"
)

func parsePossibleArray(input string) []string {
	cleanInput := strings.Trim(input, " ")
	if len(cleanInput) == 0 {
		return []string{}
	}
	if cleanInput[0] == '[' {
		return parseConfirmedArray(cleanInput)
	}
	return strings.Split(cleanInput, " ")
}

func parseConfirmedArray(input string) []string {
	// Format of ["abc", "def"]
	res := []string{}
	wordStart := 0
	input = strings.Trim(input, "[")
	input = strings.Trim(input, "]")
	for i := range input {
		if input[i] == ',' {
			cur := input[wordStart:i]
			cur = strings.Trim(cur, " ")
			cur = strings.Trim(cur, "\"")
			res = append(res, cur)
			wordStart = i + 1
		}
	}
	cur := input[wordStart : len(input)-1]
	cur = strings.Trim(cur, " ")
	cur = strings.Trim(cur, "\"")
	res = append(res, cur)
	return res
}

func parseAssign(input string) (string, string) {
	for i := range input {
		if input[i] == '=' {
			return input[:i], input[i+1:]
		}
	}
	return "", ""
}
