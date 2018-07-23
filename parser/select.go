package parser

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ulule/loukoum/lexer"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/token"
)

// ErrSelectInvalidCondition is returned when select condition cannot be parsed.
var ErrSelectInvalidCondition = fmt.Errorf("select condition is invalid")

// Parse will try to parse given query as a statement.
func parseSelect(it *lexer.Iteratee) (stmt.Select, error) {
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

			if !it.Is(token.From) && !it.Is(token.Literal) {
				return stmt.Select{}, errors.WithStack(ErrSelectInvalidCondition)
			}

		// Parse from
		case token.From:
			from, err := parseFrom(it)
			if err != nil {
				return stmt.Select{}, err
			}
			query.From = from

		}
		// TODO Remove debug
		fmt.Println(e.Type)
		fmt.Println(e)
	}

	return query, nil
}

func parseFrom(it *lexer.Iteratee) (stmt.From, error) {
	if !it.Is(token.From) && !it.Is(token.Literal) && !it.Is(token.Only) {
		return stmt.From{}, errors.WithStack(ErrSelectInvalidCondition)
	}

	query := stmt.From{}
	for it.HasNext() {
		e := it.Next()
		switch e.Type {
		case token.From:
			continue

		case token.Only:

			query.Only = true

		case token.Literal:

			table := strings.TrimSpace(e.Value)
			alias := ""

			if it.Is(token.As) {
				e = it.Next()
				if !it.Is(token.Literal) {
					break
				}
				e = it.Next()
				alias = strings.TrimSpace(e.Value)
			}

			query.Table = stmt.NewTableAlias(table, alias)
			return query, nil

		default:
			break
		}
	}

	return stmt.From{}, errors.WithStack(ErrSelectInvalidCondition)
}
