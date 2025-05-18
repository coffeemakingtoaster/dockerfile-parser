package util

import "errors"

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
