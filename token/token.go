package token

import (
	"fmt"
	"strings"
)

type Type string

func (t Type) String() string {
	return string(t)
}

const (

	// Illegal is an unknown token type.
	Illegal = Type("Illegal")

	// EOF indicates the End-Of-File for the lexer.
	EOF = Type("EOF")

	// Literal defines entities such as columns, tables, etc...
	Literal = Type("Literal")

	// Delimiters
	Comma     = Type(",")
	Semicolon = Type(";")
	Colon     = Type(":")
	LParen    = Type("(")
	RParen    = Type(")")

	Equals   = Type("=")
	Asterisk = Type("*")

	// Keywords
	Select    = Type("SELECT")
	Update    = Type("UPDATE")
	Insert    = Type("INSERT")
	Delete    = Type("DELETE")
	From      = Type("FROM")
	Where     = Type("WHERE")
	And       = Type("AND")
	Or        = Type("OR")
	Limit     = Type("LIMIT")
	Offset    = Type("OFFSET")
	Set       = Type("SET")
	As        = Type("AS")
	Inner     = Type("INNER")
	Cross     = Type("CROSS")
	Left      = Type("LEFT")
	Right     = Type("RIGHT")
	Join      = Type("JOIN")
	On        = Type("ON")
	Group     = Type("GROUP")
	By        = Type("BY")
	Having    = Type("HAVING")
	Order     = Type("ORDER")
	Distinct  = Type("DISTINCT")
	Only      = Type("ONLY")
	Returning = Type("RETURNING")
	Values    = Type("VALUES")
	Into      = Type("INTO")
)

type Token struct {
	Type  Type
	Value string
}

func (t *Token) String() string {
	return fmt.Sprintf(`{"%s": "%s"}`, t.Type, t.Value)
}

var keywords = map[string]Type{
	"SELECT":    Select,
	"UPDATE":    Update,
	"INSERT":    Insert,
	"DELETE":    Delete,
	"FROM":      From,
	"WHERE":     Where,
	"AND":       And,
	"OR":        Or,
	"LIMIT":     Limit,
	"OFFSET":    Offset,
	"SET":       Set,
	"AS":        As,
	"INNER":     Inner,
	"CROSS":     Cross,
	"LEFT":      Left,
	"RIGHT":     Right,
	"JOIN":      Join,
	"ON":        On,
	"GROUP":     Group,
	"BY":        By,
	"HAVING":    Having,
	"ORDER":     Order,
	"DISTINCT":  Distinct,
	"ONLY":      Only,
	"RETURNING": Returning,
	"VALUES":    Values,
	"INTO":      Into,
}

func Lookup(e string) Type {
	n := strings.ToUpper(e)
	t, ok := keywords[n]
	if ok {
		return t
	}
	return Literal
}

func New(t Type, v string) Token {
	return Token{
		Type:  t,
		Value: v,
	}
}
