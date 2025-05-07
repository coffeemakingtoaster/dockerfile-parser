package parser_test

import (
	"fmt"
	"testing"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/ast"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/parser"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
)

type TestCase struct {
	Input    []token.Token
	Expected *ast.StageNode
}

func compareStageNodes(expected, actual ast.StageNode) string {
	if actual.Image != expected.Image {
		return fmt.Sprintf("Stage image mismatch: Expected %s Got %s", expected.Image, actual.Image)
	}

	if actual.Identifier != expected.Identifier {
		return fmt.Sprintf("Stage image mismatch: Expected %s Got %s", expected.Image, actual.Image)
	}
	return ""

}

func TestFromParsing(t *testing.T) {
	testCases := []TestCase{
		{
			Input: []token.Token{
				{
					Kind:    token.FROM,
					Content: "alpine:lastest",
				},
			},
			Expected: &ast.StageNode{
				Identifier: "anon",
				Image:      "alpine:lastest",
			},
		},
	}

	for _, c := range testCases {
		p := parser.NewParser(c.Input)
		err := compareStageNodes(*c.Expected, p.Parse())
		if err != "" {
			t.Error(err)
		}

	}
}
