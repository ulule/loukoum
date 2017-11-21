package builder

import (
	"fmt"

	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

// Insert is a builder used for "INSERT" query.
type Insert struct {
	insert stmt.Insert
}

// NewInsert creates a new Insert.
func NewInsert() Insert {
	return Insert{
		insert: stmt.NewInsert(),
	}
}

func (b Insert) String() string {
	ctx := types.NewContext()
	b.insert.Write(ctx)
	return ctx.Query()
}

// Into sets the INTO clause of the query.
func (b Insert) Into(into interface{}) Insert {
	if !b.insert.Into.IsEmpty() {
		panic("loukoum: insert builder has into clause already defined")
	}

	switch value := into.(type) {
	case string:
		b.insert.Into = stmt.NewInto(stmt.NewTable(value))
	case stmt.Into:
		b.insert.Into = value
	case stmt.Table:
		b.insert.Into = stmt.NewInto(value)
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as into clause", into))
	}

	if b.insert.Into.IsEmpty() {
		panic("loukoum: the given into clause is undefined")
	}

	return b
}

// Columns sets the insert columns.
func (b Insert) Columns(columns ...interface{}) Insert {
	if len(columns) == 0 {
		return b
	}
	if len(b.insert.Columns) != 0 {
		panic("loukoum: insert builder has columns clause already defined")
	}

	b.insert.Columns = ToColumns(columns)

	return b
}

// Values sets the INSERT values.
func (b Insert) Values(values interface{}) Insert {
	if !b.insert.Values.IsEmpty() {
		panic("loukoum: insert builder has values clause already defined")
	}

	b.insert.Values = stmt.NewValues(stmt.NewExpression(values))

	return b
}

// Returning builds the RETURNING clause.
func (b Insert) Returning(values ...interface{}) Insert {
	if !b.insert.Returning.IsEmpty() {
		panic("loukoum: insert builder has returning clause already defined")
	}

	b.insert.Returning = stmt.NewReturning(ToColumns(values))

	return b
}
