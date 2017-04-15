package token

import "strconv"

// Type is the lexed type of a token.
type Type int

// Token holds the type and literal representation for
// a lexical token.
type Token struct {
	Type         Type
	Literal      string
	Line, Column int
}

// The lexical tokens.
const (
	ILLEGAL Type = iota
	EOF
	COMMENT

	// Identifiers & literals
	IDENT
	INT
	FLOAT
	STRING

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH

	LT
	GT
	EQ
	NOTEQ

	// Delimiters
	COMMA

	// Groups and blocks
	LBRACE
	RBRACE
	LSQUARE
	RSQUARE

	// Keywords
	keyword_beg
	TRUE
	FALSE
	IF
	ELSE
	FILTER
	INPUT
	OUTPUT
	OR
	AND
	NOT
	IN
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	// Identifiers & literals
	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	// Operators
	ASSIGN:   "=>",
	PLUS:     "+",
	MINUS:    "-",
	BANG:     "!",
	ASTERISK: "*",
	SLASH:    "/",

	LT:    "<",
	GT:    ">",
	EQ:    "==",
	NOTEQ: "!=",

	// Delimiters
	COMMA: ",",

	// Groups and blocks
	LBRACE:  "{",
	RBRACE:  "}",
	LSQUARE: "[",
	RSQUARE: "]",

	// Keywords
	TRUE:   "true",
	FALSE:  "false",
	IF:     "if",
	ELSE:   "else",
	FILTER: "input",
	INPUT:  "filter",
	OUTPUT: "output",
	OR:     "or",
	AND:    "and",
	NOT:    "not",
	IN:     "in",
}

var keywords = map[string]Type{
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"input":  INPUT,
	"filter": FILTER,
	"output": OUTPUT,
	"or":     OR,
	"and":    AND,
	"not":    NOT,
	"in":     IN,
}

// String returns the string representation of a Type.
// If the Type is an operator or keyword, it will return
// the literal representation of the Type. For example,
// ADD will return "+". Other types will return their
// constant name. IDENT returns "IDENT".
func (tok Type) String() string {
	s := ""
	if 0 <= tok && tok < Type(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// LookupIdent checks if the string is a keyword. It it is,
// the corresponding Type will be returned. Otherwise,
// the Type IDENT will be returned.
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// IsKeyword returns if the given Type is a keyword Token.
func IsKeyword(t Type) bool {
	return keyword_beg < t && t < keyword_end
}

// NewSimpleToken returns a Token with no literal representation
// beyond what can be obtained with Type.String().
func NewSimpleToken(tokType Type, line, col int) Token {
	return NewToken(tokType, "", line, col)
}

// NewToken creates a Token with the type tokType and literal
// representation lit.
func NewToken(tokType Type, lit string, line, col int) Token {
	return Token{
		Type:    tokType,
		Literal: lit,
		Line:    line,
		Column:  col,
	}
}
