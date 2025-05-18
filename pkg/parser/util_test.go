package parser

import (
	"testing"
)

func TestNoArrayParsing(t *testing.T) {
	input := "word ./path../file.txt testing 123"
	expected := []string{"word", "./path../file.txt", "testing", "123"}
	actual := parsePossibleArray(input)
	if len(expected) != len(actual) {
		t.Errorf("Parsing no array length mismatch: Expected %d Got %d", len(expected), len(actual))
	}
	for i := range actual {
		if expected[i] != actual[i] {
			t.Errorf("Parsing no array value (index: %d) mismatch: Expected %s Got %s", i, expected[i], actual[i])
		}
	}
}

func TestArrayParsing(t *testing.T) {
	input := "[\"word\", \"./path../file.txt\"  , \"testing\" ,\"123\"]"
	expected := []string{"word", "./path../file.txt", "testing", "123"}
	actual := parsePossibleArray(input)
	if len(expected) != len(actual) {
		t.Errorf("Parsing no array length mismatch: Expected %v Got %v", expected, actual)

		t.Errorf("Parsing no array length mismatch: Expected %d Got %d", len(expected), len(actual))
	}
	for i := range actual {
		if expected[i] != actual[i] {
			t.Errorf("Parsing no array value (index: %d) mismatch: Expected %s Got %s", i, expected[i], actual[i])
		}
	}
}

func TestEmptyArrayParsing(t *testing.T) {
	input := "[]"
	expected := []string{}
	actual := parsePossibleArray(input)
	if len(expected) != len(actual) {
		t.Errorf("Parsing no array length mismatch: Expected %v Got %v", expected, actual)

		t.Errorf("Parsing no array length mismatch: Expected %d Got %d", len(expected), len(actual))
	}
	for i := range actual {
		if expected[i] != actual[i] {
			t.Errorf("Parsing no array value (index: %d) mismatch: Expected %s Got %s", i, expected[i], actual[i])
		}
	}
}

func TestPaddedArrayParsing(t *testing.T) {
	input := "[ \"test\"   ]"
	expected := []string{"test"}
	actual := parsePossibleArray(input)
	if len(expected) != len(actual) {
		t.Errorf("Parsing no array length mismatch: Expected %v Got %v", expected, actual)

		t.Errorf("Parsing no array length mismatch: Expected %d Got %d", len(expected), len(actual))
	}
	for i := range actual {
		if expected[i] != actual[i] {
			t.Errorf("Parsing no array value (index: %d) mismatch: Expected %s Got %s", i, expected[i], actual[i])
		}
	}
}
