package parser

import (
	"fmt"
	"strings"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/ast"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/util"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
)

type Parser struct {
	tokens            []token.Token
	currentTokenIndex int
	rootNode          *ast.StageNode
}

func NewParser(tokens []token.Token) Parser {
	return Parser{tokens: tokens, currentTokenIndex: 0}
}

func (p *Parser) Parse() ast.StageNode {
	localRoot := p.rootNode
	for {
		if p.currentTokenIndex == len(p.tokens) {
			break
		}
		t := p.tokens[p.currentTokenIndex]
		if p.rootNode == nil {
			if p.tokens[p.currentTokenIndex].Kind != token.FROM {
				panic(fmt.Sprintf("Could not create ast. Expected first node to be FROM but was %d", t.Kind))
			}
			p.rootNode = p.parseFrom(t)
			p.currentTokenIndex += 1
			localRoot = p.rootNode
			continue
		}
		switch t.Kind {
		case token.FROM:
			node := p.parseFrom(t)
			// TODO: this should detect if stages reference it
			localRoot.Subsequent = append(localRoot.Subsequent, node)
			localRoot = node
		case token.ADD:
			node := p.parseAdd(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.ARG:
			node := p.parseArg(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.CMD:
			node := p.parseCmd(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.COPY:
			node := p.parseCopy(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.ENTRYPOINT:
			node := p.parseEntryPoint(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		default:
			fmt.Printf("Not implemented kind %d", t.Kind)
		}
		p.currentTokenIndex += 1
	}
	// Return by copy
	return *p.rootNode
}

func (p Parser) parseFrom(t token.Token) *ast.StageNode {
	if !(strings.Contains(t.Content, " AS ") || strings.Contains(t.Content, " as ")) {
		return &ast.StageNode{
			Identifier: "anon",
			Image:      t.Content,
		}
	}
	content := parsePossibleArray(t.Content)
	image := content[0]
	// as := content[1]
	identifier := strings.Join(content[2:], " ")
	return &ast.StageNode{
		Identifier: identifier,
		Image:      image,
	}
}

func (p Parser) parseAdd(t token.Token) ast.InstructionNode {
	paths := parsePossibleArray(t.Content)
	return &ast.AddInstructionNode{
		Source:      paths[0 : len(paths)-1],
		Destination: paths[len(paths)-1],
		KeepGitDir:  util.GetFromParamsWithDefault(t.Params, "keep-git-dir", "false") == "true",
		CheckSum:    util.GetFromParamsWithDefault(t.Params, "checksum", ""),
		Chown:       util.GetFromParamsWithDefault(t.Params, "chown", ""),
		Chmod:       util.GetFromParamsWithDefault(t.Params, "chmod", ""),
		Link:        util.GetFromParamsWithDefault(t.Params, "link", "false") == "true",
		Exclude:     util.GetFromParamsWithDefault(t.Params, "exclude", ""),
	}
}

func (p Parser) parseArg(t token.Token) ast.InstructionNode {
	key, value := parseAssign(t.Content)
	return &ast.ArgInstructionNode{
		Name:  key,
		Value: value,
	}
}

func (p Parser) parseCmd(t token.Token) ast.InstructionNode {
	return &ast.CmdInstructionNode{
		Cmd: parsePossibleArray(t.Content),
	}
}

func (p Parser) parseCopy(t token.Token) ast.InstructionNode {
	paths := parsePossibleArray(t.Content)
	return &ast.CopyInstructionNode{
		Source:      paths[0 : len(paths)-1],
		Destination: paths[len(paths)-1],
		KeepGitDir:  util.GetFromParamsWithDefault(t.Params, "keep-git-dir", "false") == "true",
		Chown:       util.GetFromParamsWithDefault(t.Params, "chown", ""),
		Link:        util.GetFromParamsWithDefault(t.Params, "link", "false") == "true",
	}
}

func (p Parser) parseEntryPoint(t token.Token) ast.InstructionNode {
	return &ast.EntrypointInstructionNode{
		Exec: parsePossibleArray(t.Content),
	}

}
