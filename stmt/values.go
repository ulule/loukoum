package stmt

import (
	"bytes"
)

// Values is the VALUES clause.
type Values struct {
	Statement
	Values Expression
}

// NewValues returns a new Values instance.
func NewValues(values Expression) Values {
	return Values{
		Values: values,
	}
}

// Write implements Statement interface.
func (values Values) Write(buffer *bytes.Buffer) {
	if values.IsEmpty() {
		return
	}

	buffer.WriteString("VALUES (")
	values.Values.Write(buffer)
	buffer.WriteString(")")
}

// IsEmpty implements Statement interface.
func (values Values) IsEmpty() bool {
	return values.Values == nil || (values.Values != nil && values.Values.IsEmpty())
}
