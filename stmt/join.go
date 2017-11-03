package stmt

import (
	"bytes"

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

func (join Join) Write(buffer *bytes.Buffer) {
	buffer.WriteString(join.Type.String())
	buffer.WriteString(" ")
	join.Table.Write(buffer)
	buffer.WriteString(" ")
	join.Condition.Write(buffer)
}

// IsEmpty return true if statement is undefined.
func (join Join) IsEmpty() bool {
	return join.Type == "" || join.Table.IsEmpty() || join.Condition.IsEmpty()
}
