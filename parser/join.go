package parser

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ulule/loukoum/lexer"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// ErrJoinInvalidCondition is returned when join condition cannot be parsed.
var ErrJoinInvalidCondition = fmt.Errorf("join condition is invalid")

// ParseJoin will try to parse given subquery as a join statement.
func ParseJoin(subquery string) (stmt.Join, error) {
	lexer := lexer.New(strings.NewReader(subquery))
	it := lexer.Iterator()

	join := stmt.Join{
		Type: types.InnerJoin,
	}

	for it.HasNext() {
		e := it.Next()

		// Parse join type
		if e.Type == token.Inner && it.Peek(token.Join) {
			join.Type = types.InnerJoin
		}
		if e.Type == token.Left && it.Peek(token.Join) {
			join.Type = types.LeftJoin
		}
		if e.Type == token.Right && it.Peek(token.Join) {
			join.Type = types.RightJoin
		}

		// Parse join table
		if e.Type == token.Literal && it.Peek(token.On) {
			join.Table = e.Value
		}

		// Parse join condition
		if e.Type == token.Literal && it.Peek(token.Equals) {

			// Left condition
			left := stmt.NewColumn(e.Value)

			// Check that we have a right condition
			e = it.Next()
			if e.Type != token.Equals && !it.Peek(token.Literal) {
				err := errors.Wrapf(ErrJoinInvalidCondition, "given query cannot be parsed: %s", subquery)
				return stmt.Join{}, err
			}

			// Right condition
			e = it.Next()
			right := stmt.NewColumn(e.Value)

			join.Condition = stmt.NewOn(left, right)
		}
	}

	return join, nil
}

// MustParseJoin will execute ParseJoin and panic on error.
func MustParseJoin(subquery string) stmt.Join {
	join, err := ParseJoin(subquery)
	if err != nil {
		panic(fmt.Sprintf("loukoum: %s", err))
	}
	return join
}
