package ast_test

import (
	"reflect"
	"testing"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/ast"
)

type Expected struct {
	Input    ast.StageNode
	Expected []string
}

func TestReconstructStage(t *testing.T) {
	expected := []Expected{
		{
			Input: ast.StageNode{
				Identifier: "internalandthereforeignored",
				Name:       "base",
				Image:      "img:latest",
			},
			Expected: []string{"FROM img:latest AS base"},
		},
		{
			Input: ast.StageNode{
				Identifier: "internalandthereforeignored",
				Name:       "",
				Image:      "debian:latest",
			},
			Expected: []string{"FROM debian:latest"},
		},
	}
	for _, testCase := range expected {
		actual := testCase.Input.Reconstruct()
		if !reflect.DeepEqual(actual, testCase.Expected) {
			t.Errorf("Reconstruct mismatch:\nExpected %+q\nGot %+q\n", testCase.Expected, actual)
		}
	}
}

func TestReconstructInstruction(t *testing.T) {
	expected := []Expected{
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.AddInstructionNode{
						Source:      []string{"./abc ./def"},
						Destination: "/home/new",
						KeepGitDir:  true,
						CheckSum:    "checksum",
					},
				},
			},
			Expected: []string{"ADD --keep-git-dir=true --checksum=checksum --link=false ./abc ./def /home/new"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.ArgInstructionNode{
						Pairs: map[string]string{
							"abc": "def",
						},
					},
				},
			},
			Expected: []string{"ARG abc=def"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.ArgInstructionNode{
						Pairs: map[string]string{
							"abc": "def",
							"xy":  "z",
						},
					},
				},
			},
			Expected: []string{"ARG abc=def xy=z"},
		},

		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.CmdInstructionNode{
						Cmd: []string{"curl", "ssh-coffee.dev", "&&", "whoami"},
					},
				},
			},
			Expected: []string{"CMD [\"curl\",\"ssh-coffee.dev\",\"&&\",\"whoami\"]"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.CopyInstructionNode{
						From:        "build",
						Source:      []string{"./abc ./def"},
						Destination: "/home/new",
						Link:        true,
					},
				},
			},
			Expected: []string{"COPY --keep-git-dir=false --link=true --from=build ./abc ./def /home/new"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.EntrypointInstructionNode{
						Exec: []string{"curl", "ssh-coffee.dev"},
					},
				},
			},
			Expected: []string{"ENTRYPOINT [\"curl\",\"ssh-coffee.dev\"]"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.EnvInstructionNode{
						Pairs: map[string]string{
							"xy":  "z",
							"abc": "def",
						},
					},
				},
			},
			Expected: []string{"ENV abc=def xy=z"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.ExposeInstructionNode{
						Ports: []ast.PortInfo{
							{
								Port:  "8080",
								IsTCP: false,
							},
							{
								Port:  "3000",
								IsTCP: true,
							},
						},
					},
				},
			},
			Expected: []string{"EXPOSE 8080/udp 3000/tcp"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.HealthcheckInstructionNode{
						Interval:      "31s",
						Timeout:       "32s",
						StartPeriod:   "33s",
						StartInterval: "34s",
						Retries:       3,
						Cmd:           []string{"curl", "localhost:8080/health"},
					},
				},
			},
			Expected: []string{"HEALTHCHECK --interval=31s --timeout=32s --start-period=33s --start-interval=34s --retries=3 [\"curl\",\"localhost:8080/health\"]"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.HealthcheckInstructionNode{
						CancelStatement: true,
					},
				},
			},
			Expected: []string{"HEALTHCHECK NONE"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.LabelInstructionNode{
						Pairs: map[string]string{
							"abc": "def",
							"xy":  "z",
						},
					},
				},
			},
			Expected: []string{"LABEL abc=def xy=z"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.MaintainerInstructionNode{
						Name: "Peter Lustig",
					},
				},
			},
			Expected: []string{"MAINTAINER Peter Lustig"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.OnbuildInstructionNode{
						Trigger: &ast.MaintainerInstructionNode{
							Name: "R2D2",
						},
					},
				},
			},
			Expected: []string{"ONBUILD MAINTAINER R2D2"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.OnbuildInstructionNode{
						Trigger: &ast.RunInstructionNode{
							Cmd:       []string{"EOF", "apt install curl", "curl ssh-coffee.dev", "EOF"},
							IsHeredoc: true,
						},
					},
				},
			},
			Expected: []string{"ONBUILD RUN << EOF", "apt install curl", "curl ssh-coffee.dev", "EOF"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.RunInstructionNode{
						Cmd:       []string{"curl google.com"},
						ShellForm: true,
					},
				},
			},
			Expected: []string{"RUN curl google.com"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.RunInstructionNode{
						Cmd:       []string{"curl", "google.com"},
						ShellForm: false,
					},
				},
			},
			Expected: []string{"RUN [\"curl\",\"google.com\"]"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.RunInstructionNode{
						Cmd:       []string{"EOF", "apt install curl", "curl ssh-coffee.dev", "EOF"},
						IsHeredoc: true,
					},
				},
			},
			Expected: []string{"RUN << EOF", "apt install curl", "curl ssh-coffee.dev", "EOF"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.ShellInstructionNode{
						Shell: []string{"curl", "google.com"},
					},
				},
			},
			Expected: []string{"SHELL [\"curl\",\"google.com\"]"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.StopsignalInstructionNode{
						Signal: "SIGKILL",
					},
				},
			},
			Expected: []string{"STOPSIGNAL SIGKILL"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.UserInstructionNode{
						User: "Jeff",
					},
				},
			},
			Expected: []string{"USER Jeff"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.VolumeInstructionNode{
						Mounts: []string{"/a", "/b"},
					},
				},
			},
			Expected: []string{"VOLUME [\"/a\",\"/b\"]"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.WorkdirInstructionNode{
						Path: "/a/b/c",
					},
				},
			},
			Expected: []string{"WORKDIR /a/b/c"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.CommentInstructionNode{
						Text: " Hello",
					},
				},
			},
			Expected: []string{"# Hello"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.UnknownInstructionNode{
						Text: "CUSTOM --flag command",
					},
				},
			},
			Expected: []string{"CUSTOM --flag command"},
		},
		{
			Input: ast.StageNode{
				Instructions: []ast.InstructionNode{
					&ast.OnbuildInstructionNode{
						Trigger: &ast.MaintainerInstructionNode{
							Name: "R2D2",
						},
					},
					&ast.EmptyLineNode{},
					&ast.OnbuildInstructionNode{
						Trigger: &ast.MaintainerInstructionNode{
							Name: "R2D2",
						},
					},
				},
			},
			Expected: []string{"ONBUILD MAINTAINER R2D2", "", "ONBUILD MAINTAINER R2D2"},
		},
	}
	for _, testCase := range expected {
		actual := []string{}
		for _, instruction := range testCase.Input.Instructions {
			actual = append(actual, instruction.Reconstruct()...)
		}
		if !reflect.DeepEqual(actual, testCase.Expected) {
			t.Errorf("Reconstruct mismatch:\nExpected %+q\nGot %+q\n", testCase.Expected, actual)
		}
	}
}
