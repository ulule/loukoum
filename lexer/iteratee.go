package lexer

import (
	"github.com/ulule/loukoum/token"
)

// An Iteratee contains a collection of token from a Lexer.
type Iteratee struct {
	cursor int
	list   []token.Token
}

// HasNext defines if a token is available.
func (it Iteratee) HasNext() bool {
	return it.cursor < len(it.list)
}

// Is defines if next token has given type.
func (it Iteratee) Is(next token.Type) bool {
	if !it.HasNext() {
		return false
	}
	return it.list[it.cursor].Type == next
}

// Next returns the next token.
func (it *Iteratee) Next() token.Token {
	element := it.list[it.cursor]
	it.cursor++
	return element
}
