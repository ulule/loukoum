package stmt

import (
	"github.com/ulule/loukoum/types"
)

type Table struct {
	Statement
	Name  string
	Alias string
}

func NewTable(name string) Table {
	return NewTableAlias(name, "")
}

func NewTableAlias(name, alias string) Table {
	return Table{
		Name:  name,
		Alias: alias,
	}
}

// As is used to give an alias name to the column.
func (table Table) As(alias string) Table {
	table.Alias = alias
	return table
}

func (table Table) Write(ctx *types.Context) {
	ctx.Write(table.Name)
	if table.Alias != "" {
		ctx.Write(" AS ")
		ctx.Write(table.Alias)
	}
}

// IsEmpty return true if statement is undefined.
func (table Table) IsEmpty() bool {
	return table.Name == ""
}
