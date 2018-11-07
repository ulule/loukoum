package stmt

import (
	"github.com/ulule/loukoum/v2/token"
	"github.com/ulule/loukoum/v2/types"
)

// Into is a INTO clause.
type Into struct {
	Table Table
}

// NewInto returns a new Into instance.
func NewInto(table Table) Into {
	return Into{
		Table: table,
	}
}

// Write exposes statement as a SQL query.
func (into Into) Write(ctx types.Context) {
	ctx.Write(token.Into.String())
	ctx.Write(" ")
	into.Table.Write(ctx)
}

// IsEmpty returns true if statement is undefined.
func (into Into) IsEmpty() bool {
	return into.Table.IsEmpty()
}

// Ensure that Into is a Statement
var _ Statement = Into{}
