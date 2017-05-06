package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spartanlogs/spartan/config/lexer"
	"github.com/spartanlogs/spartan/config/token"
	"github.com/spartanlogs/spartan/utils"
)

// A ParsedFile is the representation of a processed filters
// file.
type ParsedFile struct {
	// Inputs are definitions of the inputs defined in a file.
	// Input module names and options are not checked for correctness.
	Inputs []*InputDef

	// Filters are definitions of filters in the order written in
	// the file. Index 0 is the root filter in a pipeline.
	Filters []*PipelineDef

	// Outputs are definitions of outputs in the order written in
	// the file. Index 0 is the root filter in a pipeline.
	Outputs []*PipelineDef
}

// InputDef defines the module name and options map of an input.
type InputDef struct {
	Module  string
	Options utils.InterfaceMap
}

// PipelineDef defines a pipeline object (Filter/Output). It contains
// the module name, options map, and connections to the rest of the pipeline.
type PipelineDef struct {
	Module  string
	Options utils.InterfaceMap

	// Connections is a slice of index numbers corresponding to the index of
	// a pipeline definition in the parent ParsedFile struct.
	//
	// Length 0 -> end of pipeline
	// Length 1 -> normal next connection
	// Length 3 -> if statement connections
	//	Connections[0] -> true (inside if body)
	//	Connections[1] -> else (inside else body if present), nil if no else
	//	Connections[2] -> next object (after closing brace of if/else body)
	Connections []int
}

// ParseGlob will parse a collection of files in lexicographical order based
// on the supplied pattern. The files will effectivly be concatinated together
// as a single configuration.
func ParseGlob(pattern string) (*ParsedFile, error) {
	// This function works by replacing the state of the parser
	// after reading each file. The parser keeps its *ParsedFile
	// object after each parse so it simple expanded.

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return &ParsedFile{}, nil
	}

	file, err := os.Open(files[0])
	if err != nil {
		return nil, err
	}

	// First time using the parser is same as usual
	p := newParser(lexer.New(file))
	_, err = p.parse()
	if err != nil {
		file.Close()
		return nil, err
	}
	file.Close()

	// If other files need parsing, go over each file resetting the
	// state of the parse before each subsequent parse.
	for _, path := range files[1:] {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		// Reset parser state as if it was just created
		p.lexer = lexer.New(file)
		p.nextToken()
		p.nextToken()

		_, err = p.parse()
		if err != nil {
			file.Close()
			return nil, err
		}
		file.Close()
	}

	return p.file, nil
}

// ParseFile parses a single file without any globing.
func ParseFile(path string) (*ParsedFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Parse(file)
}

// ParseString parses the given string as a pipeline configuration.
func ParseString(s string) (*ParsedFile, error) {
	return Parse(strings.NewReader(s))
}

// Parse is a generic form of the other Parse* functions which reads a stream of
// tokens from a simple io.Reader.
func Parse(r io.Reader) (*ParsedFile, error) {
	return newParser(lexer.New(r)).parse()
}

type parser struct {
	file *ParsedFile

	curTok  token.Token
	peekTok token.Token

	lexer *lexer.Lexer
}

func newParser(l *lexer.Lexer) *parser {
	p := &parser{
		lexer: l,
		file:  &ParsedFile{},
	}
	p.nextToken() // Populate peekTok
	p.nextToken() // Populate curTok
	return p
}

func (p *parser) nextToken() token.Token {
	p.curTok = p.peekTok
	p.peekTok = p.lexer.NextToken()
	if p.curTok.Type == token.COMMENT {
		return p.nextToken()
	}
	return p.curTok
}

func (p *parser) parse() (*ParsedFile, error) {
parseLoop:
	for {
		var err error

		switch p.curTok.Type {
		case token.INPUT:
			p.nextToken()
			p.file.Inputs, err = p.parseInputs()
		case token.FILTER:
			p.nextToken()
			p.file.Filters, err = p.parsePipelineDefs()
		case token.OUTPUT:
			p.nextToken()
			p.file.Outputs, err = p.parsePipelineDefs()
		case token.EOF:
			break parseLoop
		default:
			err = p.tokenError(token.INPUT, token.FILTER, token.OUTPUT)
		}

		if err != nil {
			return nil, err
		}
	}

	return p.file, nil
}

func (p *parser) parseInputs() ([]*InputDef, error) {
	if p.curTok.Type != token.LBRACE {
		return nil, p.tokenError(token.LBRACE)
	}

	p.nextToken()

	inputs := make([]*InputDef, 0, 5)
	for {
		if p.curTok.Type == token.RBRACE {
			break
		}

		if p.curTok.Type != token.IDENT {
			return nil, p.tokenError(token.IDENT)
		}

		modName := p.curTok.Literal
		p.nextToken()

		options, err := p.parseMap()
		if err != nil {
			return nil, err
		}

		inputs = append(inputs, &InputDef{
			Module:  modName,
			Options: options,
		})
	}

	p.nextToken() // Consume closing }
	return inputs, nil
}

