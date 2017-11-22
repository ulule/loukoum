package stmt

import (
	"github.com/ulule/loukoum/types"
)

type Join struct {
	Statement
	Type      types.JoinType
	Table     Table
	Condition On
}

func NewJoin(kind types.JoinType, table Table, condition On) Join {
	return Join{
		Type:      kind,
		Table:     table,
		Condition: condition,
	}
}

func NewInnerJoin(table Table, condition On) Join {
	return NewJoin(types.InnerJoin, table, condition)
}

func NewLeftJoin(table Table, condition On) Join {
	return NewJoin(types.LeftJoin, table, condition)
}

func NewRightJoin(table Table, condition On) Join {
	return NewJoin(types.RightJoin, table, condition)
}

// Write expose statement as a SQL query.
func (join Join) Write(ctx *types.Context) {
	ctx.Write(join.Type.String())
	ctx.Write(" ")
	join.Table.Write(ctx)
	ctx.Write(" ")
	join.Condition.Write(ctx)
}

// IsEmpty return true if statement is undefined.
func (join Join) IsEmpty() bool {
	return join.Type == "" || join.Table.IsEmpty() || join.Condition.IsEmpty()
}
