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
	return Parser{tokens: tokens, currentTokenIndex: 0, rootNode: &ast.StageNode{Identifier: ast.GenerateStageNodeID(), ParserMetadata: make(map[string]string)}}
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
			localRoot.Subsequent = node
			if node.Name != "" {
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
					if !slices.Contains(val.ReferencedByIds, localRoot.Identifier) {
						val.ReferencedByIds = append(val.ReferencedByIds, localRoot.Identifier)
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
		case token.EMPTY_LINE:
			node := &ast.EmptyLineNode{}
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
			Identifier:     ast.GenerateStageNodeID(),
			Image:          t.Content,
			ParserMetadata: make(map[string]string),
		}
	}
	content := parsePossibleArray(t.Content)
	image := content[0]
	// as := content[1]
	name := strings.Join(content[2:], " ")
	return &ast.StageNode{
		Identifier:     ast.GenerateStageNodeID(),
		Name:           name,
		Image:          image,
		ParserMetadata: make(map[string]string),
	}
}

func (p Parser) parseAdd(t token.Token) ast.InstructionNode {
	paths := parsePossibleArray(t.Content)
	cleanedPaths := CleanSlice(paths)
	return &ast.AddInstructionNode{
		Source:      cleanedPaths[0 : len(cleanedPaths)-1],
		Destination: cleanedPaths[len(cleanedPaths)-1],
		KeepGitDir:  util.GetFromParamsWithDefault(t.Params, "keep-git-dir", []string{"false"})[0] == "true",
		CheckSum:    util.GetFromParamsWithDefault(t.Params, "checksum", []string{""})[0],
		Chown:       util.GetFromParamsWithDefault(t.Params, "chown", []string{""})[0],
		Chmod:       util.GetFromParamsWithDefault(t.Params, "chmod", []string{""})[0],
		Link:        util.GetFromParamsWithDefault(t.Params, "link", []string{"false"})[0] == "true",
		Exclude:     util.GetFromParamsWithDefault(t.Params, "exclude", []string{""})[0],
	}
}

func (p Parser) parseArg(t token.Token) ast.InstructionNode {
	pairs := util.ParseAssigns(t.Content)
	if len(pairs) == 0 {
		pairs[t.Content] = ""
	}
	return &ast.ArgInstructionNode{
		Pairs: pairs,
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
	cleanedPaths := CleanSlice(paths)

	return &ast.CopyInstructionNode{
		Source:      cleanedPaths[0 : len(cleanedPaths)-1],
		Destination: cleanedPaths[len(cleanedPaths)-1],
		KeepGitDir:  util.GetFromParamsWithDefault(t.Params, "keep-git-dir", []string{"false"})[0] == "true",
		Chown:       util.GetFromParamsWithDefault(t.Params, "chown", []string{""})[0],
		Link:        util.GetFromParamsWithDefault(t.Params, "link", []string{"false"})[0] == "true",
		From:        util.GetFromParamsWithDefault(t.Params, "from", []string{""})[0],
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
	retries, _ := strconv.Atoi(util.GetFromParamsWithDefault(t.Params, "retries", []string{"3"})[0])
	return &ast.HealthcheckInstructionNode{
		CancelStatement: false,
		Interval:        util.GetFromParamsWithDefault(t.Params, "interval", []string{"30s"})[0],
		Timeout:         util.GetFromParamsWithDefault(t.Params, "timeout", []string{"30s"})[0],
		StartPeriod:     util.GetFromParamsWithDefault(t.Params, "start-period", []string{"0s"})[0],
		StartInterval:   util.GetFromParamsWithDefault(t.Params, "start-interval", []string{"5s"})[0],
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
	var content []string
	// We do nothing with heredoc except add it directly
	if len(t.MultiLineContent) != 0 {
		content = t.MultiLineContent
	} else {
		content = parsePossibleArray(t.Content)
	}
	return &ast.RunInstructionNode{
		Cmd:       content,
		ShellForm: false,
		IsHeredoc: len(t.MultiLineContent) > 0,
		Device:    util.GetFromParamsWithDefault(t.Params, "device", []string{""})[0],
		Security:  util.GetFromParamsWithDefault(t.Params, "security", []string{""})[0], // technically the default here is sandbox...but currently this parameter only exists in labs
		Network:   util.GetFromParamsWithDefault(t.Params, "network", []string{""})[0],
		Mount:     util.GetFromParamsWithDefault(t.Params, "mount", []string{}),
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
	return &ast.WorkdirInstructionNode{Path: strings.TrimSpace(t.Content)}
}
