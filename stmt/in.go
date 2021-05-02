package stmt

import (
	"github.com/ulule/loukoum/v3/types"
)

// In is a IN expression.
type In struct {
	Expression Expression
	Operator   ComparisonOperator
	Value      Expression
}

// NewIn returns a new In instance using an inclusive operator.
func NewIn(expression Expression, value Expression) In {
	return In{
		Expression: expression,
		Operator:   NewComparisonOperator(types.In),
		Value:      value,
	}
}

// NewNotIn returns a new In instance using an exclusive operator.
func NewNotIn(expression Expression, value Expression) In {
	return In{
		Expression: expression,
		Operator:   NewComparisonOperator(types.NotIn),
		Value:      value,
	}
}

func (In) expression() {}

// Write exposes statement as a SQL query.
func (in In) Write(ctx types.Context) {
	if in.IsEmpty() {
		panic("loukoum: expression is undefined")
	}

	ctx.Write("(")
	in.Expression.Write(ctx)
	ctx.Write(" ")
	in.Operator.Write(ctx)
	ctx.Write(" (")
	if !in.Value.IsEmpty() {
		in.Value.Write(ctx)
	}
	ctx.Write("))")
}

// IsEmpty returns true if statement is undefined.
func (in In) IsEmpty() bool {
	return in.Expression.IsEmpty() || in.Operator.IsEmpty() || in.Value == nil
}

// And creates a new InfixExpression using given Expression.
func (in In) And(value Expression) InfixExpression {
	operator := NewAndOperator()
	return NewInfixExpression(in, operator, value)
}

// Or creates a new InfixExpression using given Expression.
func (in In) Or(value Expression) InfixExpression {
	operator := NewOrOperator()
	return NewInfixExpression(in, operator, value)
}

// Ensure that In is an Expression
var _ Expression = In{}
