package util

import (
	"bufio"
	"os"
	"strings"
)

// Get value from passed map with a default
func GetFromParamsWithDefault(m map[string]string, k string, d string) string {
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

func ParseAssigns(input string) map[string]string {
	m := make(map[string]string)
	parts := strings.Split(input, " ")
	var key string
	for _, p := range parts {
		k, v := ParseAssign(p)
		// If there is a key but the next assignment could not be parsed:
		// This should mean that the assigned value uses " and contains a space -> attach to previous key
		if v == "" {
			if key != "" {
				m[key] = m[key] + " " + p
			}
			continue
		}
		key = k
		m[key] = v
	}
	return m
}

func ParseAssign(input string) (string, string) {
	for i := range input {
		if input[i] == '=' {
			return input[:i], input[i+1:]
		}
	}
	return "", ""
}
