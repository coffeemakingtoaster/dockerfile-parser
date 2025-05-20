// Parser package
package parser

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/ast"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/lexer"
	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/util"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/token"
)

// Parser
type Parser struct {
	tokens            []token.Token
	currentTokenIndex int
	rootNode          *ast.StageNode
}

// Create new parser
func NewParser(tokens []token.Token) Parser {
	return Parser{tokens: tokens, currentTokenIndex: 0, rootNode: &ast.StageNode{Identifier: "root", ParserMetadata: make(map[string]string)}}
}

// Parse the token provided during init
// Return the root stage node of the ast
func (p *Parser) Parse() *ast.StageNode {
	localRoot := p.rootNode
	namedStageLookup := make(map[string]*ast.StageNode)
	for {
		if p.currentTokenIndex == len(p.tokens) {
			break
		}
		t := p.tokens[p.currentTokenIndex]
		switch t.Kind {
		case token.FROM:
			node := p.parseFrom(t)
			localRoot.Subsequent = append(localRoot.Subsequent, node)
			if node.Identifier != "anon" {
				namedStageLookup[node.Identifier] = node
			}
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
			if node.(*ast.CopyInstructionNode).From != "" {
				// Check if from actually is a stage -> can also be image
				if val, ok := namedStageLookup[node.(*ast.CopyInstructionNode).From]; ok {
					if !slices.Contains(val.Subsequent, localRoot) {
						val.Subsequent = append(val.Subsequent, localRoot)
					}
				}
			}
		case token.ENTRYPOINT:
			node := p.parseEntryPoint(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.ENV:
			node := p.parseEnv(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.EXPOSE:
			node := p.parseExpose(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.HEALTHCHECK:
			node := p.parseHealthCheck(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.LABEL:
			node := p.parseLabel(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.MAINTAINER:
			node := p.parseMaintainer(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.ONBUILD:
			node := p.parseOnBuild(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.RUN:
			node := p.parseRun(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.SHELL:
			node := p.parseShell(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.STOPSIGNAL:
			node := p.parseStop(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.USER:
			node := p.parseUser(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.WORKDIR:
			node := p.parseWorkdir(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.VOLUME:
			node := p.parseVolume(t)
			localRoot.Instructions = append(localRoot.Instructions, node)
		case token.PARSER_DIRECTIVE:
			key, value := util.ParseAssign(t.Content)
			localRoot.ParserMetadata[key] = value
		case token.COMMENT:
			node := &ast.CommentInstructionNode{Text: t.Content}
			localRoot.Instructions = append(localRoot.Instructions, node)
		default:
			fmt.Printf("Not implemented kind %d", t.Kind)
		}
		p.currentTokenIndex += 1
	}
	return p.rootNode
}

func (p Parser) parseFrom(t token.Token) *ast.StageNode {
	if !(strings.Contains(t.Content, " AS ") || strings.Contains(t.Content, " as ")) {
		return &ast.StageNode{
			Identifier:     "anon",
			Image:          t.Content,
			ParserMetadata: make(map[string]string),
		}
	}
	content := parsePossibleArray(t.Content)
	image := content[0]
	// as := content[1]
	identifier := strings.Join(content[2:], " ")
	return &ast.StageNode{
		Identifier:     identifier,
		Image:          image,
		ParserMetadata: make(map[string]string),
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
	key, value := util.ParseAssign(t.Content)
	if len(key)+len(value) == 0 {
		key = t.Content
	}
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
	if len(t.MultiLineContent) != 0 {
		fmt.Print("Heredoc copy statements are not supported as of now")
		return &ast.CopyInstructionNode{}
	}
	paths := parsePossibleArray(t.Content)
	return &ast.CopyInstructionNode{
		Source:      paths[0 : len(paths)-1],
		Destination: paths[len(paths)-1],
		KeepGitDir:  util.GetFromParamsWithDefault(t.Params, "keep-git-dir", "false") == "true",
		Chown:       util.GetFromParamsWithDefault(t.Params, "chown", ""),
		Link:        util.GetFromParamsWithDefault(t.Params, "link", "false") == "true",
		From:        util.GetFromParamsWithDefault(t.Params, "from", ""),
	}
}

func (p Parser) parseEntryPoint(t token.Token) ast.InstructionNode {
	return &ast.EntrypointInstructionNode{
		Exec: parsePossibleArray(t.Content),
	}
}

func (p Parser) parseEnv(t token.Token) ast.InstructionNode {
	return &ast.EnvInstructionNode{
		Pairs: util.ParseAssigns(t.Content),
	}
}

func (p Parser) parseExpose(t token.Token) ast.InstructionNode {
	ports := []ast.PortInfo{}

	parts := strings.Split(t.Content, " ")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		isTcp := true
		v := strings.Split(part, "/")
		// protocol is present
		if len(v) > 1 {
			isTcp = v[1] == "tcp"
		}
		ports = append(ports, ast.PortInfo{Port: v[0], IsTCP: isTcp})
	}

	return &ast.ExposeInstructionNode{
		Ports: ports,
	}
}

func (p Parser) parseHealthCheck(t token.Token) ast.InstructionNode {
	if t.Content == "NONE" {
		return &ast.HealthcheckInstructionNode{CancelStatement: true}
	}
	retries, _ := strconv.Atoi(util.GetFromParamsWithDefault(t.Params, "retries", "3"))
	return &ast.HealthcheckInstructionNode{
		CancelStatement: false,
		Interval:        util.GetFromParamsWithDefault(t.Params, "interval", "30s"),
		Timeout:         util.GetFromParamsWithDefault(t.Params, "timeout", "30s"),
		StartPeriod:     util.GetFromParamsWithDefault(t.Params, "start-period", "0s"),
		StartInterval:   util.GetFromParamsWithDefault(t.Params, "start-interval", "5s"),
		Retries:         retries,
		Cmd:             parsePossibleArray(t.Content),
	}
}

func (p Parser) parseLabel(t token.Token) ast.InstructionNode {
	return &ast.LabelInstructionNode{
		Pairs: util.ParseAssigns(t.Content),
	}
}

func (p Parser) parseMaintainer(t token.Token) ast.InstructionNode {
	return &ast.MaintainerInstructionNode{
		Name: t.Content,
	}
}

func (p Parser) parseOnBuild(t token.Token) ast.InstructionNode {
	// Easiest way to do this is by simply running the instruction through the entire lexer -> parser process
	l := lexer.NewFromInput([]string{t.Content})
	tokens, err := l.Lex()
	if err != nil {
		return &ast.OnbuildInstructionNode{
			Trigger: &ast.UnknownInstructionNode{Text: t.Content},
		}
	}
	tmpP := NewParser(tokens)
	parsed := tmpP.Parse().Instructions[0]
	return &ast.OnbuildInstructionNode{
		Trigger: parsed,
	}
}

func (p Parser) parseRun(t token.Token) ast.InstructionNode {
	// We do nothing with heredoc except add it directly
	if len(t.MultiLineContent) != 0 {
		return &ast.RunInstructionNode{
			Cmd:       t.MultiLineContent,
			ShellForm: false,
			IsHeredoc: true,
		}
	}
	return &ast.RunInstructionNode{
		Cmd:       parsePossibleArray(t.Content),
		ShellForm: false,
		IsHeredoc: false,
	}
}

func (p Parser) parseShell(t token.Token) ast.InstructionNode {
	return &ast.ShellInstructionNode{
		Shell: parsePossibleArray(t.Content),
	}
}

func (p Parser) parseStop(t token.Token) ast.InstructionNode {
	return &ast.StopsignalInstructionNode{
		Signal: t.Content,
	}
}

func (p Parser) parseUser(t token.Token) ast.InstructionNode {
	return &ast.UserInstructionNode{User: t.Content}
}

func (p Parser) parseVolume(t token.Token) ast.InstructionNode {
	return &ast.VolumeInstructionNode{Mounts: parsePossibleArray(t.Content)}
}

func (p Parser) parseWorkdir(t token.Token) ast.InstructionNode {
	return &ast.WorkdirInstructionNode{Path: t.Content}
}
