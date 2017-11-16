package stmt

import (
	"bytes"
)

type Select struct {
	Distinct bool
	Columns  []Column
	From     From
	Joins    []Join
	Where    Where
}

func NewSelect() Select {
	return Select{}
}

func (selekt Select) Write(buffer *bytes.Buffer) {
	if selekt.IsEmpty() {
		panic("loukoum: select statements must have at least one column")
	}

	// TODO Add prefixes

	buffer.WriteString("SELECT")

	if selekt.Distinct {
		buffer.WriteString(" DISTINCT")
	}

	for i := range selekt.Columns {
		if i == 0 {
			buffer.WriteString(" ")
		} else {
			buffer.WriteString(", ")
		}
		selekt.Columns[i].Write(buffer)
	}

	if !selekt.From.IsEmpty() {
		buffer.WriteString(" ")
		selekt.From.Write(buffer)
	}

	for i := range selekt.Joins {
		buffer.WriteString(" ")
		selekt.Joins[i].Write(buffer)
	}

	if !selekt.Where.IsEmpty() {
		buffer.WriteString(" ")
		selekt.Where.Write(buffer)
	}

	// TODO GROUP BY

	// TODO HAVING

	// TODO ORDER BY

	// TODO LIMIT

	// TODO OFFSET

	// TODO Add suffixes
}

// IsEmpty return true if statement is undefined.
func (selekt Select) IsEmpty() bool {
	return len(selekt.Columns) == 0
}
