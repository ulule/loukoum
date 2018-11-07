package stmt

import (
	"github.com/ulule/loukoum/v2/token"
	"github.com/ulule/loukoum/v2/types"
)

// Having is a HAVING clause.
type Having struct {
	Statement
	Condition Expression
}

// NewHaving returns a new Having instance.
func NewHaving(expression Expression) Having {
	return Having{
		Condition: expression,
	}
}

// Write exposes statement as a SQL query.
func (having Having) Write(ctx types.Context) {
	if having.IsEmpty() {
		panic("loukoum: a having clause expects at least one condition")
	}

	ctx.Write(token.Having.String())
	ctx.Write(" ")
	having.Condition.Write(ctx)
}

// IsEmpty returns true if statement is undefined.
func (having Having) IsEmpty() bool {
	return having.Condition == nil || having.Condition.IsEmpty()
}

// And appends given Expression using AND as logical operator.
func (having Having) And(right Expression) Having {
	if having.IsEmpty() {
		panic("loukoum: two conditions are required for AND statement")
	}

	left := having.Condition
	operator := NewAndOperator()
	having.Condition = NewInfixExpression(left, operator, right)
	return having
}

// Or appends given Expression using OR as logical operator.
func (having Having) Or(right Expression) Having {
	if having.IsEmpty() {
		panic("loukoum: two conditions are required for OR statement")
	}

	left := having.Condition
	operator := NewOrOperator()
	having.Condition = NewInfixExpression(left, operator, right)
	return having
}

// Ensure that Having is a Statement
var _ Statement = Having{}
