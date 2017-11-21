package types

// JoinType represents a join type.
type JoinType string

func (e JoinType) String() string {
	return string(e)
}

// Join types.
const (
	InnerJoin = JoinType("INNER JOIN")
	LeftJoin  = JoinType("LEFT JOIN")
	RightJoin = JoinType("RIGHT JOIN")
)
