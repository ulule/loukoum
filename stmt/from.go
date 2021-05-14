package stmt

import (
	"github.com/ulule/loukoum/v3/token"
	"github.com/ulule/loukoum/v3/types"
)

// From is a FROM clause.
type From struct {
	Tables []Statement
}

// NewFrom returns a new From instance.
func NewFrom(tables []Statement) From {
	return From{
		Tables: tables,
	}
}

// Write exposes statement as a SQL query.
func (from From) Write(ctx types.Context) {
	ctx.Write(token.From.String())
	ctx.Write(" ")

	for i := range from.Tables {
		from.Tables[i].Write(ctx)

		if i != len(from.Tables)-1 {
			ctx.Write(", ")
		}
	}
}

// IsEmpty returns true if statement is undefined.
func (from From) IsEmpty() bool {
	return len(from.Tables) == 0
}

// Ensure that From is a Statement
var _ Statement = From{}
