package util_test

import (
	"testing"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/util"
)

func TestStack(t *testing.T) {
	stack := util.Stack[string]{}
	_, err := stack.Peek()

	if err == nil {
		t.Error("Got no error for peeking empty stack")
	}

	_, err = stack.Pop()
	if err == nil {
		t.Error("Got no error for popping empty stack")
	}

	if stack.TopEquals("no") {
		t.Error("Topequals was true for empty stack...this should never happen")
	}

	stack.Push("1")
	stack.Push("2")
	stack.Push("3")

	if stack.Size() != 3 {
		t.Errorf("Size error! Wanted %d Got %d", 3, stack.Size())
	}

	if !stack.TopEquals("3") {
		t.Error("Top match did not work as intended")
	}

	stack.Pop()

	if !stack.TopEquals("2") {
		t.Error("Top match did not work as intended")
	}

	if stack.Size() != 2 {
		t.Errorf("Size error! Wanted %d Got %d", 2, stack.Size())
	}
}
