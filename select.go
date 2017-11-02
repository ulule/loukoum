package loukoum

import (
	"bytes"
	"fmt"

	"github.com/ulule/loukoum/stmt"
)

type SelectBuilder struct {
	columns  []stmt.Column
	distinct bool
}

func NewSelectBuilder() SelectBuilder {
	return SelectBuilder{}
}

func (builder SelectBuilder) Columns(columns []interface{}) SelectBuilder {
	if len(builder.columns) != 0 {
		panic("loukoum: select builder has columns already defined")
	}
	if len(columns) == 0 {
		columns = []interface{}{"*"}
	}

	builder.columns = []stmt.Column{}
	for i := range columns {
		switch column := columns[i].(type) {
		case string:
			builder.columns = append(builder.columns, stmt.NewColumn(column))
		case stmt.Column:
			builder.columns = append(builder.columns, column)
		default:
			panic(fmt.Sprintf("loukoum: cannot use %T as column", column))
		}
	}

	return builder
}

func (builder SelectBuilder) Distinct() SelectBuilder {
	builder.distinct = true
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

	// TODO FROM

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
