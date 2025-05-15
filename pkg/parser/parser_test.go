package parser_test

import (
	"fmt"
	"reflect"
	"testing"

	testdata "github.com/coffeemakingtoaster/dockerfile-parser/internal/pkg/test_data"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/ast"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/lexer"
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

	if len(actual.Subsequent) != len(expected.Subsequent) {
		return fmt.Sprintf("Stage subsequent length mismatch: Expected %d Got %d", len(expected.Subsequent), len(actual.Subsequent))
	}

	for i := range expected.Subsequent {
		if err := compareStageNodes(*expected.Subsequent[i], *actual.Subsequent[i]); err != "" {
			return err
		}
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
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return fmt.Sprintf("Type mismatch: Expected %v Got %v", reflect.TypeOf(expected), reflect.TypeOf(actual))
	}
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
	case *ast.EnvInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.EnvInstructionNode).Pairs, ac.Pairs) {
			return fmt.Sprintf("ENV instruction command mismatch: Expected %v Got %v", expected.(*ast.EnvInstructionNode).Pairs, ac.Pairs)
		}
	case *ast.ExposeInstructionNode:
		if *expected.(*ast.ExposeInstructionNode) != *ac {
			return fmt.Sprintf("EXPOSE instruction mismatch: Expected %v Got %v", expected, ac)
		}
	case *ast.HealthcheckInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.HealthcheckInstructionNode), ac) {
			return fmt.Sprintf("HEALTHCHECK instruction mismatch: Expected %v Got %v", expected, ac)
		}
	case *ast.LabelInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.LabelInstructionNode), ac) {
			return fmt.Sprintf("LABEL instruction mismatch: Expected %v Got %v", expected, ac)
		}
	case *ast.MaintainerInstructionNode:
		if *expected.(*ast.MaintainerInstructionNode) != *ac {
			return fmt.Sprintf("MAINTAINER instruction mismatch: Expected %v Got %v", expected, ac)
		}
	case *ast.RunInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.RunInstructionNode), ac) {
			return fmt.Sprintf("RUN instruction mismatch: Expected %v Got %v", expected, ac)
		}
	case *ast.ShellInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.ShellInstructionNode), ac) {
			return fmt.Sprintf("SHELL instruction mismatch: Expected %v Got %v", expected, ac)
		}
	case *ast.OnbuildInstructionNode:
		if err := compareInstructionNode(expected.(*ast.OnbuildInstructionNode).Trigger, ac.Trigger); err != "" {
			//return fmt.Sprintf("%v %v", reflect.TypeOf(expected.(*ast.OnbuildInstructionNode).Trigger), reflect.TypeOf(ac.Trigger))
			return fmt.Sprintf("ONBUILD instruction instruction mismatch: %s", err)
		}
	case *ast.StopsignalInstructionNode:
		if expected.(*ast.StopsignalInstructionNode).Signal != ac.Signal {
			return fmt.Sprintf("STOPSIGNAL instruction signal mismatch: Expected %s Got %s", expected.(*ast.StopsignalInstructionNode).Signal, ac.Signal)
		}
	case *ast.UserInstructionNode:
		if expected.(*ast.UserInstructionNode).User != ac.User {
			return fmt.Sprintf("USER instruction user mismatch: Expected %s Got %s", expected.(*ast.UserInstructionNode).User, ac.User)
		}
	case *ast.VolumeInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.VolumeInstructionNode), ac) {
			return fmt.Sprintf("VOLUME instruction mismatch: Expected %v Got %v", expected, ac)
		}
	case *ast.WorkdirInstructionNode:
		if expected.(*ast.WorkdirInstructionNode).Path != ac.Path {
			return fmt.Sprintf("WORKDIR instruction path mismatch: Expected %s Got %s", expected.(*ast.WorkdirInstructionNode).Path, ac.Path)
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
		{
			Input: []token.Token{
				{
					Kind:    token.FROM,
					Content: "alpine:latest AS base",
				},
				{
					Kind:    token.FROM,
					Content: "alpine:padding AS padding",
				},
				{
					Kind:    token.FROM,
					Content: "alpine:next AS next",
				},
				{
					Kind:    token.COPY,
					Params:  map[string]string{"from": "base"},
					Content: "./source1 ./source2 ../../dest",
				},
				{
					Kind:    token.COPY,
					Params:  map[string]string{"from": "base"},
					Content: "./source1 ./source2 ../../dest",
				},
			},
			Expected: &ast.StageNode{
				Identifier: "base",
				Image:      "alpine:latest",
				Subsequent: []*ast.StageNode{
					{
						Identifier: "padding",
						Image:      "alpine:padding",
						Subsequent: []*ast.StageNode{
							{
								Identifier: "next",
								Image:      "alpine:next",
								Subsequent: []*ast.StageNode{},
							},
						},
					},
					{
						Identifier: "next",
						Image:      "alpine:next",
						Subsequent: []*ast.StageNode{},
					},
				},
			},
		},
	}

	for _, c := range testCases {
		p := parser.NewParser(c.Input)
		// Pass first in because there is no need to compare the rootnode
		err := compareStageNodes(*c.Expected, *p.Parse().Subsequent[0])
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
		{
			Input: []token.Token{
				{
					Kind:    token.ENV,
					Content: "ABC=test sample=\"val1 val2\"",
				},
			},
			Expected: []ast.InstructionNode{&ast.EnvInstructionNode{
				Pairs: map[string]string{
					"ABC":    "test",
					"sample": "\"val1 val2\"",
				},
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.EXPOSE,
					Content: "3100/udp",
				},
			},
			Expected: []ast.InstructionNode{&ast.ExposeInstructionNode{
				Port:  "3100",
				IsTCP: false,
			},
			},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.EXPOSE,
					Content: "5000/tcp",
				},
			},
			Expected: []ast.InstructionNode{&ast.ExposeInstructionNode{
				Port:  "5000",
				IsTCP: true,
			},
			},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.HEALTHCHECK,
					Content: "cp test1 test2",
				},
			},
			Expected: []ast.InstructionNode{&ast.HealthcheckInstructionNode{
				CancelStatement: false,
				Interval:        "30s",
				Timeout:         "30s",
				StartPeriod:     "0s",
				StartInterval:   "5s",
				Retries:         3,
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.HEALTHCHECK,
					Content: "NONE",
				},
			},
			Expected: []ast.InstructionNode{&ast.HealthcheckInstructionNode{
				CancelStatement: true,
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.LABEL,
					Content: "A=B",
				},
			},
			Expected: []ast.InstructionNode{&ast.LabelInstructionNode{
				Pairs: map[string]string{"A": "B"},
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.MAINTAINER,
					Content: "Peter Lustig",
				},
			},
			Expected: []ast.InstructionNode{&ast.MaintainerInstructionNode{
				Name: "Peter Lustig",
			}},
		},
		/*
						This leads to weird problems with any ast node that contains variable sized fields (strings, arrays,...)
			TODO: Investigate further

									{
										Input: []token.Token{
											{
												Kind:    token.ONBUILD,
												Content: "HEALTHCHECK NONE",
											},
										},
										Expected: []ast.InstructionNode{&ast.OnbuildInstructionNode{
											// Every struct placed here that contains a variable sized type seems to not be read properly
											Trigger: &ast.HealthcheckInstructionNode{
												CancelStatement: true,
											},
										}},
									},
		*/
		{
			Input: []token.Token{
				{
					Kind:    token.RUN,
					Content: "cp ./a ./b",
				},
			},
			Expected: []ast.InstructionNode{&ast.RunInstructionNode{
				Cmd:       []string{"cp", "./a", "./b"},
				ShellForm: false,
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.SHELL,
					Content: "/bin/bash -c",
				},
			},
			Expected: []ast.InstructionNode{&ast.ShellInstructionNode{
				Shell: []string{"/bin/bash", "-c"},
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.STOPSIGNAL,
					Content: "SIGKILL",
				},
			},
			Expected: []ast.InstructionNode{&ast.StopsignalInstructionNode{Signal: "SIGKILL"}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.USER,
					Content: "root",
				},
			},
			Expected: []ast.InstructionNode{&ast.UserInstructionNode{User: "root"}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.VOLUME,
					Content: "./base ./base2",
				},
			},
			Expected: []ast.InstructionNode{&ast.VolumeInstructionNode{Mounts: []string{"./base", "./base2"}}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.VOLUME,
					Content: "[\"./base\", \"./base2\"",
				},
			},
			Expected: []ast.InstructionNode{&ast.VolumeInstructionNode{Mounts: []string{"./base", "./base2"}}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.WORKDIR,
					Content: "/app",
				},
			},
			Expected: []ast.InstructionNode{&ast.WorkdirInstructionNode{Path: "/app"}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.CMD,
					Content: "",
				},
			},
			Expected: []ast.InstructionNode{&ast.CmdInstructionNode{Cmd: []string{}}},
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

func BenchmarkParser(b *testing.B) {
	// Create these once
	l := lexer.NewFromInput(testdata.SampleDockerfile)
	tokens, err := l.Lex()
	if err != nil {
		b.Fatalf("Lexing failed: %s", err.Error())
	}
	for range b.N {
		p := parser.NewParser(tokens)
		p.Parse()
	}
}
