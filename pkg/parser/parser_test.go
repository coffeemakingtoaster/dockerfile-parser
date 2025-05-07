package parser_test

import (
	"fmt"
	"testing"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/ast"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/parser"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
)

type TestCaseStage struct {
	Input    []token.Token
	Expected *ast.StageNode
}

type TestCaseInstruction struct {
	Input    []token.Token
	Expected []ast.InstructionNode
}

var baseImageLine = token.Token{
	Kind:    token.FROM,
	Content: "alpine:lastest",
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
	testCases := []TestCaseStage{
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

func TestAddParsing(t *testing.T) {
	testCases := []TestCaseInstruction{
		{
			Input: []token.Token{
				{
					Kind:    token.ADD,
					Content: "./source1 ./source2 ../../dest",
				},
			},
			Expected: []ast.InstructionNode{&ast.AddInstructionNode{
				Source:      []string{"./source1", "./source2"},
				Destination: "../../dest",
				KeepGitDir:  false,
				CheckSum:    "",
				Chown:       "",
				Chmod:       "",
				Link:        false,
				Exclude:     "",
			}},
		},
	}

	for _, c := range testCases {
		p := parser.NewParser(append([]token.Token{baseImageLine}, c.Input...))
		instructions := p.Parse().Instructions
		for i := range instructions {
			if instructions[i] != &c.Expected[i] {
				t.Errorf("Instruction comparison error: Expected %v Got %v", c.Expected[i], instructions[i])
			}

		}
	}
}
