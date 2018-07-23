package parser

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ulule/loukoum/lexer"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/token"
)

// Parse will try to parse given query as a statement.
func Parse(query string) (stmt.Statement, error) { // nolint: gocyclo
	lexer := lexer.New(strings.NewReader(query))
	it := lexer.Iterator()

	if it.Is(token.Select) {
		st, err := parseSelect(it)
		if err != nil {
			return nil, errors.Wrapf(err, "given query cannot be parsed: %s", query)
		}
		return st, nil
	}

	for it.HasNext() {
		e := it.Next()
		fmt.Println(e)
	}

	return nil, nil
}
