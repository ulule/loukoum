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

// ErrSelectInvalidCondition is returned when select condition cannot be parsed.
var ErrSelectInvalidCondition = fmt.Errorf("select condition is invalid")

// Parse will try to parse given query as a statement.
func parseSelect(it *lexer.Iteratee) (stmt.Select, error) { // nolint: gocyclo
	if !it.Is(token.Select) {
		return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
	}

	query := stmt.Select{}
	it.Next()

	for it.HasNext() {
		e := it.Next()
		switch e.Type {
		// Parse columns
		case token.Distinct:

			query.Distinct = true
			if !it.Is(token.On) && !it.Is(token.Asterisk) && !it.Is(token.Literal) {
				return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
			}

		case token.Asterisk:

			query.Columns = append(query.Columns, stmt.NewColumn(strings.TrimSpace(e.Value)))
			if it.HasNext() && !it.Is(token.From) {
				return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
			}

		case token.Literal:

			column := strings.TrimSpace(e.Value)
			alias := ""

			if it.Is(token.As) {
				e = it.Next()
				if !it.Is(token.Literal) {
					return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
				}
				e = it.Next()
				alias = strings.TrimSpace(e.Value)
			}

			if it.Is(token.Comma) {
				it.Next()
			}

			query.Columns = append(query.Columns, stmt.NewColumnAlias(column, alias))

			if it.HasNext() && !it.Is(token.From) && !it.Is(token.Literal) {
				return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
			}

		// Parse from
		case token.From:
			if !query.From.IsEmpty() {
				return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
			}

			from, err := parseFrom(it)
			if err != nil {
				return stmt.Select{}, err
			}
			query.From = from

		// Parse joins
		case token.Inner:
			if !it.Is(token.Join) {
				return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
			}

			join, err := parseJoin(it, stmt.Join{
				Type: types.InnerJoin,
			})
			if err != nil {
				return stmt.Select{}, err
			}

			query.Joins = append(query.Joins, join)

		case token.Left:
			if !it.Is(token.Join) {
				return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
			}

			join, err := parseJoin(it, stmt.Join{
				Type: types.LeftJoin,
			})
			if err != nil {
				return stmt.Select{}, err
			}

			query.Joins = append(query.Joins, join)

		case token.Right:
			if !it.Is(token.Join) {
				return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
			}

			join, err := parseJoin(it, stmt.Join{
				Type: types.RightJoin,
			})
			if err != nil {
				return stmt.Select{}, err
			}

			query.Joins = append(query.Joins, join)

		case token.Where:
			if !query.Where.IsEmpty() {
				return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
			}

			where, err := parseWhere(it)
			if err != nil {
				return stmt.Select{}, err
			}
			query.Where = where

		}
	}

	return query, nil
}
