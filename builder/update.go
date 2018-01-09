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

// Only sets the ONLY clause.
func (b Update) Only() Update {
	b.query.Only = true
	return b
}

// Set adds a SET clause.
func (b Update) Set(args ...interface{}) Update {
	if len(args) == 0 {
		panic("loukoum: update set clause requires at least one argument")
	}

	b.query.Set = MergeSet(b.query.Set, args)
	return b
}

// Using assigns the result of the given expression to
// the columns defined in Set.
func (b Update) Using(args ...interface{}) Update {
	if b.query.Set.Pairs.Mode != stmt.PairArrayMode {
		panic("loukoum: you can only use Using with column-list syntax")
	}

	if len(args) == 0 {
		panic("loukoum: using clause requires a column or an expression")
	}

	for i := range args {
		b.query.Set.Pairs.Use(stmt.NewExpression(args[i]))
	}

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
