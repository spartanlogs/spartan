package lexer

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/spartanlogs/spartan/config/token"
)

type Lexer struct {
	input  *bufio.Reader
	curCh  byte // current char under examination
	peekCh byte // peek character
	line   int  // line in source file
	column int  // column in current line
}

func New(reader io.Reader) *Lexer {
	l := &Lexer{
		input:  bufio.NewReader(reader),
		line:   1,
		column: 1,
	}
	// Populate both current and peek char
	l.readChar()
	l.readChar()
	return l
}

func NewString(input string) *Lexer {
	return New(strings.NewReader(input))
}

func (l *Lexer) readChar() {
	l.curCh = l.peekCh

	var err error
	l.peekCh, err = l.input.ReadByte()
	if err != nil {
		l.peekCh = 0
	}

	if l.curCh == '\r' {
		l.readChar()
	} else if l.curCh == '\n' {
		l.line++
		l.column = 0
	}
	l.column++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.devourWhitespace()

	switch l.curCh {
	// Operators
	case '+':
		tok = token.NewSimpleToken(token.PLUS, l.line, l.column)
	case '-':
		tok = token.NewSimpleToken(token.MINUS, l.line, l.column)
	case '*':
		tok = token.NewSimpleToken(token.ASTERISK, l.line, l.column)
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok = token.NewToken(token.COMMENT, l.readSingleLineComment(), l.line, l.column)
		} else if l.peekChar() == '*' {
			l.readChar()
			tok = token.NewToken(token.COMMENT, l.readMultiLineComment(), l.line, l.column)
		} else {
			tok = token.NewSimpleToken(token.SLASH, l.line, l.column)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.NewSimpleToken(token.NOTEQ, l.line, l.column)
		} else {
			tok = token.NewSimpleToken(token.BANG, l.line, l.column)
		}

	// Equality
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.NewSimpleToken(token.EQ, l.line, l.column)
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = token.NewSimpleToken(token.ASSIGN, l.line, l.column)
		} else {
			tok = token.NewSimpleToken(token.ILLEGAL, l.line, l.column)
		}
	case '<':
		tok = token.NewSimpleToken(token.LT, l.line, l.column)
	case '>':
		tok = token.NewSimpleToken(token.GT, l.line, l.column)

	// Control characters
	case ',':
		tok = token.NewSimpleToken(token.COMMA, l.line, l.column)

	// Groupings
	case '{':
		tok = token.NewSimpleToken(token.LBRACE, l.line, l.column)
	case '}':
		tok = token.NewSimpleToken(token.RBRACE, l.line, l.column)
	case '[':
		tok = token.NewSimpleToken(token.LSQUARE, l.line, l.column)
	case ']':
		tok = token.NewSimpleToken(token.RSQUARE, l.line, l.column)

	case '"':
		tok = token.NewToken(token.STRING, l.readString(), l.line, l.column)
	case '#':
		tok = token.NewToken(token.COMMENT, l.readSingleLineComment(), l.line, l.column)
	case 0:
		tok = token.NewSimpleToken(token.EOF, l.line, l.column)

	default:
		if isLetter(l.curCh) {
			lit := l.readIdentifier()
			tokType := token.LookupIdent(lit)
			if token.IsKeyword(tokType) { // No need to save the literal keyword
				tok = token.NewSimpleToken(tokType, l.line, l.column)
			} else {
				tok = token.NewToken(tokType, lit, l.line, l.column)
			}
			return tok
		} else if isDigit(l.curCh) {
			tok = l.readNumber()
			return tok
		}

		tok = token.NewSimpleToken(token.ILLEGAL, l.line, l.column)
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	return l.peekCh
}

func (l *Lexer) readIdentifier() string {
	var ident bytes.Buffer
	for isLetter(l.curCh) {
		ident.WriteByte(l.curCh)
		l.readChar()
	}
	return ident.String()
}

// TODO: Support escape sequences, standard Go should be fine, or PHP.
func (l *Lexer) readString() string {
	var ident bytes.Buffer
	l.readChar() // Go past the starting double quote

	for l.curCh != '"' {
		ident.WriteByte(l.curCh)
		l.readChar()
	}

	return ident.String()
}

func (l *Lexer) readNumber() token.Token {
	var ident bytes.Buffer
	numTokenType := token.INT

	for isDigit(l.curCh) {
		// The parser will handle bad floats
		if l.curCh == '.' && numTokenType == token.INT {
			numTokenType = token.FLOAT
		}

		ident.WriteByte(l.curCh)
		l.readChar()
	}

	return token.NewToken(numTokenType, ident.String(), l.line, l.column)
}

func (l *Lexer) readSingleLineComment() string {
	var com bytes.Buffer
	l.readChar() // Go over # or / characters

	for l.curCh != '\n' && l.curCh != 0 {
		com.WriteByte(l.curCh)
		l.readChar()
	}
	return strings.TrimSpace(com.String())
}

func (l *Lexer) readMultiLineComment() string {
	var com bytes.Buffer
	l.readChar() // Go over * character

	for l.curCh != 0 {
		if l.curCh == '*' && l.peekChar() == '/' {
			l.readChar() // Skip *
			break
		}

		com.WriteByte(l.curCh)
		l.readChar()
	}
	return strings.TrimSpace(com.String())
}

func (l *Lexer) devourWhitespace() {
	for isWhitespace(l.curCh) {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9') || ch == '.'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
