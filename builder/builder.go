package builder

import (
	"fmt"
	"strings"

	"github.com/ulule/loukoum/stmt"
)

// Builder defines a generic methods available for Select, Insert, Update and Delete builders.
type Builder interface {
	// String returns the underlying query as a raw statement.
	String()
	// Prepare returns the underlying query as a named statement.
	Prepare() (string, map[string]interface{})
	// Statement return underlying statement.
	Statement() stmt.Statement
}

func rawify(query string, args map[string]interface{}) string {
	for key, arg := range args {
		switch value := arg.(type) {
		case string:
			query = strings.Replace(query, key, fmt.Sprint("'", value, "'"), 1)
		default:
			query = strings.Replace(query, key, fmt.Sprint(value), 1)
		}
	}
	return query
}

// ToColumns takes a list of empty interfaces and returns a slice of Column instance.
func ToColumns(values []interface{}) []stmt.Column {
	columns := make([]stmt.Column, 0, len(values))

	for i := range values {
		column := stmt.Column{}

		switch value := values[i].(type) {
		case string:
			column = stmt.NewColumn(value)
		case stmt.Column:
			column = value
		default:
			panic(fmt.Sprintf("loukoum: cannot use %T as column", column))
		}
		if column.IsEmpty() {
			panic("loukoum: a column was undefined")
		}

		columns = append(columns, column)
	}

	return columns
}
