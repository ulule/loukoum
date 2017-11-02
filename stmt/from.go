package stmt

import (
	"bytes"
)

type From struct {
	Statement
	Table string
}

func NewFrom(table string) From {
	return From{
		Table: table,
	}
}

func (from From) Write(buffer *bytes.Buffer) {
	buffer.WriteString(from.Table)
}

// IsEmpty return true if statement is undefined.
func (from From) IsEmpty() bool {
	return from.Table == ""
}
