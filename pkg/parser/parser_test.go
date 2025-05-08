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
		return fmt.Sprintf("Stage identifier mismatch: Expected %s Got %s", expected.Identifier, actual.Identifier)
	}
	return ""
}

func compareAddInstructionNode(expected, actual ast.AddInstructionNode) string {
	if expected.CheckSum != actual.CheckSum {
		return fmt.Sprintf("Add instruction checksum param mismatch: Expected %s Got %s", expected.CheckSum, actual.CheckSum)
	}
	if expected.Link != actual.Link {
		return fmt.Sprintf("Add instruction link param mismatch: Expected %v Got %v", expected.Link, actual.Link)
	}
	if expected.KeepGitDir != actual.KeepGitDir {
		return fmt.Sprintf("Add instruction keep git dir param mismatch: Expected %v Got %v", expected.KeepGitDir, actual.KeepGitDir)
	}
	if expected.Chmod != actual.Chmod {
		return fmt.Sprintf("Add instruction chmod param mismatch: Expected %s Got %s", expected.Chmod, actual.Chmod)
	}
	if expected.Chown != actual.Chown {
		return fmt.Sprintf("Add instruction chown param mismatch: Expected %s Got %s", expected.Chown, actual.Chown)
	}
	if expected.Exclude != actual.Exclude {
		return fmt.Sprintf("Add instruction exclude param mismatch: Expected %s Got %s", expected.Exclude, actual.Exclude)
	}
	if expected.Destination != actual.Destination {
		return fmt.Sprintf("Add instruction destination mismatch: Expected %s Got %s", expected.Destination, actual.Destination)
	}
	if len(expected.Source) != len(actual.Source) {
		return fmt.Sprintf("Add instruction source length mismatch: Expected %d Got %d", len(expected.Source), len(actual.Source))
	}
	for i := range expected.Source {
		if expected.Source[i] != actual.Source[i] {
			return fmt.Sprintf("Add isntruction source mismatch: Expected %v Got %v", expected.Source, actual.Source)
		}
	}
	return ""
}

func compareInstructionNode(expected, actual ast.InstructionNode) string {
	switch ac := actual.(type) {
	case ast.AddInstructionNode:
		return compareAddInstructionNode(expected.(ast.AddInstructionNode), ac)
	default:
		return "Unknown ast node type"
	}
}

func TestFromParsing(t *testing.T) {
	testCases := []TestCaseStage{
		{
			Input: []token.Token{
				{
					Kind:    token.FROM,
					Content: "alpine:latest",
				},
			},
			Expected: &ast.StageNode{
				Identifier: "anon",
				Image:      "alpine:latest",
			},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.FROM,
					Content: "alpine:latest AS base",
				},
			},
			Expected: &ast.StageNode{
				Identifier: "base",
				Image:      "alpine:latest",
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
			Expected: []ast.InstructionNode{ast.AddInstructionNode{
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
			err := compareInstructionNode(instructions[i], c.Expected[i])
			if err != "" {
				t.Error(err)
			}

		}
	}
}
