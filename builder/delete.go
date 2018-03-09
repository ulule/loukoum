package builder

import (
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

// Delete is a builder used for "SELECT" query.
type Delete struct {
	query stmt.Delete
}

// NewDelete creates a new Delete.
func NewDelete() Delete {
	return Delete{}
}

// From sets the FROM clause of the query.
func (b Delete) From(arg interface{}) Delete {
	if !b.query.From.IsEmpty() {
		panic("loukoum: delete builder has from clause already defined")
	}

	only := b.query.From.Only
	b.query.From = ToFrom(arg)
	if only {
		b.query.From.Only = only
	}

	return b
}

// Only adds a ONLY clause to the query.
func (b Delete) Only() Delete {
	b.query.From.Only = true
	return b
}

// Using adds a ONLY clause to the query.
func (b Delete) Using(args ...interface{}) Delete {
	if !b.query.Using.IsEmpty() {
		panic("loukoum: delete builder has using clause already defined")
	}

	tables := ToTables(args)
	b.query.Using = stmt.NewUsing(tables)

	return b
}

// Where adds WHERE clauses.
func (b Delete) Where(condition stmt.Expression) Delete {
	if b.query.Where.IsEmpty() {
		b.query.Where = stmt.NewWhere(condition)
		return b
	}

	return b.And(condition)
}

// And adds AND WHERE conditions.
func (b Delete) And(condition stmt.Expression) Delete {
	b.query.Where = b.query.Where.And(condition)
	return b
}

// Or adds OR WHERE conditions.
func (b Delete) Or(condition stmt.Expression) Delete {
	b.query.Where = b.query.Where.Or(condition)
	return b
}

// Returning adds a RETURNING clause.
func (b Delete) Returning(values ...interface{}) Delete {
	if !b.query.Returning.IsEmpty() {
		panic("loukoum: delete builder has returning clause already defined")
	}

	b.query.Returning = stmt.NewReturning(ToColumns(values))

	return b
}

// String returns the underlying query as a raw statement.
func (b Delete) String() string {
	var ctx types.RawContext
	b.query.Write(&ctx)
	return ctx.Query()
}

// NamedQuery returns the underlying query as a named statement.
func (b Delete) NamedQuery() (string, map[string]interface{}) {
	var ctx types.NamedContext
	b.query.Write(&ctx)
	return ctx.Query(), ctx.Values()
}

// Query returns the underlying query as a regular statement.
func (b Delete) Query() (string, []interface{}) {
	var ctx types.StdContext
	b.query.Write(&ctx)
	return ctx.Query(), ctx.Values()
}

// Statement returns underlying statement.
func (b Delete) Statement() stmt.Statement {
	return b.query
}

// Ensure that Delete is a Builder
var _ Builder = Delete{}
