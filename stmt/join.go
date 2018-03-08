package stmt

import (
	"github.com/ulule/loukoum/types"
)

// Join is a JOIN clause.
type Join struct {
	Type      types.JoinType
	Table     Table
	Condition On
}

// NewJoin returns a new Join instance.
func NewJoin(kind types.JoinType, table Table, condition On) Join {
	return Join{
		Type:      kind,
		Table:     table,
		Condition: condition,
	}
}

// NewInnerJoin returns a new Join instance using an INNER JOIN.
func NewInnerJoin(table Table, condition On) Join {
	return NewJoin(types.InnerJoin, table, condition)
}

// NewLeftJoin returns a new Join instance using a LEFT JOIN.
func NewLeftJoin(table Table, condition On) Join {
	return NewJoin(types.LeftJoin, table, condition)
}

// NewRightJoin returns a new Join instance using a RIGHT JOIN.
func NewRightJoin(table Table, condition On) Join {
	return NewJoin(types.RightJoin, table, condition)
}

// Write exposes statement as a SQL query.
func (join Join) Write(ctx types.Context) {
	ctx.Write(join.Type.String())
	ctx.Write(" ")
	join.Table.Write(ctx)
	ctx.Write(" ")
	join.Condition.Write(ctx)
}

// IsEmpty returns true if statement is undefined.
func (join Join) IsEmpty() bool {
	return join.Type == "" || join.Table.IsEmpty() || join.Condition.IsEmpty()
}

// Ensure that Join is a Statement
var _ Statement = Join{}
