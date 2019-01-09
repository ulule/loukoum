package stmt

import (
	"github.com/ulule/loukoum/v3/token"
	"github.com/ulule/loukoum/v3/types"
)

// Using is a USING clause.
type Using struct {
	Tables []Table
}

// NewUsing returns a new Using instance.
func NewUsing(tables []Table) Using {
	return Using{
		Tables: tables,
	}
}

// Write exposes statement as a SQL query.
func (using Using) Write(ctx types.Context) {
	if using.IsEmpty() {
		return
	}

	ctx.Write(token.Using.String())
	ctx.Write(" ")

	for i := range using.Tables {
		if i != 0 {
			ctx.Write(", ")
		}
		using.Tables[i].Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (using Using) IsEmpty() bool {
	return len(using.Tables) == 0
}

// Ensure that Using is a Statement
var _ Statement = Using{}
