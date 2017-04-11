package config

type tokenType string

type token struct {
	Type    tokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	COMMENT = "COMMENT"

	// Identifiers & literals
	IDENT  = "IDENT"  // add, foobar, x, y
	INT    = "INT"    // 1343546
	FLOAT  = "FLOAT"  // 12.52
	STRING = "STRING" // "some text"

	// Operators
	ASSIGN   = "=>"
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA = ","

	// Groups and blocks
	LBRACE  = "{"
	RBRACE  = "}"
	LSQUARE = "["
	RSQUARE = "]"

	// Keywords
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	FILTER = "FILTER"
	INPUT  = "INPUT"
	OUTPUT = "OUTPUT"
	OR     = "OR"
	AND    = "AND"
	NOT    = "NOT"
	IN     = "IN"
)

var keywords = map[string]tokenType{
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

func lookupIdent(ident string) tokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func newByteToken(tokType tokenType, ch byte) token {
	return newToken(tokType, string(ch))
}

func newToken(tokType tokenType, lit string) token {
	return token{Type: tokType, Literal: lit}
}
