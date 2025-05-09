package lexer

import (
	"reflect"
	"testing"
)

func TestMergeLine(t *testing.T) {
	input := []string{"do a \\", "#comment", "do B", "do C"}
	expected := []string{"do a do B", "do C"}
	actual := mergeLines(input)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Merged lines mismatch: Expected %+q Got %+q", expected, actual)
	}
}
