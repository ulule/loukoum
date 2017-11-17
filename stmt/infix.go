package stmt

import (
	"github.com/ulule/loukoum/types"
)

type InfixExpression struct {
	Expression
	Left     Expression
	Operator Operator
	Right    Expression
}

func NewInfixExpression(left Expression, operator Operator, right Expression) InfixExpression {
	return InfixExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (expression InfixExpression) Write(ctx *types.Context) {
	if expression.IsEmpty() {
		panic("loukoum: expression is undefined")
	}

	ctx.Write("(")
	expression.Left.Write(ctx)
	ctx.Write(" ")
	expression.Operator.Write(ctx)
	ctx.Write(" ")
	expression.Right.Write(ctx)
	ctx.Write(")")
}

// IsEmpty return true if statement is undefined.
func (expression InfixExpression) IsEmpty() bool {
	return expression.Left == nil || expression.Operator == nil || expression.Right == nil ||
		expression.Left.IsEmpty() || expression.Operator.IsEmpty() || expression.Right.IsEmpty()
}

// And creates a new InfixExpression using given Expression.
func (expression InfixExpression) And(value Expression) InfixExpression {
	operator := NewAndOperator()
	return NewInfixExpression(expression, operator, value)
}

// Or creates a new InfixExpression using given Expression.
func (expression InfixExpression) Or(value Expression) InfixExpression {
	operator := NewOrOperator()
	return NewInfixExpression(expression, operator, value)
}
