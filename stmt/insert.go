package stmt

import "bytes"

// Insert is the INSERT statement.
type Insert struct {
	Statement
	Into      Into
	Columns   []Column
	Values    Values
	Returning Returning
}

// NewInsert returns a new Insert instance.
func NewInsert() Insert {
	return Insert{}
}

// Write implements Statement interface.
func (insert Insert) Write(buffer *bytes.Buffer) {
	if insert.IsEmpty() {
		panic("loukoum: an insert statement must have at least one column")
	}

	buffer.WriteString("INSERT ")
	insert.Into.Write(buffer)

	if len(insert.Columns) > 0 {
		nbColumns := len(insert.Columns)
		for i := range insert.Columns {
			if i == 0 {
				buffer.WriteString(" (")
			} else {
				buffer.WriteString(", ")
			}
			insert.Columns[i].Write(buffer)
			if i == nbColumns-1 {
				buffer.WriteString(")")
			}
		}
	}

	if !insert.Values.IsEmpty() {
		buffer.WriteString(" ")
		insert.Values.Write(buffer)
	}

	if !insert.Returning.IsEmpty() {
		buffer.WriteString(" ")
		insert.Returning.Write(buffer)
	}
}

// IsEmpty implements Statement interface.
func (insert Insert) IsEmpty() bool {
	return insert.Into.IsEmpty()
}
