package builder

import (
	"fmt"

	"github.com/ulule/loukoum/stmt"
)

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
