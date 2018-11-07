package stmt

import (
	"github.com/ulule/loukoum/v2/token"
	"github.com/ulule/loukoum/v2/types"
)

// Where is a WHERE clause.
type Where struct {
	Condition Expression
}

// NewWhere returns a new Where instance.
func NewWhere(expression Expression) Where {
	return Where{
		Condition: NewWrapper(expression),
	}
}

// Write exposes statement as a SQL query.
func (where Where) Write(ctx types.Context) {
	if where.IsEmpty() {
		panic("loukoum: a where clause expects at least one condition")
	}

	ctx.Write(token.Where.String())
	ctx.Write(" ")
	where.Condition.Write(ctx)
}

// IsEmpty returns true if statement is undefined.
func (where Where) IsEmpty() bool {
	return where.Condition == nil || where.Condition.IsEmpty()
}

// And appends given Expression using AND as logical operator.
func (where Where) And(right Expression) Where {
	if where.IsEmpty() {
		panic("loukoum: two conditions are required for AND statement")
	}

	left := where.Condition
	operator := NewAndOperator()
	where.Condition = NewInfixExpression(left, operator, NewWrapper(right))
	return where
}

// Or appends given Expression using OR as logical operator.
func (where Where) Or(right Expression) Where {
	if where.IsEmpty() {
		panic("loukoum: two conditions are required for OR statement")
	}

	left := where.Condition
	operator := NewOrOperator()
	where.Condition = NewInfixExpression(left, operator, NewWrapper(right))
	return where
}

// Ensure that Where is a Statement
var _ Statement = Where{}
