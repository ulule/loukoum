package stmt

import (
	"github.com/ulule/loukoum/v3/token"
	"github.com/ulule/loukoum/v3/types"
)

// Column is a column identifier.
type Column struct {
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

// Asc is used to transform a column to an order expression.
func (column Column) Asc() Order {
	expression := column.Name
	if column.Alias != "" {
		expression = column.Alias
	}

	return NewOrder(expression, types.Asc)
}

// Desc is used to transform a column to an order expression.
func (column Column) Desc() Order {
	expression := column.Name
	if column.Alias != "" {
		expression = column.Alias
	}

	return NewOrder(expression, types.Desc)
}

// Write exposes statement as a SQL query.
func (column Column) Write(ctx types.Context) {
	ctx.Write(quote(column.Name))
	if column.Alias != "" {
		ctx.Write(" ")
		ctx.Write(token.As.String())
		ctx.Write(" ")
		ctx.Write(quote(column.Alias))
	}
}

// IsEmpty returns true if statement is undefined.
func (column Column) IsEmpty() bool {
	return column.Name == ""
}

func (Column) selectExpression() {}

// Ensure that Column is a SelectExpression.
var _ SelectExpression = Column{}
