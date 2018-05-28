package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// On is a ON clause.
type On struct {
	Left  Column
	Right Column
}

// NewOn returns a new On instance.
func NewOn(left, right Column) On {
	return On{
		Left:  left,
		Right: right,
	}
}

func (On) expression() {}

// And creates a new InfixExpression using given On statement.
func (on On) And(value On) InfixExpression {
	operator := NewAndOperator()
	return NewInfixExpression(on, operator, value)
}

// Or creates a new InfixExpression using given On statement.
func (on On) Or(value On) InfixExpression {
	operator := NewOrOperator()
	return NewInfixExpression(on, operator, value)
}

// Write exposes statement as a SQL query.
func (on On) Write(ctx types.Context) {
	ctx.Write(on.Left.Name)
	ctx.Write(" ")
	ctx.Write(token.Equals.String())
	ctx.Write(" ")
	ctx.Write(on.Right.Name)
}

// IsEmpty returns true if statement is undefined.
func (on On) IsEmpty() bool {
	return on.Left.IsEmpty() || on.Right.IsEmpty()
}

// Ensure that On is an Expression
var _ Expression = On{}
