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

				condition, err := parseJoinCondition(it, e.Value)
				if err != nil {
					return stmt.Join{}, err
				}

				join.Condition = condition

				return parseListJoinCondition(it, join)
			}

		case token.On:
			continue

		case token.LParen:
			op, err := parseWrappedJoinCondition(it, 0)
			if err != nil {
				return stmt.Join{}, err
			}

			e = it.Next()
			if e.Type != token.RParen {
				return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
			}

			join.Condition = op
			return parseListJoinCondition(it, join)

		default:
			return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
		}
	}

	return join, nil
}

func parseJoinCondition(it *lexer.Iteratee, name string) (stmt.OnClause, error) {
	name = strings.TrimSpace(name)
	if !it.HasNext() {
		return stmt.OnClause{}, errors.WithStack(ErrJoinInvalidCondition)
	}

	// Left condition
	left := stmt.NewColumn(name)

	// Check that we have a right condition
	e := it.Next()
	if e.Type != token.Equals || !it.Is(token.Literal) {
		return stmt.OnClause{}, errors.WithStack(ErrJoinInvalidCondition)
	}

	// Right condition
	e = it.Next()
	right := stmt.NewColumn(e.Value)

	condition := stmt.NewOnClause(left, right)
	return condition, nil
}

func parseWrappedJoinCondition(it *lexer.Iteratee, level int) (stmt.OnExpression, error) { // nolint: gocyclo
	var conditions stmt.OnExpression
	for it.HasNext() {
		e := it.Next()
		switch e.Type {
		case token.LParen:

			val, err := parseWrappedJoinCondition(it, level+1)
			if err != nil {
				return nil, err
			}
			if !it.Is(token.RParen) {
				return nil, errors.WithStack(ErrJoinInvalidCondition)
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
				return nil, errors.WithStack(ErrJoinInvalidCondition)
			}

			right, err := parseWrappedJoinCondition(it, level+1)
			if err != nil {
				return nil, err
			}

			left := conditions
			operator := stmt.NewOrOperator()
			conditions = stmt.NewInfixOnExpression(left, operator, right)

			if it.Is(token.RParen) {
				return conditions, nil
			}

		case token.And:
			if !it.Is(token.LParen) && !it.Is(token.Literal) {
				return nil, errors.WithStack(ErrJoinInvalidCondition)
			}

			right, err := parseWrappedJoinCondition(it, level+1)
			if err != nil {
				return nil, err
			}

			left := conditions
			operator := stmt.NewAndOperator()
			conditions = stmt.NewInfixOnExpression(left, operator, right)

			if it.Is(token.RParen) {
				return conditions, nil
			}

		case token.Literal:
			val, err := parseJoinCondition(it, e.Value)
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
	return nil, errors.WithStack(ErrJoinInvalidCondition)
}

func parseListJoinCondition(it *lexer.Iteratee, join stmt.Join) (stmt.Join, error) { // nolint: gocyclo
	for it.Is(token.And) || it.Is(token.Or) {
		// We have an AND operator
		if it.Is(token.And) {
			e := it.Next()
			if e.Type != token.And || !it.HasNext() {
				return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
			}
			e = it.Next()
			if e.Type == token.Literal {
				condition, err := parseJoinCondition(it, e.Value)
				if err != nil {
					return stmt.Join{}, err
				}
				join.Condition = join.Condition.And(condition)
			} else if e.Type == token.LParen {
				conditions, err := parseWrappedJoinCondition(it, 0)
				if err != nil {
					return stmt.Join{}, err
				}
				join.Condition = join.Condition.And(conditions)
			} else {
				return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
			}
		}
		// We have an OR operator
		if it.Is(token.Or) {
			e := it.Next()
			if e.Type != token.Or || !it.HasNext() {
				return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
			}
			e = it.Next()
			if e.Type == token.Literal {
				condition, err := parseJoinCondition(it, e.Value)
				if err != nil {
					return stmt.Join{}, err
				}
				join.Condition = join.Condition.Or(condition)
			} else if e.Type == token.LParen {
				conditions, err := parseWrappedJoinCondition(it, 0)
				if err != nil {
					return stmt.Join{}, err
				}
				join.Condition = join.Condition.Or(conditions)
			} else {
				return stmt.Join{}, errors.WithStack(ErrJoinInvalidCondition)
			}
		}
	}
	return join, nil
}
