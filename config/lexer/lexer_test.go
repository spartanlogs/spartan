package lexer

import (
	"testing"

	"github.com/spartanlogs/spartan/config/token"
)

var testCases = []struct {
	input  string
	output []token.Token
}{
	{
		input: `input { }`,
		output: []token.Token{
			token.NewSimpleToken(token.INPUT, 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
		},
	},
	{
		input: `filter { }`,
		output: []token.Token{
			token.NewSimpleToken(token.FILTER, 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
		},
	},
	{
		input: `output { }`,
		output: []token.Token{
			token.NewSimpleToken(token.OUTPUT, 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
		},
	},
	{
		input: `# Hello, I'm a comment`,
		output: []token.Token{
			token.NewToken(token.COMMENT, "Hello, I'm a comment", 0, 0),
		},
	},
	{
		input: `// Hello, I'm a comment`,
		output: []token.Token{
			token.NewToken(token.COMMENT, "Hello, I'm a comment", 0, 0),
		},
	},
	{
		input: `/* Hello, I'm a
multiline comment */`,
		output: []token.Token{
			token.NewToken(token.COMMENT, "Hello, I'm a\nmultiline comment", 0, 0),
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
			token.NewSimpleToken(token.INPUT, 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewToken(token.IDENT, "file", 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewToken(token.IDENT, "path", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewToken(token.STRING, "/tmp/test", 0, 0),
			token.NewToken(token.IDENT, "tail_only", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewSimpleToken(token.TRUE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
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
			token.NewSimpleToken(token.FILTER, 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewToken(token.IDENT, "grok", 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewToken(token.IDENT, "field", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewToken(token.STRING, "message", 0, 0),
			token.NewToken(token.IDENT, "patterns", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewSimpleToken(token.LSQUARE, 0, 0),
			token.NewToken(token.STRING, "%{PATTERN1}", 0, 0),
			token.NewSimpleToken(token.COMMA, 0, 0),
			token.NewToken(token.STRING, "%{PATTERN2}", 0, 0),
			token.NewSimpleToken(token.RSQUARE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
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
			token.NewSimpleToken(token.FILTER, 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewToken(token.IDENT, "mutate", 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewToken(token.IDENT, "add_field", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewToken(token.STRING, "field1", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewToken(token.STRING, "new value", 0, 0),
			token.NewSimpleToken(token.COMMA, 0, 0),
			token.NewToken(token.STRING, "field2", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewToken(token.STRING, "other new value", 0, 0),
			token.NewSimpleToken(token.COMMA, 0, 0),
			token.NewToken(token.STRING, "field3", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewSimpleToken(token.FALSE, 0, 0),
			token.NewSimpleToken(token.COMMA, 0, 0),
			token.NewToken(token.STRING, "field4", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewToken(token.INT, "123456789", 0, 0),
			token.NewSimpleToken(token.COMMA, 0, 0),
			token.NewToken(token.STRING, "field5", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewToken(token.FLOAT, "3.14159", 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
		},
	},
	{
		input: `filter{mutate{add_field=>{"field1"=>"new value","field2"=>"other new value","field3"=>false,"field4"=>123456789,"field5"=>3.14159}}}`,
		output: []token.Token{
			token.NewSimpleToken(token.FILTER, 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewToken(token.IDENT, "mutate", 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewToken(token.IDENT, "add_field", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewSimpleToken(token.LBRACE, 0, 0),
			token.NewToken(token.STRING, "field1", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewToken(token.STRING, "new value", 0, 0),
			token.NewSimpleToken(token.COMMA, 0, 0),
			token.NewToken(token.STRING, "field2", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewToken(token.STRING, "other new value", 0, 0),
			token.NewSimpleToken(token.COMMA, 0, 0),
			token.NewToken(token.STRING, "field3", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewSimpleToken(token.FALSE, 0, 0),
			token.NewSimpleToken(token.COMMA, 0, 0),
			token.NewToken(token.STRING, "field4", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewToken(token.INT, "123456789", 0, 0),
			token.NewSimpleToken(token.COMMA, 0, 0),
			token.NewToken(token.STRING, "field5", 0, 0),
			token.NewSimpleToken(token.ASSIGN, 0, 0),
			token.NewToken(token.FLOAT, "3.14159", 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
			token.NewSimpleToken(token.RBRACE, 0, 0),
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
