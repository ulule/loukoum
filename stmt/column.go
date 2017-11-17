package stmt

import (
	"fmt"

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

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

// ToColumns takes a list of empty interfaces and returns a slice of Column instance.
func ToColumns(values []interface{}) []Column {
	columns := make([]Column, 0, len(values))

	for i := range values {
		column := Column{}

		switch value := values[i].(type) {
		case string:
			column = NewColumn(value)
		case Column:
			column = value
		default:
			panic(fmt.Sprintf("loukoum: cannot use %T as column", column))
		}
		if column.IsEmpty() {
			panic("loukoum: a column was undefined")
		}

		columns = append(columns, column)
	}

	return columns
}
