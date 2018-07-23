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

	join, err := parseJoin(it, stmt.Join{
		Type: types.InnerJoin,
	})
	if err != nil {
		return stmt.Join{}, errors.Wrapf(err, "given query cannot be parsed: %s", subquery)
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

func parseJoin(it *lexer.Iteratee, join stmt.Join) (stmt.Join, error) { // nolint: gocyclo
	parenIdent := 0
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
					return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
				}

				// Right condition
				e = it.Next()
				right := stmt.NewColumn(e.Value)

				join.Condition = stmt.NewOnClause(left, right)

				for it.Is(token.RParen) {
					it.Next()
					parenIdent--
					if parenIdent < 0 {
						return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
					}
				}

				for it.Is(token.And) || it.Is(token.Or) {
					// We have an AND operator
					if it.Is(token.And) {
						e := it.Next()
						if e.Type != token.And || !it.Is(token.Literal) {
							return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
						}

						// Left condition
						e = it.Next()
						left := stmt.NewColumn(e.Value)

						// Check that we have a right condition
						e = it.Next()
						if e.Type != token.Equals || !it.Is(token.Literal) {
							return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
						}

						// Right condition
						e = it.Next()
						right := stmt.NewColumn(e.Value)

						join.Condition = join.Condition.And(stmt.NewOnClause(left, right))

						for it.Is(token.RParen) {
							it.Next()
							parenIdent--
							if parenIdent < 0 {
								return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
							}
						}
					}
					// We have an OR operator
					if it.Is(token.Or) {
						e := it.Next()
						if e.Type != token.Or || !it.Is(token.Literal) {
							return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
						}

						// Left condition
						e = it.Next()
						left := stmt.NewColumn(e.Value)

						// Check that we have a right condition
						e = it.Next()
						if e.Type != token.Equals || !it.Is(token.Literal) {
							return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
						}

						// Right condition
						e = it.Next()
						right := stmt.NewColumn(e.Value)

						join.Condition = join.Condition.Or(stmt.NewOnClause(left, right))
						for it.Is(token.RParen) {
							it.Next()
							parenIdent--
							if parenIdent < 0 {
								return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
							}
						}
					}
				}

				return join, nil
			}

		case token.On:
			continue

		case token.LParen:
			parenIdent++

		case token.RParen:
			parenIdent--
			if parenIdent < 0 {
				return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
			}

		default:
			return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
		}
	}

	return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
}