func (p *parser) parsePipelineDefs() ([]*PipelineDef, error) {
	if p.curTok.Type != token.LBRACE {
		return nil, p.tokenError(token.LBRACE)
	}

	p.nextToken()

	modules := make([]*PipelineDef, 0, 5)
	for {
		if p.curTok.Type == token.RBRACE {
			break
		}

		if p.curTok.Type != token.IDENT {
			return nil, p.tokenError(token.IDENT)
		}

		modName := p.curTok.Literal
		p.nextToken()

		options, err := p.parseMap()
		if err != nil {
			return nil, err
		}

		modules = append(modules, &PipelineDef{
			Module:  modName,
			Options: options,
		})
	}

	lastIndex := len(modules) - 1
	for i, mod := range modules {
		if i == lastIndex {
			break
		}
		mod.Connections = []int{i + 1}
	}

	p.nextToken() // Consume closing }
	return modules, nil
}

func (p *parser) parseMap() (utils.InterfaceMap, error) {
	if p.curTok.Type != token.LBRACE {
		return nil, p.tokenError(token.LBRACE)
	}

	p.nextToken()

	m := utils.NewInterfaceMap()
mapLoop:
	for {
		switch p.curTok.Type {
		case token.RBRACE:
			break mapLoop
		case token.COMMA: // Commas are optional
			p.nextToken()
			continue mapLoop
		case token.IDENT:
			fallthrough
		case token.STRING:
			key := p.curTok.Literal
			p.nextToken()

			if p.curTok.Type != token.ASSIGN {
				return nil, p.tokenError(token.ASSIGN)
			}

			p.nextToken()
			val := p.curTok

			switch val.Type {
			case token.TRUE:
				m.Set(key, true)
			case token.FALSE:
				m.Set(key, false)
			case token.STRING:
				m.Set(key, val.Literal)
			case token.INT:
				valInt, err := strconv.Atoi(val.Literal)
				if err != nil {
					return nil, fmt.Errorf("invalid integer %s", val.Literal)
				}
				m.Set(key, valInt)
			case token.FLOAT:
				valFloat, err := strconv.ParseFloat(val.Literal, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid floating point number %s", val.Literal)
				}
				m.Set(key, valFloat)
			case token.LSQUARE:
				array, err := p.parseArray()
				if err != nil {
					return nil, err
				}
				m.Set(key, array)
				continue mapLoop
			case token.LBRACE:
				subMap, err := p.parseMap()
				if err != nil {
					return nil, err
				}
				m.Set(key, subMap)
				continue mapLoop
			default:
				return nil, p.tokenError(
					token.STRING,
					token.INT,
					token.FLOAT,
					token.TRUE,
					token.FALSE,
					token.LBRACE,
					token.LSQUARE,
				)
			}
		default:
			return nil, fmt.Errorf("map key must be a string: %s", p.curTok.Type)
		}

		p.nextToken()
	}

	p.nextToken() // Consume closing }
	return m, nil
}

func (p *parser) parseArray() (interface{}, error) {
	if p.curTok.Type != token.LSQUARE {
		return nil, p.tokenError(token.LSQUARE)
	}

	p.nextToken()

	var array interface{}
	var err error
	switch p.curTok.Type {
	case token.STRING:
		array, err = p.parseStringArray()
	case token.INT:
		array, err = p.parseIntArray()
	case token.RSQUARE:
		break
	default:
		return nil, p.tokenError(token.STRING, token.INT, token.RSQUARE)
	}

	p.nextToken() // Consume closing ]

	if array == nil {
		array = make([]interface{}, 0)
	}
	return array, err
}

func (p *parser) parseStringArray() ([]string, error) {
	array := make([]string, 0, 5) // Make room for a conservitive number of elements

	for {
		array = append(array, p.curTok.Literal)
		p.nextToken()

		if p.curTok.Type == token.COMMA { // Commas are optional
			p.nextToken()
		}

		if p.curTok.Type == token.RSQUARE {
			break
		}

		if p.curTok.Type != token.STRING {
			return nil, p.tokenError(token.STRING)
		}
	}

	return array, nil
}

func (p *parser) parseIntArray() ([]int, error) {
	array := make([]int, 0, 5) // Make room for a conservitive number of elements

	for {
		val, err := strconv.Atoi(p.curTok.Literal)
		if err != nil {
			return nil, p.tokenError(token.INT)
		}

		array = append(array, val)
		p.nextToken()

		if p.curTok.Type == token.COMMA { // Commas are optional
			p.nextToken()
		}

		if p.curTok.Type == token.RSQUARE {
			break
		}

		if p.curTok.Type != token.INT {
			return nil, p.tokenError(token.INT)
		}
	}

	return array, nil
}

func (p *parser) tokenError(expected ...token.Type) error {
	allExpected := make([]string, len(expected))
	for i, t := range expected {
		allExpected[i] = t.String()
	}

	return fmt.Errorf("expected %s on line %d, got %s",
		strings.Join(allExpected, " or "), p.curTok.Line, p.curTok.Type.String())
}
