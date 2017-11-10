package stmt

import (
	"bytes"
)

// Where is a WHERE clause.
type Where struct {
	Statement
	Condition Expression
}

// NewWhere returns a new WHERE clause.
func NewWhere(expression Expression) Where {
	return Where{
		Condition: expression,
	}
}

// Write writes WHERE clause into the given buffer.
func (where Where) Write(buffer *bytes.Buffer) {
	if where.IsEmpty() {
		panic("loukoum: a where clause expects at least one condition")
	}

	buffer.WriteString("WHERE ")
	where.Condition.Write(buffer)
}

// IsEmpty return true if statement is undefined.
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
	where.Condition = NewInfixExpression(left, operator, right)
	return where
}

// Or appends given Expression using OR as logical operator.
func (where Where) Or(right Expression) Where {
	if where.IsEmpty() {
		panic("loukoum: two conditions are required for OR statement")
	}

	left := where.Condition
	operator := NewOrOperator()
	where.Condition = NewInfixExpression(left, operator, right)
	return where
}
