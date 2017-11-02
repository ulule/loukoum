package loukoum

import (
	"bytes"
	"fmt"

	"github.com/ulule/loukoum/stmt"
)

// SelectBuilder is a builder used for "SELECT" query.
type SelectBuilder struct {
	distinct bool
	columns  []stmt.Column
	from     stmt.From
}

// NewSelectBuilder creates a new SelectBuilder.
func NewSelectBuilder() SelectBuilder {
	return SelectBuilder{}
}

// Distinct adds a DISTINCT clause to the query.
func (builder SelectBuilder) Distinct() SelectBuilder {
	builder.distinct = true
	return builder
}

// Columns adds result columns to the query.
func (builder SelectBuilder) Columns(columns []interface{}) SelectBuilder {
	if len(builder.columns) != 0 {
		panic("loukoum: select builder has columns already defined")
	}
	if len(columns) == 0 {
		columns = []interface{}{"*"}
	}

	builder.columns = []stmt.Column{}
	for i := range columns {
		var column stmt.Column
		switch value := columns[i].(type) {
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
		builder.columns = append(builder.columns, column)
	}

	return builder
}

// From sets the FROM clause of the query.
func (builder SelectBuilder) From(from interface{}) SelectBuilder {
	if !builder.from.IsEmpty() {
		panic("loukoum: select builder has from clause already defined")
	}

	switch value := from.(type) {
	case string:
		builder.from = stmt.NewFrom(value)
	case stmt.From:
		builder.from = value
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as from clause", from))
	}

	if builder.from.IsEmpty() {
		panic("loukoum: given from clause is undefined")
	}

	return builder
}

func (builder SelectBuilder) String() string {
	if len(builder.columns) == 0 {
		panic("loukoum: select statements must have at least one column")
	}

	buffer := &bytes.Buffer{}

	// TODO Add prefixes

	buffer.WriteString("SELECT ")

	if builder.distinct {
		buffer.WriteString("DISTINCT ")
	}

	for i := range builder.columns {
		if i != 0 {
			buffer.WriteString(", ")
		}
		builder.columns[i].Write(buffer)
	}

	if !builder.from.IsEmpty() {
		buffer.WriteString(" FROM ")
		builder.from.Write(buffer)
	}

	// TODO JOINS

	// TODO WHERE

	// TODO GROUP BY

	// TODO HAVING

	// TODO ORDER BY

	// TODO LIMIT

	// TODO OFFSET

	// TODO Add suffixes

	return buffer.String()
}
