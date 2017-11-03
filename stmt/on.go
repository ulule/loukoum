package stmt

import (
	"bytes"
)

type On struct {
	Left  Column
	Right Column
}

func NewOn(left, right Column) On {
	return On{
		Left:  left,
		Right: right,
	}
}

func (on On) Write(buffer *bytes.Buffer) {
	buffer.WriteString("ON ")
	buffer.WriteString(on.Left.Name)
	buffer.WriteString(" = ")
	buffer.WriteString(on.Right.Name)
}

// IsEmpty return true if statement is undefined.
func (on On) IsEmpty() bool {
	return on.Left.IsEmpty() || on.Right.IsEmpty()
}
