package util_test

import (
	"reflect"
	"testing"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/util"
)

func TestAssignmentsParsing(t *testing.T) {
	input := "ABC=def test=\"test1 test2\" A=${SAMPLE:-placeholder}"
	expected := map[string]string{
		"ABC":  "def",
		"test": "\"test1 test2\"",
		"A":    "${SAMPLE:-placeholder}",
	}
	actual := util.ParseAssigns(input)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Parsing result mismatch: Expected %v Got %v", expected, actual)
	}
}
