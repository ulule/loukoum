package stmt

import (
	"github.com/ulule/loukoum/v3/token"
	"github.com/ulule/loukoum/v3/types"
)

// DistinctOn is a distinct on expression.
type DistinctOn struct {
	Columns []Column
}

// NewDistinctOn returns a new DistinctOn instance.
func NewDistinctOn(columns []Column) DistinctOn {
	return DistinctOn{
		Columns: columns,
	}
}

// Write exposes statement as a SQL query.
func (distinctOn DistinctOn) Write(ctx types.Context) {
	if distinctOn.IsEmpty() {
		return
	}
	ctx.Write(token.DistinctOn.String())
	ctx.Write(" (")
	for i := range distinctOn.Columns {
		if i != 0 {
			ctx.Write(", ")
		}
		distinctOn.Columns[i].Write(ctx)
	}
	ctx.Write(")")
}

// IsEmpty returns true if statement is undefined.
func (distinctOn DistinctOn) IsEmpty() bool {
	return len(distinctOn.Columns) == 0
}

// Ensure that DistinctOn is a Statement
var _ Statement = DistinctOn{}
