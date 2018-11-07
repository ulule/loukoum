package stmt

import (
	"github.com/ulule/loukoum/v2/token"
	"github.com/ulule/loukoum/v2/types"
)

// From is a FROM clause.
type From struct {
	Only  bool
	Table Table
}

// NewFrom returns a new From instance.
func NewFrom(table Table, only bool) From {
	return From{
		Only:  only,
		Table: table,
	}
}

// Write exposes statement as a SQL query.
func (from From) Write(ctx types.Context) {
	ctx.Write(token.From.String())
	if from.Only {
		ctx.Write(" ")
		ctx.Write(token.Only.String())
	}
	ctx.Write(" ")
	from.Table.Write(ctx)
}

// IsEmpty returns true if statement is undefined.
func (from From) IsEmpty() bool {
	return from.Table.IsEmpty()
}

// Ensure that From is a Statement
var _ Statement = From{}
