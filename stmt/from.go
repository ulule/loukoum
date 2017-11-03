package stmt

import (
	"bytes"
)

type From struct {
	Statement
	Table Table
}

func NewFrom(table Table) From {
	return From{
		Table: table,
	}
}

func (from From) Write(buffer *bytes.Buffer) {
	buffer.WriteString("FROM ")
	from.Table.Write(buffer)
}

// IsEmpty return true if statement is undefined.
func (from From) IsEmpty() bool {
	return from.Table.IsEmpty()
}
