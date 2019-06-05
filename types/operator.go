package types

// LogicalOperator represents a logical operator.
type LogicalOperator string

func (e LogicalOperator) String() string {
	return string(e)
}

// Logical operators.
const (
	And = LogicalOperator("AND")
	Or  = LogicalOperator("OR")
	Not = LogicalOperator("NOT")
)

// ComparisonOperator represents a comparison operator.
type ComparisonOperator string

func (e ComparisonOperator) String() string {
	return string(e)
}

// Comparison operators.
const (
	Equal              = ComparisonOperator("=")
	NotEqual           = ComparisonOperator("!=")
	Is                 = ComparisonOperator("IS")
	IsNot              = ComparisonOperator("IS NOT")
	GreaterThan        = ComparisonOperator(">")
	GreaterThanOrEqual = ComparisonOperator(">=")
	LessThan           = ComparisonOperator("<")
	LessThanOrEqual    = ComparisonOperator("<=")
	In                 = ComparisonOperator("IN")
	NotIn              = ComparisonOperator("NOT IN")
	Like               = ComparisonOperator("LIKE")
	NotLike            = ComparisonOperator("NOT LIKE")
	ILike              = ComparisonOperator("ILIKE")
	NotILike           = ComparisonOperator("NOT ILIKE")
	Between            = ComparisonOperator("BETWEEN")
	NotBetween         = ComparisonOperator("NOT BETWEEN")
	Contains           = ComparisonOperator("@>")
	IsContainedBy      = ComparisonOperator("<@")
	Overlap            = ComparisonOperator("&&")
)
