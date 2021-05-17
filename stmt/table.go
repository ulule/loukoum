package stmt

import (
	"github.com/ulule/loukoum/v3/token"
	"github.com/ulule/loukoum/v3/types"
)

// Table is a table identifier.
type Table struct {
	Alias string
	Name  string
	only  bool
}

// NewTable returns a new Table instance.
func NewTable(name string) Table {
	return NewTableAlias(name, "")
}

// Only sets ONLY clause to the table.
func (table Table) Only() Table {
	table.only = true
	return table
}

// NewTableAlias returns a new Table instance with an alias.
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

// Write exposes statement as a SQL query.
func (table Table) Write(ctx types.Context) {
	if table.only {
		ctx.Write(token.Only.String())
		ctx.Write(" ")
	}
	ctx.Write(quote(table.Name))
	if table.Alias != "" {
		ctx.Write(" ")
		ctx.Write(token.As.String())
		ctx.Write(" ")
		ctx.Write(quote(table.Alias))
	}
}

// IsEmpty returns true if statement is undefined.
func (table Table) IsEmpty() bool {
	return table.Name == ""
}

// Ensure that Table is a Statement
var _ Statement = Table{}
