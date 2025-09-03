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
	if len(input) == 0 {
		return res
	}
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
func CleanSlice(input []string) []string {
	result := []string{}
	for i := range input {
		cleanPath := strings.TrimSpace(input[i])
		if len(cleanPath) == 0 {
			continue
		}
		result = append(result, cleanPath)
	}
	return result
}
