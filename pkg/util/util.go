package util

import (
	"bufio"
	"os"
	"strings"
)

// Get value from passed map with a default
func GetFromParamsWithDefault(m map[string][]string, k string, d []string) []string {
	if v, ok := m[k]; ok {
		return v
	}
	return d
}

// Read the lines of a file into a slice
func ReadFileLines(path string) ([]string, error) {
	lines := []string{}

	file, err := os.Open(path)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

// TODO: Simplify this! For now this is GPT
func ParseAssigns(input string) map[string]string {
	m := make(map[string]string)
	parts := tokenize(input)

	var key string
	for i := 0; i < len(parts); i++ {
		p := parts[i]

		// Try normal KEY=VAL first
		k, v := ParseAssign(p)
		if k != "" && v != "" {
			m[k] = v
			continue
		}

		// Otherwise alternate KEY VAL
		if key == "" {
			key = p
		} else {
			m[key] = p
			key = ""
		}
	}
	if key != "" {
		m[key] = "" // dangling key
	}
	return m
}

func tokenize(input string) []string {
	var tokens []string
	var buf strings.Builder
	inQuotes := false

	for i := 0; i < len(input); i++ {
		c := input[i]

		switch c {
		case '"':
			inQuotes = !inQuotes
			buf.WriteByte(c) // keep the quote
		case ' ':
			if inQuotes {
				buf.WriteByte(c)
			} else {
				if buf.Len() > 0 {
					tokens = append(tokens, buf.String())
					buf.Reset()
				}
			}
		default:
			buf.WriteByte(c)
		}
	}

	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}
	return tokens
}

func ParseAssign(input string) (string, string) {
	for i := range input {
		if input[i] == '=' {
			return input[:i], input[i+1:]
		}
	}
	return "", ""
}
