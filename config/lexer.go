package config

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type lexer struct {
	input  *bufio.Reader
	curCh  byte // current char under examination
	peekCh byte // peek character
}

func new(reader io.Reader) *lexer {
	l := &lexer{input: bufio.NewReader(reader)}
	// Populate both current and peek char
	l.readChar()
	l.readChar()
	return l
}

func newString(input string) *lexer {
	return new(strings.NewReader(input))
}

func (l *lexer) readChar() {
	l.curCh = l.peekCh

	var err error
	l.peekCh, err = l.input.ReadByte()
	if err != nil {
		l.peekCh = 0
	}
}

func (l *lexer) nextToken() token {
	var tok token

	l.devourWhitespace()

	switch l.curCh {
	// Operators
	case '+':
		tok = newByteToken(PLUS, l.curCh)
		break
	case '-':
		tok = newByteToken(MINUS, l.curCh)
		break
	case '*':
		tok = newByteToken(ASTERISK, l.curCh)
		break
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok = newToken(COMMENT, l.readSingleLineComment())
		} else if l.peekChar() == '*' {
			l.readChar()
			tok = newToken(COMMENT, l.readMultiLineComment())
		} else {
			tok = newByteToken(SLASH, l.curCh)
		}
		break
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(NOT_EQ, "!=")
		} else {
			tok = newByteToken(BANG, l.curCh)
		}
		break

	// Equality
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(EQ, "==")
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = newToken(ASSIGN, "=>")
		} else {
			tok = newByteToken(ILLEGAL, l.curCh)
		}
		break
	case '<':
		tok = newByteToken(LT, l.curCh)
		break
	case '>':
		tok = newByteToken(GT, l.curCh)
		break

	// Control characters
	case ',':
		tok = newByteToken(COMMA, l.curCh)
		break

	// Groupings
	case '{':
		tok = newByteToken(LBRACE, l.curCh)
		break
	case '}':
		tok = newByteToken(RBRACE, l.curCh)
		break
	case '[':
		tok = newByteToken(LSQUARE, l.curCh)
		break
	case ']':
		tok = newByteToken(RSQUARE, l.curCh)
		break

	case '"':
		tok = newToken(STRING, l.readString())
		break
	case '#':
		tok = newToken(COMMENT, l.readSingleLineComment())
		break
	case 0:
		tok = newToken(EOF, "")
		break

	default:
		if isLetter(l.curCh) {
			lit := l.readIdentifier()
			tok = newToken(lookupIdent(lit), lit)
			return tok
		} else if isDigit(l.curCh) {
			tok = l.readNumber()
			return tok
		}

		tok = newByteToken(ILLEGAL, l.curCh)
	}

	l.readChar()
	return tok
}

func (l *lexer) peekChar() byte {
	return l.peekCh
}

func (l *lexer) readIdentifier() string {
	var ident bytes.Buffer
	for isLetter(l.curCh) {
		ident.WriteByte(l.curCh)
		l.readChar()
	}
	return ident.String()
}

// TODO: Support escape sequences, standard Go should be fine, or PHP.
func (l *lexer) readString() string {
	var ident bytes.Buffer
	l.readChar() // Go past the starting double quote

	for l.curCh != '"' {
		ident.WriteByte(l.curCh)
		l.readChar()
	}

	return ident.String()
}

func (l *lexer) readNumber() token {
	var ident bytes.Buffer
	numTokenType := INT

	for isDigit(l.curCh) {
		// The parser will handle bad floats
		if l.curCh == '.' && numTokenType == INT {
			numTokenType = FLOAT
		}

		ident.WriteByte(l.curCh)
		l.readChar()
	}

	return newToken(tokenType(numTokenType), ident.String())
}

func (l *lexer) readSingleLineComment() string {
	var com bytes.Buffer
	l.readChar() // Go over # or / characters

	for l.curCh != '\n' && l.curCh != 0 {
		com.WriteByte(l.curCh)
		l.readChar()
	}
	return strings.TrimSpace(com.String())
}

func (l *lexer) readMultiLineComment() string {
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

func (l *lexer) devourWhitespace() {
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
