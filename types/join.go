package types

// JoinType represents a join type.
type JoinType string

func (e JoinType) String() string {
	return string(e)
}

// Join types.
const (
	// InnerJoin has a "INNER JOIN" type.
	InnerJoin = JoinType("INNER JOIN")
	// LeftJoin has a "LEFT JOIN" type.
	LeftJoin = JoinType("LEFT JOIN")
	// RightJoin has a "RIGHT JOIN" type.
	RightJoin = JoinType("RIGHT JOIN")
)
