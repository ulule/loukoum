package loukoum

import (
	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

const (
	// InnerJoin is used for "INNER JOIN" in join statement.
	InnerJoin = types.InnerJoin
	// LeftJoin is used for "LEFT JOIN" in join statement.
	LeftJoin = types.LeftJoin
	// RightJoin is used for "RIGHT JOIN" in join statement.
	RightJoin = types.RightJoin
	// Asc is used for "ORDER BY" statement.
	Asc = types.Asc
	// Desc is used for "ORDER BY" statement.
	Desc = types.Desc
)

// Map is a key/value map.
type Map = types.Map

// Pair takes a key and its related value and returns a Pair.
func Pair(key, value interface{}) types.Pair {
	return types.Pair{Key: key, Value: value}
}

// Select starts a SelectBuilder using the given columns.
func Select(columns ...interface{}) builder.Select {
	return builder.NewSelect().Columns(columns...)
}

// Column is a wrapper to create a new Column statement.
func Column(name string) stmt.Column {
	return stmt.NewColumn(name)
}

// Table is a wrapper to create a new Table statement.
func Table(name string) stmt.Table {
	return stmt.NewTable(name)
}

// On is a wrapper to create a new On statement.
func On(left string, right string) stmt.On {
	return stmt.NewOn(stmt.NewColumn(left), stmt.NewColumn(right))
}

// Condition is a wrapper to create a new Identifier statement.
func Condition(column string) stmt.Identifier {
	return stmt.NewIdentifier(column)
}

// Order is a wrapper to create a new Order statement.
func Order(column string, option ...types.OrderType) stmt.Order {
	order := types.Asc
	if len(option) > 0 {
		order = option[0]
	}
	return stmt.NewOrder(column, order)
}

// And is a wrapper to create a new InfixExpression statement.
func And(left stmt.Expression, right stmt.Expression) stmt.InfixExpression {
	return stmt.NewInfixExpression(left, stmt.NewLogicalOperator(types.And), right)
}

// Or is a wrapper to create a new InfixExpression statement.
func Or(left stmt.Expression, right stmt.Expression) stmt.InfixExpression {
	return stmt.NewInfixExpression(left, stmt.NewLogicalOperator(types.Or), right)
}

// Raw is a wrapper to create a new Raw expression.
func Raw(value string) stmt.Raw {
	return stmt.NewRaw(value)
}

// Insert starts an InsertBuilder using the given table as into clause.
func Insert(into interface{}) builder.Insert {
	return builder.NewInsert().Into(into)
}

// Delete starts a DeleteBuilder using the given table as from clause.
func Delete(from interface{}) builder.Delete {
	return builder.NewDelete().From(from)
}

// Update starts an Update builder using the given table.
func Update(table interface{}) builder.Update {
	return builder.NewUpdate(table)
}
