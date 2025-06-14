package wrapper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coffeemakingtoaster/dockerfile-parser/internal/pkg/display"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/ast"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/lexer"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/parser"
)

func ParsePath(path string, recursive, output bool) int {
	isFile, err := isFile(path)
	if err != nil {
		panic(err)
	}
	if isFile {
		parseAndDisplayFileList([]string{path}, output)
		return 1
	} else {
		paths := buildDirPathList(path, recursive)
		parseAndDisplayFileList(paths, output)
		return len(paths)
	}
}

// Only looks for files ending with .Dockerfile
func buildDirPathList(basePath string, recursive bool) []string {
	res := []string{}
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return res
	}
	for _, entry := range entries {
		fullPath := filepath.Join(basePath, entry.Name())
		if entry.IsDir() {
			if recursive {
				subFiles := buildDirPathList(fullPath, recursive)
				res = append(res, subFiles...)
			}
		} else if strings.HasSuffix(entry.Name(), ".Dockerfile") {
			res = append(res, fullPath)
		}
	}
	return res
}

func isFile(path string) (bool, error) {
	// This returns an *os.FileInfo type
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	// IsDir is short for fileInfo.Mode().IsDir()
	if fileInfo.IsDir() {
		return false, nil
	} else {
		return true, nil
	}
}

func parseAndDisplayFileList(paths []string, output bool) {
	for _, path := range paths {
		fmt.Printf("---\t%s\t---\n", path)
		l, err := lexer.NewFromFile(path)
		if err != nil {
			panic(err)
		}
		tokens, err := l.Lex()
		if err != nil {
			panic(err)
		}
		p := parser.NewParser(tokens)
		root := p.Parse()

		if root == nil {
			fmt.Printf("Dockerfile at path %s contains no valid instruction of no FROM", path)
			continue
		}
		display.DisplayAst(root)
		if output {
			outputReconstructed(root, filepath.Base(path))
		}
	}
}

func outputReconstructed(root *ast.StageNode, filename string) {
	os.MkdirAll("./out", 0755)
	content := root.Reconstruct()
	data := strings.Join(content, "\n")
	err := os.WriteFile(filepath.Join("./out", filename), []byte(data), 0755)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err.Error())
	}
}
