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

	if actual.Name != expected.Name {
		return fmt.Sprintf("Stage name mismatch: Expected %s Got %s", expected.Name, actual.Name)
	}

	if actual.Identifier != fmt.Sprintf("%s-identifier", actual.Name) {
		return fmt.Sprintf("Stage identifier mismatch: Expected %s (%s) Got %s (%s)", fmt.Sprintf("%s-identifier", actual.Name), expected.Name, actual.Identifier, expected.Name)
	}

	if expected.Subsequent == nil || actual.Subsequent == nil {
		if expected.Subsequent != actual.Subsequent {
			return fmt.Sprintf("Subsequent mismatch: Expected %v Got %v", expected.Subsequent, actual.Subsequent)
		}

	} else {
		if err := compareStageNodes(*expected.Subsequent, *actual.Subsequent); err != "" {
			return fmt.Sprintf("Nested mismatch: %s", err)
		}
	}

	if !reflect.DeepEqual(actual.ReferencedByIds, expected.ReferencedByIds) {
		return fmt.Sprintf("Stage reference ids mismatch: Expected %v Got %v", expected.ReferencedByIds, actual.ReferencedByIds)
	}

	if !reflect.DeepEqual(expected.ParserMetadata, actual.ParserMetadata) {
		return fmt.Sprintf("Stage parser metadata mismatch: Expected %v Got %v", expected.ParserMetadata, actual.ParserMetadata)
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
	case *ast.CommentInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.CommentInstructionNode), ac) {
			return fmt.Sprintf("COMMENT content mismatch: Expected %v Got %v", expected.(*ast.CommentInstructionNode).Text, ac.Text)
		}
	case *ast.EntrypointInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.EntrypointInstructionNode).Exec, ac.Exec) {
			return fmt.Sprintf("ENTRYPOINT instruction command mismatch: Expected %v Got %v", expected.(*ast.EntrypointInstructionNode).Exec, ac.Exec)
		}
	case *ast.EnvInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.EnvInstructionNode).Pairs, ac.Pairs) {
			return fmt.Sprintf("ENV instruction command mismatch: Expected %v Got %v", expected.(*ast.EnvInstructionNode).Pairs, ac.Pairs)
		}
	case *ast.ExposeInstructionNode:
		if !reflect.DeepEqual(expected.(*ast.ExposeInstructionNode), ac) {
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
				Name:           "",
				Identifier:     "-identifier",
				Image:          "alpine:latest",
				ParserMetadata: make(map[string]string),
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
				Name:           "base",
				Identifier:     "base-identifier",
				Image:          "alpine:latest",
				ParserMetadata: make(map[string]string),
			},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.FROM,
					Content: "alpine:latest AS base",
				},
				{
					Kind:    token.PARSER_DIRECTIVE,
					Content: "syntax=docker/dockerfile:1",
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
					Params:  map[string][]string{"from": {"base"}},
					Content: "./source1 ./source2 ../../dest",
				},
				{
					Kind:    token.COPY,
					Params:  map[string][]string{"from": {"base"}},
					Content: "./source1 ./source2 ../../dest",
				},
			},
			Expected: &ast.StageNode{
				Name:           "base",
				Image:          "alpine:latest",
				Identifier:     "base-identifier",
				ParserMetadata: map[string]string{"syntax": "docker/dockerfile:1"},
				Subsequent: &ast.StageNode{
					Name:       "padding",
					Image:      "alpine:padding",
					Identifier: "padding-identifier",

					Subsequent: &ast.StageNode{
						Name:       "next",
						Identifier: "next-identifier",

						Image:          "alpine:next",
						Subsequent:     nil,
						ParserMetadata: make(map[string]string),
					},
					ParserMetadata: make(map[string]string),
				},
			},
		},
	}
	for _, c := range testCases {
		p := parser.NewParser(c.Input)
		actual := p.Parse()
		curr := actual.Subsequent
		// overwrite generated ids with predictable ids
		for curr != nil {
			curr.Identifier = fmt.Sprintf("%s-identifier", curr.Name)
			curr = curr.Subsequent
		}
		// Pass first in because there is no need to compare the rootnode
		err := compareStageNodes(*c.Expected, *actual.Subsequent)
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
					Kind:    token.ARG,
					Content: "test",
				},
			},
			Expected: []ast.InstructionNode{&ast.ArgInstructionNode{
				Name:  "test",
				Value: "",
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
				Ports: []ast.PortInfo{{
					Port:  "3100",
					IsTCP: false,
				}},
			},
			},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.EXPOSE,
					Content: " 5000/tcp",
				},
			},
			Expected: []ast.InstructionNode{&ast.ExposeInstructionNode{
				Ports: []ast.PortInfo{{
					Port:  "5000",
					IsTCP: true,
				}},
			},
			},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.EXPOSE,
					Content: "5000/tcp 3000/udp",
				},
			},
			Expected: []ast.InstructionNode{&ast.ExposeInstructionNode{
				Ports: []ast.PortInfo{
					{
						Port:  "5000",
						IsTCP: true,
					},
					{
						Port:  "3000",
						IsTCP: false,
					},
				},
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
				Cmd:             []string{"cp", "test1", "test2"},
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
		{
			Input: []token.Token{
				{
					Kind:    token.ONBUILD,
					Content: "EXPOSE 5000/tcp 3000/udp",
				},
			},
			Expected: []ast.InstructionNode{&ast.OnbuildInstructionNode{
				Trigger: &ast.ExposeInstructionNode{
					Ports: []ast.PortInfo{
						{
							Port:  "5000",
							IsTCP: true,
						},
						{
							Port:  "3000",
							IsTCP: false,
						},
					},
				}}},
		},
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
				IsHeredoc: false,
				Mount:     []string{},
				Network:   "",
				Security:  "",
				Device:    "",
			}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.RUN,
					Content: "cp ./a ./b",
					Params: map[string][]string{
						"security": {"sandbox"},
						"device":   {"gpu"},
						"mount":    {"test1", "test2"},
						"network":  {"nono"},
					},
				},
			},
			Expected: []ast.InstructionNode{&ast.RunInstructionNode{
				Cmd:       []string{"cp", "./a", "./b"},
				ShellForm: false,
				IsHeredoc: false,
				Mount:     []string{"test1", "test2"},
				Network:   "nono",
				Security:  "sandbox",
				Device:    "gpu",
			}},
		},

		{
			Input: []token.Token{
				{
					Kind:               token.RUN,
					Content:            "",
					MultiLineContent:   []string{"EOT bash", "set -ex", "apt-get update", "apt-get install -y vim", "EOT"},
					HereDocRedirection: false,
				},
			},
			Expected: []ast.InstructionNode{&ast.RunInstructionNode{
				Cmd:       []string{"EOT bash", "set -ex", "apt-get update", "apt-get install -y vim", "EOT"},
				ShellForm: false,
				IsHeredoc: true,
				Mount:     []string{},
				Network:   "",
				Security:  "",
				Device:    "",
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
		{
			Input: []token.Token{
				{
					Kind:    token.COMMENT,
					Content: "some helpful comment",
				},
			},
			Expected: []ast.InstructionNode{&ast.CommentInstructionNode{Text: "some helpful comment"}},
		},
		{
			Input: []token.Token{
				{
					Kind:    token.PARSER_DIRECTIVE,
					Content: "syntax=docker/dockerfile:1",
				},
			},
			Expected: []ast.InstructionNode{},
		},
	}
	for _, c := range testCases {
		p := parser.NewParser(c.Input)
		instructions := p.Parse().Instructions
		if len(instructions) != len(c.Expected) {
			t.Errorf("Instruction count mismatch: Expected %d Got %d", len(c.Expected), len(instructions))
		}
		for i := range instructions {
			err := compareInstructionNode(c.Expected[i], instructions[i])
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
