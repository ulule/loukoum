package stmt

import (
	"github.com/ulule/loukoum/types"
)

// In is a IN expression.
type In struct {
	Identifier Identifier
	Operator   ComparisonOperator
	Value      Expression
}

// NewIn returns a new In instance using an inclusive operator.
func NewIn(identifier Identifier, value Expression) In {
	return In{
		Identifier: identifier,
		Operator:   NewComparisonOperator(types.In),
		Value:      value,
	}
}

// NewNotIn returns a new In instance using an exclusive operator.
func NewNotIn(identifier Identifier, value Expression) In {
	return In{
		Identifier: identifier,
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
	in.Identifier.Write(ctx)
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
	return in.Identifier.IsEmpty() || in.Operator.IsEmpty() || in.Value == nil
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
