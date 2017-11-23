package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// Column is a column identifier.
type Column struct {
	Statement
	Name  string
	Alias string
}

// NewColumn returns a new Column instance.
func NewColumn(name string) Column {
	return NewColumnAlias(name, "")
}

// NewColumnAlias returns a new Column instance with an alias.
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

// Write exposes statement as a SQL query.
func (column Column) Write(ctx *types.Context) {
	ctx.Write(column.Name)
	if column.Alias != "" {
		ctx.Write(" ")
		ctx.Write(token.As.String())
		ctx.Write(" ")
		ctx.Write(column.Alias)
	}
}

// IsEmpty returns true if statement is undefined.
func (column Column) IsEmpty() bool {
	return column.Name == ""
}
