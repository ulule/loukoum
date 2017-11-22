package stmt

import (
	"github.com/ulule/loukoum/types"
)

type From struct {
	Statement
	Table Table
}

func NewFrom(table Table) From {
	return From{
		Table: table,
	}
}

// Write expose statement as a SQL query.
func (from From) Write(ctx *types.Context) {
	ctx.Write("FROM ")
	from.Table.Write(ctx)
}

// IsEmpty return true if statement is undefined.
func (from From) IsEmpty() bool {
	return from.Table.IsEmpty()
}
