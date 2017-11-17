package stmt

import "bytes"

// Returning is the RETURNING clause.
type Returning struct {
	Statement
	Columns []Column
}

// NewReturning returns a new Returning instance.
func NewReturning(columns []Column) Returning {
	return Returning{
		Columns: columns,
	}
}

// Write implements Statement interface.
func (returning Returning) Write(buffer *bytes.Buffer) {
	buffer.WriteString("RETURNING ")

	l := len(returning.Columns)
	if l > 1 {
		buffer.WriteString("(")
	}

	for i := range returning.Columns {
		if i > 0 {
			buffer.WriteString(", ")
		}
		returning.Columns[i].Write(buffer)
	}

	if l > 1 {
		buffer.WriteString(")")
	}
}

// IsEmpty implements Statement interface.
func (returning Returning) IsEmpty() bool {
	return len(returning.Columns) == 0
}
