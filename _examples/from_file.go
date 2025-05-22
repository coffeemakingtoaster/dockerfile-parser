package main

import (
	"fmt"
	"reflect"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/lexer"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/parser"
)

func main() {
	l, err := lexer.NewFromFile("./example.Dockerfile")
	if err != nil {
		panic(err)
	}
	tokens, err := l.Lex()
	if err != nil {
		panic(err)
	}
	p := parser.NewParser(tokens)
	rootNode := p.Parse()
	for rootNode != nil {
		fmt.Printf("Stage:  %s\n", rootNode.Identifier)
		if len(rootNode.Instructions) > 0 {
			for _, instr := range rootNode.Instructions {
				fmt.Printf("InstructionNode: %s\n", reflect.TypeOf(instr))
			}
		}
		if len(rootNode.Subsequent) == 0 {
			rootNode = nil
			continue
		}
		rootNode = rootNode.Subsequent[0]
	}
}
