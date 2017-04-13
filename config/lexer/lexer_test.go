package lexer

import (
	"testing"

	"github.com/lfkeitel/spartan/config/token"
)

var testCases = []struct {
	input  string
	output []token.Token
}{
	{
		input: `input { }`,
		output: []token.Token{
			token.NewSimpleToken(token.INPUT),
			token.NewSimpleToken(token.LBRACE),
			token.NewSimpleToken(token.RBRACE),
		},
	},
	{
		input: `filter { }`,
		output: []token.Token{
			token.NewSimpleToken(token.FILTER),
			token.NewSimpleToken(token.LBRACE),
			token.NewSimpleToken(token.RBRACE),
		},
	},
	{
		input: `output { }`,
		output: []token.Token{
			token.NewSimpleToken(token.OUTPUT),
			token.NewSimpleToken(token.LBRACE),
			token.NewSimpleToken(token.RBRACE),
		},
	},
	{
		input: `# Hello, I'm a comment`,
		output: []token.Token{
			token.NewToken(token.COMMENT, "Hello, I'm a comment"),
		},
	},
	{
		input: `// Hello, I'm a comment`,
		output: []token.Token{
			token.NewToken(token.COMMENT, "Hello, I'm a comment"),
		},
	},
	{
		input: `/* Hello, I'm a
multiline comment */`,
		output: []token.Token{
			token.NewToken(token.COMMENT, "Hello, I'm a\nmultiline comment"),
		},
	},
	{
		input: `input {
	file {
		path => "/tmp/test"
		tail_only => true
	}
}`,
		output: []token.Token{
			token.NewSimpleToken(token.INPUT),
			token.NewSimpleToken(token.LBRACE),
			token.NewToken(token.IDENT, "file"),
			token.NewSimpleToken(token.LBRACE),
			token.NewToken(token.IDENT, "path"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewToken(token.STRING, "/tmp/test"),
			token.NewToken(token.IDENT, "tail_only"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewSimpleToken(token.TRUE),
			token.NewSimpleToken(token.RBRACE),
			token.NewSimpleToken(token.RBRACE),
		},
	},
	{
		input: `filter {
	grok {
		field => "message"
		patterns => [
			"%{PATTERN1}",
			"%{PATTERN2}"
		]
	}
}`,
		output: []token.Token{
			token.NewSimpleToken(token.FILTER),
			token.NewSimpleToken(token.LBRACE),
			token.NewToken(token.IDENT, "grok"),
			token.NewSimpleToken(token.LBRACE),
			token.NewToken(token.IDENT, "field"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewToken(token.STRING, "message"),
			token.NewToken(token.IDENT, "patterns"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewSimpleToken(token.LSQUARE),
			token.NewToken(token.STRING, "%{PATTERN1}"),
			token.NewSimpleToken(token.COMMA),
			token.NewToken(token.STRING, "%{PATTERN2}"),
			token.NewSimpleToken(token.RSQUARE),
			token.NewSimpleToken(token.RBRACE),
			token.NewSimpleToken(token.RBRACE),
		},
	},
	{
		input: `filter {
	mutate {
		add_field => {
			"field1" => "new value",
			"field2" => "other new value",
			"field3" => false,
			"field4" => 123456789,
			"field5" => 3.14159
		}
	}
}`,
		output: []token.Token{
			token.NewSimpleToken(token.FILTER),
			token.NewSimpleToken(token.LBRACE),
			token.NewToken(token.IDENT, "mutate"),
			token.NewSimpleToken(token.LBRACE),
			token.NewToken(token.IDENT, "add_field"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewSimpleToken(token.LBRACE),
			token.NewToken(token.STRING, "field1"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewToken(token.STRING, "new value"),
			token.NewSimpleToken(token.COMMA),
			token.NewToken(token.STRING, "field2"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewToken(token.STRING, "other new value"),
			token.NewSimpleToken(token.COMMA),
			token.NewToken(token.STRING, "field3"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewSimpleToken(token.FALSE),
			token.NewSimpleToken(token.COMMA),
			token.NewToken(token.STRING, "field4"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewToken(token.INT, "123456789"),
			token.NewSimpleToken(token.COMMA),
			token.NewToken(token.STRING, "field5"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewToken(token.FLOAT, "3.14159"),
			token.NewSimpleToken(token.RBRACE),
			token.NewSimpleToken(token.RBRACE),
			token.NewSimpleToken(token.RBRACE),
		},
	},
	{
		input: `filter{mutate{add_field=>{"field1"=>"new value","field2"=>"other new value","field3"=>false,"field4"=>123456789,"field5"=>3.14159}}}`,
		output: []token.Token{
			token.NewSimpleToken(token.FILTER),
			token.NewSimpleToken(token.LBRACE),
			token.NewToken(token.IDENT, "mutate"),
			token.NewSimpleToken(token.LBRACE),
			token.NewToken(token.IDENT, "add_field"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewSimpleToken(token.LBRACE),
			token.NewToken(token.STRING, "field1"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewToken(token.STRING, "new value"),
			token.NewSimpleToken(token.COMMA),
			token.NewToken(token.STRING, "field2"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewToken(token.STRING, "other new value"),
			token.NewSimpleToken(token.COMMA),
			token.NewToken(token.STRING, "field3"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewSimpleToken(token.FALSE),
			token.NewSimpleToken(token.COMMA),
			token.NewToken(token.STRING, "field4"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewToken(token.INT, "123456789"),
			token.NewSimpleToken(token.COMMA),
			token.NewToken(token.STRING, "field5"),
			token.NewSimpleToken(token.ASSIGN),
			token.NewToken(token.FLOAT, "3.14159"),
			token.NewSimpleToken(token.RBRACE),
			token.NewSimpleToken(token.RBRACE),
			token.NewSimpleToken(token.RBRACE),
		},
	},
}

func TestLexer(t *testing.T) {
	for testNum, test := range testCases {
		l := NewString(test.input)
		tokenCount := 0
		for {
			next := l.NextToken()
			if next.Type == token.EOF {
				break
			}

			tok := test.output[tokenCount]
			if next.Type != tok.Type {
				t.Errorf(
					"Test %d errored on token %d. Expected type \"%s\", got \"%s\"",
					testNum+1, tokenCount+1, tok.Type, next.Type,
				)
			}

			if next.Literal != tok.Literal {
				t.Errorf(
					"Test %d errored on token %d. Expected literal \"%s\", got \"%s\"",
					testNum+1, tokenCount+1, tok.Literal, next.Literal,
				)
			}
			tokenCount++
		}

		// The extra token accounts for EOF
		if tokenCount != len(test.output) {
			t.Errorf("Test %d errored. Expected %d tokens, got %d",
				testNum+1, len(test.output), tokenCount+1,
			)
		}
	}
}
