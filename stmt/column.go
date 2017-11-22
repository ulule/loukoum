package stmt

import (
	"github.com/ulule/loukoum/types"
)

type Column struct {
	Statement
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

// Write expose statement as a SQL query.
func (column Column) Write(ctx *types.Context) {
	ctx.Write(column.Name)
	if column.Alias != "" {
		ctx.Write(" AS ")
		ctx.Write(column.Alias)
	}
}

// IsEmpty return true if statement is undefined.
func (column Column) IsEmpty() bool {
	return column.Name == ""
}
