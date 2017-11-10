package stmt

import (
	"bytes"

	"github.com/ulule/loukoum/types"
)

type Operator interface {
	Statement
	operator()
}

type LogicalOperator struct {
	Statement
	Operator types.LogicalOperator
}

func NewAndOperator() LogicalOperator {
	return NewLogicalOperator(types.And)
}

func NewOrOperator() LogicalOperator {
	return NewLogicalOperator(types.Or)
}

func NewLogicalOperator(operator types.LogicalOperator) LogicalOperator {
	return LogicalOperator{
		Operator: operator,
	}
}

func (LogicalOperator) operator() {}

func (operator LogicalOperator) Write(buffer *bytes.Buffer) {
	buffer.WriteString(operator.Operator.String())
}

// IsEmpty return true if statement is undefined.
func (operator LogicalOperator) IsEmpty() bool {
	return operator.Operator == ""
}

type ComparisonOperator struct {
	Statement
	Operator types.ComparisonOperator
}

func NewComparisonOperator(operator types.ComparisonOperator) ComparisonOperator {
	return ComparisonOperator{
		Operator: operator,
	}
}

func (ComparisonOperator) operator() {}

func (operator ComparisonOperator) Write(buffer *bytes.Buffer) {
	buffer.WriteString(operator.Operator.String())
}

// IsEmpty return true if statement is undefined.
func (operator ComparisonOperator) IsEmpty() bool {
	return operator.Operator == ""
}
