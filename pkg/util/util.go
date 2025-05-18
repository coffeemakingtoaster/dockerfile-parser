package util

import (
	"bufio"
	"errors"
	"os"
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

type Stack[T comparable] struct {
	items []T
}

func (s Stack[T]) TopEquals(item T) bool {
	if i, err := s.Peek(); err == nil {
		return i == item
	}
	return false
}

func (s *Stack[T]) Peek() (T, error) {
	if len(s.items) == 0 {
		var res T
		return res, errors.New("No item to peek")
	}
	return s.items[len(s.items)-1], nil
}

func (s Stack[T]) Size() int {
	return len(s.items)
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, error) {
	res, err := s.Peek()
	if err != nil {
		return res, err
	}
	s.items = s.items[:len(s.items)-1]
	return res, nil
}
