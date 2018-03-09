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

// Columns sets the query columns.
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

// Values sets the query values.
func (b Insert) Values(values ...interface{}) Insert {
	if !b.insert.Values.IsEmpty() {
		panic("loukoum: insert builder has values clause already defined")
	}

	b.insert.Values = stmt.NewValues(stmt.NewArrayExpression(values...))

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

// OnConflict builds the ON CONFLICT clause.
func (b Insert) OnConflict(args ...interface{}) Insert {
	if !b.insert.OnConflict.IsEmpty() {
		panic("loukoum: insert builder has on conflict clause already defined")
	}

	if len(args) == 0 {
		panic("loukoum: on conflict clause requires arguments")
	}

	for i := range args {
		switch value := args[i].(type) {
		case string, stmt.Column:
			b.insert.OnConflict.Target.Columns = append(b.insert.OnConflict.Target.Columns, ToColumn(value))
		case stmt.ConflictNoAction:
			b.insert.OnConflict.Action = value
			return b
		case stmt.ConflictUpdateAction:
			if b.insert.OnConflict.Target.IsEmpty() {
				panic("loukoum: on conflict update clause requires at least one target")
			}
			b.insert.OnConflict.Action = value
			return b
		default:
			panic(fmt.Sprintf("loukoum: cannot use %T as on conflict clause", args[i]))
		}
	}

	panic("loukoum: on conflict clause requires an action")
}

// Set is a wrapper that defines columns and values clauses using a pair.
func (b Insert) Set(args ...interface{}) Insert {
	if len(b.insert.Columns) != 0 {
		panic("loukoum: insert builder has columns clause already defined")
	}
	if !b.insert.Values.IsEmpty() {
		panic("loukoum: insert builder has values clause already defined")
	}

	pairs := ToSet(args).Pairs
	columns, expressions := pairs.Values()

	array := stmt.NewArray()
	array.AddValues(expressions)

	values := stmt.NewValues(array)

	b.insert.Columns = columns
	b.insert.Values = values

	return b
}

// String returns the underlying query as a raw statement.
func (b Insert) String() string {
	var ctx types.RawContext
	b.insert.Write(&ctx)
	return ctx.Query()
}

// Prepare returns the underlying query as a named statement.
func (b Insert) Prepare() (string, map[string]interface{}) {
	var ctx types.NamedContext
	b.insert.Write(&ctx)
	return ctx.Query(), ctx.Values()
}

// Query returns the underlying query as a regular statement.
func (b Insert) Query() (string, []interface{}) {
	var ctx types.StdContext
	b.insert.Write(&ctx)
	return ctx.Query(), ctx.Values()
}

// Statement returns underlying statement.
func (b Insert) Statement() stmt.Statement {
	return b.insert
}

// Ensure that Insert is a Builder
var _ Builder = Insert{}
