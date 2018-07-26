package parser

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ulule/loukoum/lexer"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/token"
)

// ErrFromInvalidCondition is returned when from condition cannot be parsed.
var ErrFromInvalidCondition = fmt.Errorf("from condition is invalid")

func parseFrom(it *lexer.Iteratee) (stmt.From, error) { // nolint: gocyclo
	if !it.Is(token.From) && !it.Is(token.Literal) && !it.Is(token.Only) {
		return stmt.From{}, errors.WithStack(ErrFromInvalidCondition)
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
				// We expect a litteral after to define table alias.
				if !it.Is(token.Literal) {
					return stmt.From{}, errors.WithStack(ErrFromInvalidCondition)
				}
			}
			if it.Is(token.Literal) {
				e = it.Next()
				alias = strings.TrimSpace(e.Value)
			}

			query.Table = stmt.NewTableAlias(table, alias)
			return query, nil

		default:
			return stmt.From{}, errors.WithStack(ErrFromInvalidCondition)
		}
	}

	return stmt.From{}, errors.WithStack(ErrFromInvalidCondition)
}
