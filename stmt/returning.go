package stmt

import (
	"sort"

	"github.com/ulule/loukoum/v3/token"
	"github.com/ulule/loukoum/v3/types"
)

// Returning is a RETURNING clause.
type Returning struct {
	Columns []Column
}

// NewReturning returns a new Returning instance.
func NewReturning(columns []Column) Returning {
	return Returning{
		Columns: columns,
	}
}

// Write exposes statement as a SQL query.
func (returning Returning) Write(ctx types.Context) {
	ctx.Write(token.Returning.String())
	ctx.Write(" ")

	sort.Slice(returning.Columns, func(i, j int) bool {
		return returning.Columns[i].Name < returning.Columns[j].Name ||
			returning.Columns[i].Alias < returning.Columns[j].Alias
	})

	for i := range returning.Columns {
		if i > 0 {
			ctx.Write(", ")
		}
		returning.Columns[i].Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (returning Returning) IsEmpty() bool {
	return len(returning.Columns) == 0
}

// Ensure that Returning is a Statement
var _ Statement = Returning{}
