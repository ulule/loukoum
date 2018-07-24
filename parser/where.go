package parser

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ulule/loukoum/lexer"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/token"
)

// ErrWhereInvalidCondition is returned when where condition cannot be parsed.
var ErrWhereInvalidCondition = fmt.Errorf("where condition is invalid")

func parseWhere(it *lexer.Iteratee) (stmt.Where, error) { // nolint: gocyclo
	if !it.Is(token.Where) && !it.Is(token.Literal) && !it.Is(token.LParen) {
		return stmt.Where{}, errors.WithStack(ErrWhereInvalidCondition)
	}

	query := stmt.Where{}

	for it.HasNext() {
		e := it.Next()
		switch e.Type {
		case token.Where:
			continue

		case token.Literal:

			op, err := parseWhereCondition(it, e.Value)
			if err != nil {
				return stmt.Where{}, err
			}

			query.Condition = op
			if !it.HasNext() {
				return query, nil
			}

			for it.Is(token.And) || it.Is(token.Or) {
				e := it.Next()
				// We have an AND operator
				if e.Type == token.And {
					e := it.Next()
					switch e.Type {
					case token.LParen:
						op, err := parseWrappedWhereCondition(it, 0)
						if err != nil {
							return stmt.Where{}, err
						}

						e = it.Next()
						if e.Type != token.RParen {
							return stmt.Where{}, errors.WithStack(ErrWhereInvalidCondition)
						}

						query = query.And(op)

					case token.Literal:
						op, err := parseWhereCondition(it, e.Value)
						if err != nil {
							return stmt.Where{}, err
						}

						query = query.And(op)

					default:
						return stmt.Where{}, errors.WithStack(ErrWhereInvalidCondition)
					}
				}
				// We have an OR operator
				if e.Type == token.Or {
					e := it.Next()
					switch e.Type {
					case token.LParen:
						op, err := parseWrappedWhereCondition(it, 0)
						if err != nil {
							return stmt.Where{}, err
						}

						e = it.Next()
						if e.Type != token.RParen {
							return stmt.Where{}, errors.WithStack(ErrWhereInvalidCondition)
						}

						query = query.Or(op)

					case token.Literal:
						op, err := parseWhereCondition(it, e.Value)
						if err != nil {
							return stmt.Where{}, err
						}

						query = query.Or(op)

					default:
						return stmt.Where{}, errors.WithStack(ErrWhereInvalidCondition)
					}
				}
			}
			if !it.HasNext() {
				return query, nil
			}

		case token.LParen:
			op, err := parseWrappedWhereCondition(it, 0)
			if err != nil {
				return stmt.Where{}, err
			}

			e = it.Next()
			if e.Type != token.RParen {
				return stmt.Where{}, errors.WithStack(ErrWhereInvalidCondition)
			}

			query.Condition = op
			if !it.HasNext() {
				// TODO Fix
				return query, nil
			}

		default:
			return stmt.Where{}, errors.WithStack(ErrWhereInvalidCondition)
		}
	}

	return stmt.Where{}, errors.WithStack(ErrWhereInvalidCondition)
}

func parseWhereCondition(it *lexer.Iteratee, name string) (stmt.Expression, error) {
	name = strings.TrimSpace(name)
	if !it.HasNext() {
		return nil, errors.WithStack(ErrWhereInvalidCondition)
	}

	op := it.Next()
	switch op.Type {
	case token.Equals:
		if !it.HasNext() || !it.Is(token.Literal) {
			e := it.Next()
			// TODO Remove debug
			fmt.Println(e.Type)
			fmt.Println(e)
			return nil, errors.WithStack(ErrWhereInvalidCondition)
		}

		right := it.Next()
		val := stmt.NewIdentifier(name).Equal(stmt.NewRaw(strings.TrimSpace(right.Value)))
		return val, nil

	default:
		return nil, errors.WithStack(ErrWhereInvalidCondition)
	}
}

func parseWrappedWhereCondition(it *lexer.Iteratee, level int) (stmt.Expression, error) { // nolint: gocyclo
	var conditions stmt.Expression
	for it.HasNext() {
		e := it.Next()
		switch e.Type {
		case token.LParen:

			val, err := parseWrappedWhereCondition(it, level+1)
			if err != nil {
				return nil, err
			}
			if !it.Is(token.RParen) {
				return nil, errors.WithStack(ErrWhereInvalidCondition)
			}

			it.Next()
			if conditions == nil {
				conditions = val
			}
			if it.Is(token.RParen) {
				return conditions, nil
			}

		case token.Or:
			if !it.Is(token.LParen) && !it.Is(token.Literal) {
				return nil, errors.WithStack(ErrWhereInvalidCondition)
			}

			right, err := parseWrappedWhereCondition(it, level+1)
			if err != nil {
				return nil, err
			}

			left := conditions
			operator := stmt.NewOrOperator()
			conditions = stmt.NewInfixExpression(left, operator, right)

			if it.Is(token.RParen) {
				return conditions, nil
			}

		case token.And:
			if !it.Is(token.LParen) && !it.Is(token.Literal) {
				return nil, errors.WithStack(ErrWhereInvalidCondition)
			}

			right, err := parseWrappedWhereCondition(it, level+1)
			if err != nil {
				return nil, err
			}

			left := conditions
			operator := stmt.NewAndOperator()
			conditions = stmt.NewInfixExpression(left, operator, right)

			if it.Is(token.RParen) {
				return conditions, nil
			}

		case token.Literal:
			val, err := parseWhereCondition(it, e.Value)
			if err != nil {
				return nil, err
			}
			if conditions == nil {
				conditions = val
			}
			if it.Is(token.RParen) {
				return conditions, nil
			}
		}
	}
	return nil, errors.WithStack(ErrWhereInvalidCondition)
}
