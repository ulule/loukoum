package stmt

import (
	"github.com/ulule/loukoum/types"
)

type GroupBy struct {
	Statement
	Columns []Column
}

func NewGroupBy(columns []Column) GroupBy {
	return GroupBy{
		Columns: columns,
	}
}

func (group GroupBy) Write(ctx *types.Context) {
	ctx.Write("GROUP BY ")
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
