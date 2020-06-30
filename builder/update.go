package builder

import (
	"github.com/ulule/loukoum/v3/stmt"
	"github.com/ulule/loukoum/v3/types"
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

// With adds WITH clauses.
func (b Update) With(args ...stmt.WithQuery) Update {
	if b.query.With.IsEmpty() {
		b.query.With = stmt.NewWith(args)
		return b
	}

	b.query.With.Queries = append(b.query.With.Queries, args...)
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

// Comment adds comment to the query.
func (b Update) Comment(comment string) Update {
	b.query.Comment = stmt.NewComment(comment)

	return b
}

// String returns the underlying query as a raw statement.
// This function should be used for debugging since it doesn't escape anything and is completely
// vulnerable to SQL injection.
// You should use either NamedQuery() or Query()...
func (b Update) String() string {
	ctx := &types.RawContext{}
	b.query.Write(ctx)
	return ctx.Query()
}

// NamedQuery returns the underlying query as a named statement.
func (b Update) NamedQuery() (string, map[string]interface{}) {
	ctx := &types.NamedContext{}
	b.query.Write(ctx)
	return ctx.Query(), ctx.Values()
}

// Query returns the underlying query as a regular statement.
func (b Update) Query() (string, []interface{}) {
	ctx := &types.StdContext{}
	return ctx.Query(), ctx.Values()
}

func (b Update) Write(ctx types.Context) {
	b.query.Write(ctx)
}

// Statement returns underlying statement.
func (b Update) Statement() stmt.Statement {
	return b.query
}

// Ensure that Update is a Builder
var _ Builder = Update{}
