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
		tok = newToken(PLUS, l.curCh)
		break
	case '-':
		tok = newToken(MINUS, l.curCh)
		break
	case '*':
		tok = newToken(ASTERISK, l.curCh)
		break
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok = token{
				Type:    COMMENT,
				Literal: l.readSingleLineComment(),
			}
		} else if l.peekChar() == '*' {
			l.readChar()
			tok = token{
				Type:    COMMENT,
				Literal: l.readMultiLineComment(),
			}
		} else {
			tok = newToken(SLASH, l.curCh)
		}
		break
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token{
				Type:    NOT_EQ,
				Literal: "!=",
			}
		} else {
			tok = newToken(BANG, l.curCh)
		}
		break

	// Equality
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token{
				Type:    EQ,
				Literal: "==",
			}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = token{
				Type:    ASSIGN,
				Literal: "=>",
			}
		} else {
			tok = newToken(ILLEGAL, l.curCh)
		}
		break
	case '<':
		tok = newToken(LT, l.curCh)
		break
	case '>':
		tok = newToken(GT, l.curCh)
		break

	// Control characters
	case ',':
		tok = newToken(COMMA, l.curCh)
		break
	case ';':
		tok = newToken(SEMICOLON, l.curCh)
		break
	case ':':
		tok = newToken(COLON, l.curCh)
		break

	// Groupings
	case '(':
		tok = newToken(LPAREN, l.curCh)
		break
	case ')':
		tok = newToken(RPAREN, l.curCh)
		break
	case '{':
		tok = newToken(LBRACE, l.curCh)
		break
	case '}':
		tok = newToken(RBRACE, l.curCh)
		break
	case '[':
		tok = newToken(LSQUARE, l.curCh)
		break
	case ']':
		tok = newToken(RSQUARE, l.curCh)
		break

	case '"':
		tok.Literal = l.readString()
		tok.Type = STRING
		break
	case '#':
		tok.Literal = l.readSingleLineComment()
		tok.Type = COMMENT
		break
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		break

	default:
		if isLetter(l.curCh) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			break
		} else if isDigit(l.curCh) {
			tok = l.readNumber()
			break
		}

		tok = newToken(ILLEGAL, l.curCh)
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

	return token{
		Type:    tokenType(numTokenType),
		Literal: ident.String(),
	}
}

func (l *lexer) readSingleLineComment() string {
	var com bytes.Buffer
	l.readChar() // Go over # or / characters

	for l.curCh != '\n' {
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
	return com.String()
}

func (l *lexer) devourWhitespace() {
	for isWhitespace(l.curCh) {
		l.readChar()
	}
}

func newToken(tokType tokenType, ch byte) token {
	return token{Type: tokType, Literal: string(ch)}
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
