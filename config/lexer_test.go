package config

import (
	"testing"
)

var testCases = []struct {
	input  string
	output []token
}{
	{
		input: `input { }`,
		output: []token{
			newToken(INPUT, "input"),
			newByteToken(LBRACE, '{'),
			newByteToken(RBRACE, '}'),
		},
	},
	{
		input: `filter { }`,
		output: []token{
			newToken(FILTER, "filter"),
			newByteToken(LBRACE, '{'),
			newByteToken(RBRACE, '}'),
		},
	},
	{
		input: `output { }`,
		output: []token{
			newToken(OUTPUT, "output"),
			newByteToken(LBRACE, '{'),
			newByteToken(RBRACE, '}'),
		},
	},
	{
		input: `# Hello, I'm a comment`,
		output: []token{
			newToken(COMMENT, "Hello, I'm a comment"),
		},
	},
	{
		input: `// Hello, I'm a comment`,
		output: []token{
			newToken(COMMENT, "Hello, I'm a comment"),
		},
	},
	{
		input: `/* Hello, I'm a
multiline comment */`,
		output: []token{
			newToken(COMMENT, "Hello, I'm a\nmultiline comment"),
		},
	},
	{
		input: `input {
	file {
		path => "/tmp/test"
		tail_only => true
	}
}`,
		output: []token{
			newToken(INPUT, "input"),
			newByteToken(LBRACE, '{'),
			newToken(IDENT, "file"),
			newByteToken(LBRACE, '{'),
			newToken(IDENT, "path"),
			newToken(ASSIGN, "=>"),
			newToken(STRING, "/tmp/test"),
			newToken(IDENT, "tail_only"),
			newToken(ASSIGN, "=>"),
			newToken(TRUE, "true"),
			newByteToken(RBRACE, '}'),
			newByteToken(RBRACE, '}'),
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
		output: []token{
			newToken(FILTER, "filter"),
			newByteToken(LBRACE, '{'),
			newToken(IDENT, "grok"),
			newByteToken(LBRACE, '{'),
			newToken(IDENT, "field"),
			newToken(ASSIGN, "=>"),
			newToken(STRING, "message"),
			newToken(IDENT, "patterns"),
			newToken(ASSIGN, "=>"),
			newByteToken(LSQUARE, '['),
			newToken(STRING, "%{PATTERN1}"),
			newByteToken(COMMA, ','),
			newToken(STRING, "%{PATTERN2}"),
			newByteToken(RSQUARE, ']'),
			newByteToken(RBRACE, '}'),
			newByteToken(RBRACE, '}'),
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
		output: []token{
			newToken(FILTER, "filter"),
			newByteToken(LBRACE, '{'),
			newToken(IDENT, "mutate"),
			newByteToken(LBRACE, '{'),
			newToken(IDENT, "add_field"),
			newToken(ASSIGN, "=>"),
			newByteToken(LBRACE, '{'),
			newToken(STRING, "field1"),
			newToken(ASSIGN, "=>"),
			newToken(STRING, "new value"),
			newByteToken(COMMA, ','),
			newToken(STRING, "field2"),
			newToken(ASSIGN, "=>"),
			newToken(STRING, "other new value"),
			newByteToken(COMMA, ','),
			newToken(STRING, "field3"),
			newToken(ASSIGN, "=>"),
			newToken(FALSE, "false"),
			newByteToken(COMMA, ','),
			newToken(STRING, "field4"),
			newToken(ASSIGN, "=>"),
			newToken(INT, "123456789"),
			newByteToken(COMMA, ','),
			newToken(STRING, "field5"),
			newToken(ASSIGN, "=>"),
			newToken(FLOAT, "3.14159"),
			newByteToken(RBRACE, '}'),
			newByteToken(RBRACE, '}'),
			newByteToken(RBRACE, '}'),
		},
	},
	{
		input: `filter{mutate{add_field=>{"field1"=>"new value","field2"=>"other new value","field3"=>false,"field4"=>123456789,"field5"=>3.14159}}}`,
		output: []token{
			newToken(FILTER, "filter"),
			newByteToken(LBRACE, '{'),
			newToken(IDENT, "mutate"),
			newByteToken(LBRACE, '{'),
			newToken(IDENT, "add_field"),
			newToken(ASSIGN, "=>"),
			newByteToken(LBRACE, '{'),
			newToken(STRING, "field1"),
			newToken(ASSIGN, "=>"),
			newToken(STRING, "new value"),
			newByteToken(COMMA, ','),
			newToken(STRING, "field2"),
			newToken(ASSIGN, "=>"),
			newToken(STRING, "other new value"),
			newByteToken(COMMA, ','),
			newToken(STRING, "field3"),
			newToken(ASSIGN, "=>"),
			newToken(FALSE, "false"),
			newByteToken(COMMA, ','),
			newToken(STRING, "field4"),
			newToken(ASSIGN, "=>"),
			newToken(INT, "123456789"),
			newByteToken(COMMA, ','),
			newToken(STRING, "field5"),
			newToken(ASSIGN, "=>"),
			newToken(FLOAT, "3.14159"),
			newByteToken(RBRACE, '}'),
			newByteToken(RBRACE, '}'),
			newByteToken(RBRACE, '}'),
		},
	},
}

func TestLexer(t *testing.T) {
	for testNum, test := range testCases {
		l := newString(test.input)
		tokenCount := 0
		for {
			next := l.nextToken()
			if next.Type == EOF {
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
