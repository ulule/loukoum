package stmt

import (
	"bytes"
)

type Column struct {
	Name  string
	Alias string
}

func NewColumn(name string) Column {
	return NewColumnAlias(name, "")
}

func NewColumnAlias(name, alias string) Column {
	return Column{
		Name:  name,
		Alias: alias,
	}
}

// As is used to give an alias name to the column.
func (column Column) As(alias string) Column {
	column.Alias = alias
	return column
}

func (column Column) Write(buffer *bytes.Buffer) {
	buffer.WriteString(column.Name)
	if column.Alias != "" {
		buffer.WriteString(" AS ")
		buffer.WriteString(column.Alias)
	}
}
