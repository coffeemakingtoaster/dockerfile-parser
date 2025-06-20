package display

import (
	"fmt"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/ast"
)

func DisplayAst(root *ast.StageNode) {
	for root != nil {
		fmt.Println(root.ToString())
		for _, instruction := range root.Instructions {
			fmt.Println(fmt.Sprintf(" > %s", instruction.ToString()))
		}
		root = root.Subsequent
	}
}
