package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// Returning is the RETURNING clause.
type Returning struct {
	Statement
	Columns []Column
}

// NewReturning returns a new Returning instance.
func NewReturning(columns []Column) Returning {
	return Returning{
		Columns: columns,
	}
}

// Write implements Statement interface.
func (returning Returning) Write(ctx *types.Context) {
	ctx.Write(token.Returning.String())
	ctx.Write(" ")

	l := len(returning.Columns)
	if l > 1 {
		ctx.Write("(")
	}

	for i := range returning.Columns {
		if i > 0 {
			ctx.Write(", ")
		}
		returning.Columns[i].Write(ctx)
	}

	if l > 1 {
		ctx.Write(")")
	}
}

// IsEmpty implements Statement interface.
func (returning Returning) IsEmpty() bool {
	return len(returning.Columns) == 0
}
