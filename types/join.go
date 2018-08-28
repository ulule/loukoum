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
	// LeftOuterJoin has a "LEFT OUTER JOIN" type.
	LeftOuterJoin = JoinType("LEFT OUTER JOIN")
	// RightOuterJoin has a "RIGHT OUTER JOIN" type.
	RightOuterJoin = JoinType("RIGHT OUTER JOIN")
)
