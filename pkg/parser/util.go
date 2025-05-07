package parser

import "strings"

func parsePossibleArray(input string) []string {
	cleanInput := strings.Trim(input, " ")
	if len(cleanInput) == 0 {
		return []string{}
	}
	if cleanInput[0] == '[' {
		panic("Arrays are not supported yet")
	}
	return strings.Split(cleanInput, " ")
}
