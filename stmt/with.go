package stmt

import (
	"github.com/ulule/loukoum/v2/token"
	"github.com/ulule/loukoum/v2/types"
)

// With is a WITH clause.
type With struct {
	Queries []WithQuery
}

// NewWith returns a new With instance.
func NewWith(queries []WithQuery) With {
	return With{
		Queries: queries,
	}
}

// Write exposes statement as a SQL query.
func (with With) Write(ctx types.Context) {
	if with.IsEmpty() {
		return
	}
	ctx.Write(token.With.String())
	ctx.Write(" ")
	for i := range with.Queries {
		if i != 0 {
			ctx.Write(", ")
		}
		with.Queries[i].Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (with With) IsEmpty() bool {
	return len(with.Queries) == 0
}

// WithQuery is a statement in a With clause.
type WithQuery struct {
	Name     string
	Subquery Expression
}

// Write exposes statement as a SQL query.
func (with WithQuery) Write(ctx types.Context) {
	if with.IsEmpty() {
		return
	}
	ctx.Write(with.Name)
	ctx.Write(" ")
	ctx.Write(token.As.String())
	ctx.Write(" (")
	with.Subquery.Write(ctx)
	ctx.Write(")")
}

// IsEmpty returns true if statement is undefined.
func (with WithQuery) IsEmpty() bool {
	return with.Name == "" || with.Subquery == nil || (with.Subquery != nil && with.Subquery.IsEmpty())
}

// NewWithQuery returns a new WithQuery instance.
func NewWithQuery(name string, value interface{}) WithQuery {
	return WithQuery{
		Name:     name,
		Subquery: NewExpression(value),
	}
}
