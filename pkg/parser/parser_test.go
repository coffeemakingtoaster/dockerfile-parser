package parser_test

import (
	"fmt"
	"reflect"
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

func compareAddInstructionNode(expected, actual *ast.AddInstructionNode) string {
	if expected.CheckSum != actual.CheckSum {
		return fmt.Sprintf("ADD instruction checksum param mismatch: Expected %s Got %s", expected.CheckSum, actual.CheckSum)
	}
	if expected.Link != actual.Link {
		return fmt.Sprintf("ADD instruction link param mismatch: Expected %v Got %v", expected.Link, actual.Link)
	}
	if expected.KeepGitDir != actual.KeepGitDir {
		return fmt.Sprintf("ADD instruction keep git dir param mismatch: Expected %v Got %v", expected.KeepGitDir, actual.KeepGitDir)
	}
	if expected.Chmod != actual.Chmod {
		return fmt.Sprintf("ADD instruction chmod param mismatch: Expected %s Got %s", expected.Chmod, actual.Chmod)
	}
	if expected.Chown != actual.Chown {
		return fmt.Sprintf("ADD instruction chown param mismatch: Expected %s Got %s", expected.Chown, actual.Chown)
	}
	if expected.Exclude != actual.Exclude {
		return fmt.Sprintf("ADD instruction exclude param mismatch: Expected %s Got %s", expected.Exclude, actual.Exclude)
	}
	if expected.Destination != actual.Destination {
		return fmt.Sprintf("ADD instruction destination mismatch: Expected %s Got %s", expected.Destination, actual.Destination)
	}
	if !reflect.DeepEqual(expected.Source, actual.Source) {
		return fmt.Sprintf("ADD instruction source mismatch: Expected %v Got %v", expected.Source, actual.Source)
	}
	return ""
}

func compareCopyInstructionNode(expected, actual *ast.CopyInstructionNode) string {
	if expected.Link != actual.Link {
		return fmt.Sprintf("COPY instruction link param mismatch: Expected %v Got %v", expected.Link, actual.Link)
	}
	if expected.KeepGitDir != actual.KeepGitDir {
		return fmt.Sprintf("COPY instruction keep git dir param mismatch: Expected %v Got %v", expected.KeepGitDir, actual.KeepGitDir)
	}
	if expected.Chown != actual.Chown {
		return fmt.Sprintf("COPY instruction chown param mismatch: Expected %s Got %s", expected.Chown, actual.Chown)
	}
	if expected.Destination != actual.Destination {
		return fmt.Sprintf("COPY instruction destination mismatch: Expected %s Got %s", expected.Destination, actual.Destination)
	}
	if !reflect.DeepEqual(expected.Source, actual.Source) {
		return fmt.Sprintf("COPY instruction source mismatch: Expected %v Got %v", expected.Source, actual.Source)
	}
	return ""
}

func compareInstructionNode(expected, actual ast.InstructionNode) string {
	switch ac := actual.(type) {
	case *ast.AddInstructionNode:
		return compareAddInstructionNode(expected.(*ast.AddInstructionNode), ac)
	case *ast.ArgInstructionNode:
		if *expected.(*ast.ArgInstructionNode) != *ac {
			return fmt.Sprintf("ARG instruction mismatch: Expected %v Got %v", expected, ac)
		}
	case *ast.CmdInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.CmdInstructionNode).Cmd, ac.Cmd) {
			return fmt.Sprintf("CMD instruction command mismatch: Expected %v Got %v", expected.(*ast.CmdInstructionNode).Cmd, ac.Cmd)
		}
	case *ast.CopyInstructionNode:
		return compareCopyInstructionNode(expected.(*ast.CopyInstructionNode), ac)
	case *ast.EntrypointInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.EntrypointInstructionNode).Exec, ac.Exec) {
			return fmt.Sprintf("ENTRYPOINT instruction command mismatch: Expected %v Got %v", expected.(*ast.EntrypointInstructionNode).Exec, ac.Exec)
		}
	default:
		return "Unknown ast node type"
	}
	return ""
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

func TestInstructionParsing(t *testing.T) {
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
		{
			Input: []token.Token{
				{
					Kind:    token.ARG,
					Content: "test=value",
				},
			},
			Expected: []ast.InstructionNode{&ast.ArgInstructionNode{
				Name:  "test",
				Value: "value",
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.CMD,
					Content: "echo hello testing",
				},
			},
			Expected: []ast.InstructionNode{&ast.CmdInstructionNode{
				Cmd: []string{"echo", "hello", "testing"},
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.COPY,
					Content: "./source1 ./source2 ../../dest",
				},
			},
			Expected: []ast.InstructionNode{&ast.CopyInstructionNode{
				Source:      []string{"./source1", "./source2"},
				Destination: "../../dest",
				KeepGitDir:  false,
				Chown:       "",
				Link:        false,
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.ENTRYPOINT,
					Content: "cp ./source1 ./source2 ../../dest",
				},
			},
			Expected: []ast.InstructionNode{&ast.EntrypointInstructionNode{
				Exec: []string{"cp", "./source1", "./source2", "../../dest"},
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.ENTRYPOINT,
					Content: "[\"cp\", \"./source1\", \"./source2\", \"../../dest\"]",
				},
			},
			Expected: []ast.InstructionNode{&ast.EntrypointInstructionNode{
				Exec: []string{"cp", "./source1", "./source2", "../../dest"},
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
