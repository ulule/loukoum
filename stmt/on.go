package stmt

import (
	"github.com/ulule/loukoum/v3/token"
	"github.com/ulule/loukoum/v3/types"
)

// OnExpression is a SQL expression for a ON statement.
type OnExpression interface {
	Statement
	And(value OnExpression) OnExpression
	Or(value OnExpression) OnExpression
	onExpression()
}

// OnClause is a ON clause.
type OnClause struct {
	Left  Column
	Right Column
}

// NewOnClause returns a new On instance.
func NewOnClause(left, right Column) OnClause {
	return OnClause{
		Left:  left,
		Right: right,
	}
}

func (OnClause) onExpression() {}

// And creates a new InfixOnExpression using given OnExpression.
func (on OnClause) And(value OnExpression) OnExpression {
	operator := NewAndOperator()
	return NewInfixOnExpression(on, operator, value)
}

// Or creates a new InfixOnExpression using given OnExpression.
func (on OnClause) Or(value OnExpression) OnExpression {
	operator := NewOrOperator()
	return NewInfixOnExpression(on, operator, value)
}

// Write exposes statement as a SQL query.
func (on OnClause) Write(ctx types.Context) {
	ctx.Write(on.Left.Name)
	ctx.Write(" ")
	ctx.Write(token.Equals.String())
	ctx.Write(" ")
	ctx.Write(on.Right.Name)
}

// IsEmpty returns true if statement is undefined.
func (on OnClause) IsEmpty() bool {
	return on.Left.IsEmpty() || on.Right.IsEmpty()
}

// Ensure that OnClause is an OnExpression
var _ OnExpression = OnClause{}

// InfixOnExpression is an OnExpression that has a left and right clauses with a logical operator
// for an ON statement.
type InfixOnExpression struct {
	Left     OnExpression
	Operator LogicalOperator
	Right    OnExpression
}

// NewInfixOnExpression returns a new InfixOnExpression instance.
func NewInfixOnExpression(left OnExpression, operator LogicalOperator, right OnExpression) InfixOnExpression {
	return InfixOnExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (InfixOnExpression) onExpression() {}

// Write exposes statement as a SQL query.
func (expression InfixOnExpression) Write(ctx types.Context) {
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

// IsEmpty returns true if statement is undefined.
func (expression InfixOnExpression) IsEmpty() bool {
	return expression.Left == nil || expression.Right == nil ||
		expression.Left.IsEmpty() || expression.Operator.IsEmpty() || expression.Right.IsEmpty()
}

// And creates a new InfixOnExpression using given OnExpression.
func (expression InfixOnExpression) And(value OnExpression) OnExpression {
	operator := NewAndOperator()
	return NewInfixOnExpression(expression, operator, value)
}

// Or creates a new InfixOnExpression using given OnExpression.
func (expression InfixOnExpression) Or(value OnExpression) OnExpression {
	operator := NewOrOperator()
	return NewInfixOnExpression(expression, operator, value)
}

// Ensure that InfixOnExpression is an OnExpression
var _ OnExpression = InfixOnExpression{}
