package types

type JoinType string

func (e JoinType) String() string {
	return string(e)
}

const (
	InnerJoin = JoinType("INNER JOIN")
	LeftJoin  = JoinType("LEFT JOIN")
	RightJoin = JoinType("RIGHT JOIN")
)
