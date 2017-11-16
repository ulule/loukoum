package stmt

import "bytes"

// Into is the INTO clause.
type Into struct {
	Statement
	Table Table
}

// NewInto returns a new Into instance.
func NewInto(table Table) Into {
	return Into{
		Table: table,
	}
}

// Write implements Statement interface.
func (into Into) Write(buffer *bytes.Buffer) {
	buffer.WriteString("INTO ")
	into.Table.Write(buffer)
}

// IsEmpty implements Statement interface.
func (into Into) IsEmpty() bool {
	return into.Table.IsEmpty()
}
