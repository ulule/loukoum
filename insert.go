package loukoum

import (
	"bytes"
	"fmt"

	"github.com/ulule/loukoum/stmt"
)

// InsertBuilder is a builder used for "INSERT" query.
type InsertBuilder struct {
	insert stmt.Insert
}

// NewInsertBuilder creates a new InsertBuilder.
func NewInsertBuilder() InsertBuilder {
	return InsertBuilder{
		insert: stmt.NewInsert(),
	}
}

func (builder InsertBuilder) String() string {
	buffer := &bytes.Buffer{}
	builder.insert.Write(buffer)
	return buffer.String()
}

// Into sets the INTO clause of the query.
func (builder InsertBuilder) Into(into interface{}) InsertBuilder {
	if !builder.insert.Into.IsEmpty() {
		panic("loukoum: insert builder has into clause already defined")
	}

	switch value := into.(type) {
	case string:
		builder.insert.Into = stmt.NewInto(stmt.NewTable(value))
	case stmt.Into:
		builder.insert.Into = value
	case stmt.Table:
		builder.insert.Into = stmt.NewInto(value)
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as into clause", into))
	}

	if builder.insert.Into.IsEmpty() {
		panic("loukoum: the given into clause is undefined")
	}

	return builder
}

// Columns sets the insert columns.
func (builder InsertBuilder) Columns(columns ...interface{}) InsertBuilder {
	if len(columns) == 0 {
		return builder
	}

	builder.insert.Columns = make([]stmt.Column, 0, len(columns))

	for i := range columns {
		column := stmt.Column{}
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
		builder.insert.Columns = append(builder.insert.Columns, column)
	}

	return builder
}

// Values sets the INSERT values.
func (builder InsertBuilder) Values(values interface{}) InsertBuilder {
	if !builder.insert.Values.IsEmpty() {
		return builder
	}

	builder.insert.Values = stmt.NewValues(stmt.NewExpression(values))

	return builder
}

// Returning builds the RETURNING clause.
func (builder InsertBuilder) Returning(values ...interface{}) InsertBuilder {
	if !builder.insert.Returning.IsEmpty() {
		return builder
	}

	builder.insert.Returning = stmt.NewReturning(stmt.ToColumns(values))

	return builder
}
