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
func ParseJoin(subquery string) (stmt.Join, error) { // nolint: gocyclo
	lexer := lexer.New(strings.NewReader(subquery))
	it := lexer.Iterator()

	join := stmt.Join{
		Type: types.InnerJoin,
	}

	for it.HasNext() {
		e := it.Next()

		switch e.Type {
		// Parse join type
		case token.Join:
			continue
		case token.Inner:
			if it.Is(token.Join) {
				it.Next()
				join.Type = types.InnerJoin
				continue
			}
		case token.Left:
			if it.Is(token.Join) {
				it.Next()
				join.Type = types.LeftJoin
				continue
			}
		case token.Right:
			if it.Is(token.Join) {
				it.Next()
				join.Type = types.RightJoin
				continue
			}
		case token.Literal:
			// Parse join table
			if it.Is(token.On) {
				it.Next()
				join.Table = stmt.NewTable(e.Value)
				continue
			}

			// Parse join condition
			if it.Is(token.Equals) {

				// Left condition
				left := stmt.NewColumn(e.Value)

				// Check that we have a right condition
				e = it.Next()
				if e.Type != token.Equals || !it.Is(token.Literal) {
					err := errors.Wrapf(ErrJoinInvalidCondition, "given query cannot be parsed: %s", subquery)
					return stmt.Join{}, err
				}

				// Right condition
				e = it.Next()
				right := stmt.NewColumn(e.Value)

				join.Condition = stmt.NewOn(left, right)

				for it.Is(token.And) || it.Is(token.Or) {
					// We have an AND operator
					if it.Is(token.And) {
						e := it.Next()
						if e.Type != token.And || !it.Is(token.Literal) {
							err := errors.Wrapf(ErrJoinInvalidCondition, "given query cannot be parsed: %s", subquery)
							return stmt.Join{}, err
						}

						// Left condition
						e = it.Next()
						left := stmt.NewColumn(e.Value)

						// Check that we have a right condition
						e = it.Next()
						if e.Type != token.Equals || !it.Is(token.Literal) {
							err := errors.Wrapf(ErrJoinInvalidCondition, "given query cannot be parsed: %s", subquery)
							return stmt.Join{}, err
						}

						// Right condition
						e = it.Next()
						right := stmt.NewColumn(e.Value)

						join.Condition = stmt.NewInfixExpression(
							join.Condition, stmt.NewAndOperator(), stmt.NewOn(left, right),
						)

					}
					// We have an OR operator
					if it.Is(token.Or) {
						e := it.Next()
						if e.Type != token.Or || !it.Is(token.Literal) {
							err := errors.Wrapf(ErrJoinInvalidCondition, "given query cannot be parsed: %s", subquery)
							return stmt.Join{}, err
						}

						// Left condition
						e = it.Next()
						left := stmt.NewColumn(e.Value)

						// Check that we have a right condition
						e = it.Next()
						if e.Type != token.Equals || !it.Is(token.Literal) {
							err := errors.Wrapf(ErrJoinInvalidCondition, "given query cannot be parsed: %s", subquery)
							return stmt.Join{}, err
						}

						// Right condition
						e = it.Next()
						right := stmt.NewColumn(e.Value)

						join.Condition = stmt.NewInfixExpression(
							join.Condition, stmt.NewOrOperator(), stmt.NewOn(left, right),
						)
					}
				}

				continue
			}

		case token.On, token.LParen, token.RParen:
			continue
		}

		// Ignore invalid token and stop iterating.
		break
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
