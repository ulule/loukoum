package builder

import (
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

// Update is a builder used for "UPDATE" query.
type Update struct {
	query stmt.Update
}

// NewUpdate creates a new Update.
func NewUpdate(arg interface{}) Update {
	return Update{
		query: stmt.NewUpdate(ToTable(arg)),
	}
}

// Set adds a SET clause.
func (b Update) Set(m types.Map) Update {
	if !b.query.Set.IsEmpty() {
		panic("loukoum: update builder has set clause already defined")
	}

	b.query.Set = ToSet(m)

	return b
}

// Where adds WHERE clauses.
func (b Update) Where(condition stmt.Expression) Update {
	if b.query.Where.IsEmpty() {
		b.query.Where = stmt.NewWhere(condition)
		return b
	}

	return b.And(condition)
}

// And adds AND WHERE conditions.
func (b Update) And(condition stmt.Expression) Update {
	b.query.Where = b.query.Where.And(condition)
	return b
}

// Or adds OR WHERE conditions.
func (b Update) Or(condition stmt.Expression) Update {
	b.query.Where = b.query.Where.Or(condition)
	return b
}

// From sets the FROM clause of the query.
func (b Update) From(arg interface{}) Update {
	if !b.query.From.IsEmpty() {
		panic("loukoum: update builder has from clause already defined")
	}

	b.query.From = ToFrom(arg)

	return b
}

// Returning adds a RETURNING clause.
func (b Update) Returning(values ...interface{}) Update {
	if !b.query.Returning.IsEmpty() {
		panic("loukoum: update builder has returning clause already defined")
	}

	b.query.Returning = stmt.NewReturning(ToColumns(values))

	return b
}

// String returns the underlying query as a raw statement.
func (b Update) String() string {
	return rawify(b.Prepare())
}

// Prepare returns the underlying query as a named statement.
func (b Update) Prepare() (string, map[string]interface{}) {
	ctx := types.NewContext()
	b.query.Write(ctx)
	return ctx.Query(), ctx.Values()
}

// Statement returns underlying statement.
func (b Update) Statement() stmt.Statement {
	return b.query
}

// Ensure that Update is a Builder
var _ Builder = Update{}
