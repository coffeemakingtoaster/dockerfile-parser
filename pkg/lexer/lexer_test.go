package lexer_test

import (
	"fmt"
	"testing"

	testdata "github.com/coffeemakingtoaster/dockerfile-parser/internal/pkg/test_data"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/lexer"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
)

type TestCase struct {
	Input          []string
	ExpectedOutput []token.Token
}

func compareTokens(expected, actual token.Token) string {
	if expected.Kind != actual.Kind {
		return fmt.Sprintf("Token kind mismatch: Expected %d Got %d", expected.Kind, actual.Kind)
	}

	if len(expected.Params) == len(actual.Params) {
		for k, v := range expected.Params {
			if v != actual.Params[k] {
				return fmt.Sprintf("Token param value mismatch for key %s: Expected %s Got %s", k, v, actual.Params[k])
			}
		}
	} else {
		return fmt.Sprintf("Token param count mismatch: Expected %d Got %d", len(expected.Params), len(actual.Params))
	}

	if expected.Content != expected.Content {
		return fmt.Sprintf("Token content mismatch: Expected %s Got %s", expected.Content, actual.Content)
	}
	return ""
}

func TestInstructionParse(t *testing.T) {
	testCases := []TestCase{
		{
			Input: []string{"FROM test:latest"},
			ExpectedOutput: []token.Token{
				{
					Kind:    token.FROM,
					Params:  make(map[string]string),
					Content: "test:latest",
				},
			},
		},
		{
			Input: []string{"COPY --from=build /hello /"},
			ExpectedOutput: []token.Token{
				{
					Kind: token.COPY,
					Params: map[string]string{
						"from": "build",
					},
					Content: "test:latest",
				},
			},
		},
		{
			Input: []string{"RUN --mount=type=bind,target=. go build -o /myapp ./cmd"},
			ExpectedOutput: []token.Token{
				{
					Kind: token.RUN,
					Params: map[string]string{
						"mount": "type=bind,target=.",
					},
					Content: "go build -o /myapp ./cmd",
				},
			},
		}, {
			Input: []string{"COPY --link --from=build /foo /bar"},
			ExpectedOutput: []token.Token{
				{
					Kind: token.COPY,
					Params: map[string]string{
						"link": "true",
						"from": "build",
					},
					Content: "/foo /bar",
				},
			},
		},
		{
			Input: []string{"RUN --test=test command --commandflag=notmeantfordocker"},
			ExpectedOutput: []token.Token{
				{
					Kind: token.RUN,
					Params: map[string]string{
						"test": "test",
					},
					Content: "command --commandflag=notmeantfordocker",
				},
			},
		},
		{
			Input: []string{"CMD"},
			ExpectedOutput: []token.Token{
				{
					Kind:    token.CMD,
					Params:  map[string]string{},
					Content: "",
				},
			},
		},
		{
			Input: []string{"#test"},
			ExpectedOutput: []token.Token{
				{
					Kind:    token.COMMENT,
					Params:  map[string]string{},
					Content: "test",
				},
			},
		},
		{
			Input: []string{"# test"},
			ExpectedOutput: []token.Token{
				{
					Kind:    token.COMMENT,
					Params:  map[string]string{},
					Content: "test",
				},
			},
		},
		{
			Input: []string{"RUN echo a # test"},
			ExpectedOutput: []token.Token{
				{
					Kind:          token.RUN,
					Params:        map[string]string{},
					Content:       "echo a",
					InlineComment: "test",
				},
			},
		},
	}
	for _, v := range testCases {
		l := lexer.NewFromInput(v.Input)
		got, err := l.Lex()
		if err != nil {
			t.Fatalf("Failed to lex: %s", err.Error())
		}
		for i := range got {
			err := compareTokens(v.ExpectedOutput[i], got[i])
			if err != "" {
				t.Errorf("%s (%s)", err, v.Input)
			}
		}
	}
}

func BenchmarkLexer(b *testing.B) {
	for range b.N {
		// Lexer creation performs action on startup -> run this in the loop
		l := lexer.NewFromInput(testdata.SampleDockerfile)
		_, err := l.Lex()
		if err != nil {
			b.Fatalf("Lexing failed: %s", err.Error())
		}
	}
}
