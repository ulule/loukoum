package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// GroupBy is a GROUP BY clause.
type GroupBy struct {
	Statement
	Columns []Column
}

// NewGroupBy returns a new GroupBy instance.
func NewGroupBy(columns []Column) GroupBy {
	return GroupBy{
		Columns: columns,
	}
}

// Write expose statement as a SQL query.
func (group GroupBy) Write(ctx *types.Context) {
	ctx.Write(token.Group.String())
	ctx.Write(" ")
	ctx.Write(token.By.String())
	ctx.Write(" ")
	for i := range group.Columns {
		if i != 0 {
			ctx.Write(", ")
		}
		group.Columns[i].Write(ctx)
	}
}

// IsEmpty return true if statement is undefined.
func (group GroupBy) IsEmpty() bool {
	return len(group.Columns) == 0
}
