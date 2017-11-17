package loukoum

import (
	"fmt"

	"github.com/ulule/loukoum/stmt"
)

// Select start a SelectBuilder using given columns.
func Select(columns ...interface{}) SelectBuilder {
	return NewSelectBuilder().Columns(columns)
}

// SelectDistinct start a SelectBuilder using given columns and "DISTINCT" option.
func SelectDistinct(columns ...interface{}) SelectBuilder {
	return Select(columns...).Distinct()
}

// Insert starts an InsertBuilder using the given table as into clause.
func Insert(into string) InsertBuilder {
	return NewInsertBuilder().Into(into)
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
