package stmt

import (
	"bytes"
)

// Values is the VALUES clause.
type Values struct {
	Statement
	Expression Expression
}

// NewValues returns a new Values instance.
func NewValues(expr Expression) Values {
	return Values{
		Expression: expr,
	}
}

// Write implements Statement interface.
func (values Values) Write(buffer *bytes.Buffer) {
	buffer.WriteString("VALUES (")
	values.Expression.Write(buffer)
	buffer.WriteString(")")
}

// IsEmpty implements Statement interface.
func (values Values) IsEmpty() bool {
	return values.Expression != nil && values.Expression.IsEmpty()
}
