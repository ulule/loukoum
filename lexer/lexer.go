package lexer

import (
	"bufio"
	"io"

	"github.com/ulule/loukoum/token"
)

const (
	// EOF sentinel
	eof = 0
	// Newline sentinel
	newline = '\n'
)

// A Lexer will return a list of Token from a reader.
type Lexer struct {
	input  *bufio.Reader
	e0, en rune // current/next rune in reader
}

// New return a new Lexer from given source.
func New(input io.Reader) Lexer {

	buffer := bufio.NewReader(input)

	l := Lexer{}
	l.input = buffer

	// Read two runes so current and next items are both available.
	l.read()
	l.read()

	return l
}

func (l *Lexer) read() rune {
	e, _, err := l.input.ReadRune()

	if err != nil && err != io.EOF {
		//l.error(err.Error())
		// TODO FIX ME
		panic(err)
	}

	if err == io.EOF {
		e = eof
	}

	l.e0 = l.en

	l.en = e
	return l.e0
}

// current return the current rune in reader.
func (l *Lexer) current() rune {
	return l.e0
}

// next return the next rune in reader.
func (l *Lexer) next() rune {
	return l.en
}

// Iterator returns an Iteratee from reader.
func (l *Lexer) Iterator() *Iteratee {
	list := []token.Token{}
	for {
		current := l.Next()
		if current.Type == token.EOF {
			break
		}
		list = append(list, current)
	}
	return &Iteratee{list: list}
}

// Next will return the next token on reader.
func (l *Lexer) Next() token.Token {

	l.skipWhitespace()

	if l.current() == eof {
		return l.getToken(token.EOF)
	}

	if l.current() == newline {
		return l.getToken(token.Semicolon)
	}

	t, ok := l.getOperatorToken()
	if ok {
		return t
	}

	t, ok = l.getDelimiterToken()
	if ok {
		return t
	}

	return l.getDefaultToken()
}

func (l *Lexer) getToken(t token.Type) token.Token {
	defer l.read()
	switch t {
	case token.EOF:
		return token.New(t, "")
	case token.Semicolon:
		return token.New(t, ";")
	default:
		return token.New(t, string(l.current()))
	}
}

func (l *Lexer) getOperatorToken() (token.Token, bool) {
	switch l.current() {
	case '*':
		return l.getToken(token.Asterisk), true
	case '=':
		return l.getToken(token.Equals), true
	default:
		return token.Token{}, false
	}
}

func (l *Lexer) getDelimiterToken() (token.Token, bool) {
	switch l.current() {
	case ';':
		return l.getToken(token.Semicolon), true
	case ',':
		return l.getToken(token.Comma), true
	case ':':
		return l.getToken(token.Colon), true
	case '(':
		return l.getToken(token.LParen), true
	case ')':
		return l.getToken(token.RParen), true
	default:
		return token.Token{}, false
	}
}

func (l *Lexer) getDefaultToken() token.Token {

	if isLetter(l.current()) || isDigit(l.current()) {
		return l.getIdentifier()
	}

	return l.getToken(token.Illegal)
}

func (l *Lexer) getIdentifier() token.Token {

	t := token.Token{}
	v, ok := l.readIdentifier()

	t.Value = v
	t.Type = token.Lookup(t.Value)

	if !ok {
		t.Type = token.Illegal
	}

	return t
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.current()) {
		l.read()
	}
}

func (l *Lexer) readIdentifier() (string, bool) {

	buffer := []rune{}

	for isLetter(l.current()) || isDigit(l.current()) {
		buffer = append(buffer, l.current())
		l.read()
	}

	return string(buffer), true
}

func isLetter(e rune) bool {
	return 'a' <= e && e <= 'z' || 'A' <= e && e <= 'Z' || e == '_' || e == '.'
}

func isWhitespace(e rune) bool {
	return e == ' ' || e == '\t' || e == '\n' || e == '\r'
}

func isDigit(e rune) bool {
	return '0' <= e && e <= '9'
}

func isHex(e rune) bool {
	return isDigit(e) || 'a' <= e && e <= 'f' || 'A' <= e && e <= 'F'
}
