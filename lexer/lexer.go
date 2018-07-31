package lexer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pkg/errors"

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
	e0, en rune  // current/next rune in reader
	err    error // error while reading
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
		l.error(errors.WithStack(err))
	}

	if err == io.EOF {
		e = eof
	}

	l.e0 = l.en

	l.en = e
	return l.e0
}

func (l *Lexer) error(err error) {
	l.err = err
}

// current return the current rune in reader.
func (l *Lexer) current() rune {
	return l.e0
}

// next return the next rune in reader.
// nolint: deadcode, megacheck
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

// Err will returns the error, if any, that was encountered during reading.
func (l *Lexer) Err() error {
	return l.err
}

// Next will return the next token on reader.
func (l *Lexer) Next() token.Token {
	if l.err != nil {
		return l.getToken(token.EOF)
	}

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

func (l *Lexer) getDelimiterToken() (token.Token, bool) { // nolint: gocyclo
	switch l.current() {
	case ';':
		return l.getToken(token.Semicolon), true
	case ',':
		return l.getToken(token.Comma), true
	case ':':
		if isLetter(l.next()) || isDigit(l.next()) {
			return l.getEscapedValue(), true
		}
		return l.getToken(token.Colon), true
	case '$':
		return l.getEscapedValue(), true
	case '(':
		return l.getToken(token.LParen), true
	case ')':
		return l.getToken(token.RParen), true
	case '"':
		return l.unwrapDelimiter('"', l.getIdentifier), true
	case '\'':
		return l.getString(), true
	default:
		return token.Token{}, false
	}
}

func (l *Lexer) unwrapDelimiter(delimiter rune, handler func() token.Token) token.Token {
	if l.current() != delimiter {
		l.error(errors.Errorf("identifier not started by %c", delimiter))
		return token.Token{
			Type:  token.Illegal,
			Value: fmt.Sprintf("%c", l.current()),
		}
	}

	l.read()
	t := handler()
	if t.Type == token.Illegal {
		return t
	}

	if l.current() != delimiter {
		l.error(errors.New("identifier not terminated"))
		return token.Token{
			Type:  token.Illegal,
			Value: fmt.Sprintf("%c", l.current()),
		}
	}

	l.read()
	return t
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

func (l *Lexer) getString() token.Token {
	t := token.Token{}
	v, ok := l.readString()

	t.Value = v
	t.Type = token.Literal
	if !ok {
		t.Type = token.Illegal
	}

	return t
}

func (l *Lexer) getEscapedValue() token.Token {
	t := token.Token{}
	v, ok := l.readEscapedValue()

	t.Value = v
	t.Type = token.Literal
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

func (l *Lexer) readEscapedValue() (string, bool) {

	buffer := []rune{l.current()}
	l.read()

	for isLetter(l.current()) || isDigit(l.current()) {
		buffer = append(buffer, l.current())
		l.read()
	}

	return string(buffer), true
}

func (l *Lexer) readString() (string, bool) {

	buffer := []rune{}
	failure := false

	l.read()
	buffer = append(buffer, '\'')

	for l.current() != '\'' && l.current() != '\n' && l.current() != eof {

		s, ok := l.escape('\'')
		if !ok {
			failure = true
		}

		buffer = append(buffer, s)
		l.read()
	}

	buffer = append(buffer, '\'')
	s := string(buffer)

	if l.current() != '\'' {
		l.error(errors.New("string not terminated"))
		return "\"" + s, false
	}

	l.read()
	return s, !failure
}

// nolint: gocyclo
func (l *Lexer) escape(quote rune) (rune, bool) {

	c := l.current()
	if c != '\\' {
		return c, true
	}

	switch l.next() {
	case 'n':
		c = '\n'
	case '\\':
		c = '\\'
	case 'a':
		c = '\a'
	case 'b':
		c = '\b'
	case 'f':
		c = '\f'
	case 'r':
		c = '\r'
	case 't':
		c = '\t'
	case 'v':
		c = '\v'
	case quote:
		c = quote
	default:
		l.error(errors.Errorf("unknown escape sequence: \\%s", l.next()))
		return c, false
	}

	l.read()
	return c, true
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

// nolint: deadcode, megacheck
func isHex(e rune) bool {
	return isDigit(e) || 'a' <= e && e <= 'f' || 'A' <= e && e <= 'F'
}
