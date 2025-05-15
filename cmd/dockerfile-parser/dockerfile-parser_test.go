package main_test

import (
	"testing"

	testdata "github.com/coffeemakingtoaster/dockerfile-parser/internal/pkg/test_data"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/lexer"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/parser"
)

func BenchmarkTotal(b *testing.B) {
	for range b.N {
		l := lexer.NewFromInput(testdata.SampleDockerfile)
		tokens, err := l.Lex()
		if err != nil {
			b.Fatalf("Lexing failed: %s", err.Error())
		}
		p := parser.NewParser(tokens)
		p.Parse()
	}
}
