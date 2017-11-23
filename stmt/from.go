package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// From is a FROM clause.
type From struct {
	Statement
	Table Table
}

// NewFrom returns a new From instance.
func NewFrom(table Table) From {
	return From{
		Table: table,
	}
}

// Write exposes statement as a SQL query.
func (from From) Write(ctx *types.Context) {
	ctx.Write(token.From.String())
	ctx.Write(" ")
	from.Table.Write(ctx)
}

// IsEmpty returns true if statement is undefined.
func (from From) IsEmpty() bool {
	return from.Table.IsEmpty()
}
