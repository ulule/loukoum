package stmt

import (
	"github.com/ulule/loukoum/types"
)

// Operator are used to compose expressions.
type Operator interface {
	Statement
	operator()
}

// LogicalOperator are used to evaluate two expressions using a logical operator.
type LogicalOperator struct {
	Statement
	Operator types.LogicalOperator
}

// NewAndOperator returns a new AND LogicalOperator instance.
func NewAndOperator() LogicalOperator {
	return NewLogicalOperator(types.And)
}

// NewOrOperator returns a new OR LogicalOperator instance.
func NewOrOperator() LogicalOperator {
	return NewLogicalOperator(types.Or)
}

// NewLogicalOperator returns a new LogicalOperator instance.
func NewLogicalOperator(operator types.LogicalOperator) LogicalOperator {
	return LogicalOperator{
		Operator: operator,
	}
}

func (LogicalOperator) operator() {}

// Write expose statement as a SQL query.
func (operator LogicalOperator) Write(ctx *types.Context) {
	ctx.Write(operator.Operator.String())
}

// IsEmpty return true if statement is undefined.
func (operator LogicalOperator) IsEmpty() bool {
	return operator.Operator == ""
}

// ComparisonOperator are used to evaluate two expressions using a comparison operator.
type ComparisonOperator struct {
	Statement
	Operator types.ComparisonOperator
}

// NewComparisonOperator returns a new ComparisonOperator instance.
func NewComparisonOperator(operator types.ComparisonOperator) ComparisonOperator {
	return ComparisonOperator{
		Operator: operator,
	}
}

func (ComparisonOperator) operator() {}

// Write expose statement as a SQL query.
func (operator ComparisonOperator) Write(ctx *types.Context) {
	ctx.Write(operator.Operator.String())
}

// IsEmpty return true if statement is undefined.
func (operator ComparisonOperator) IsEmpty() bool {
	return operator.Operator == ""
}
