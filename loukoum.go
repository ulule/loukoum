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
	// LeftOuterJoin is used for "LEFT OUTER JOIN" in join statement.
	LeftOuterJoin = types.LeftOuterJoin
	// RightOuterJoin is used for "RIGHT OUTER JOIN" in join statement.
	RightOuterJoin = types.RightOuterJoin
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
func On(left string, right string) stmt.OnClause {
	return stmt.NewOnClause(stmt.NewColumn(left), stmt.NewColumn(right))
}

// AndOn is a wrapper to create a new On statement using an infix expression.
func AndOn(left stmt.OnExpression, right stmt.OnExpression) stmt.OnExpression {
	return stmt.NewInfixOnExpression(left, stmt.NewLogicalOperator(types.And), right)
}

// OrOn is a wrapper to create a new On statement using an infix expression.
func OrOn(left stmt.OnExpression, right stmt.OnExpression) stmt.OnExpression {
	return stmt.NewInfixOnExpression(left, stmt.NewLogicalOperator(types.Or), right)
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

// Exists is a wrapper to create a new Exists expression.
func Exists(value interface{}) stmt.Exists {
	return stmt.NewExists(value)
}

// NotExists is a wrapper to create a new NotExists expression.
func NotExists(value interface{}) stmt.NotExists {
	return stmt.NewNotExists(value)
}

// Count is a wrapper to create a new Count expression.
func Count(value string) stmt.Count {
	return stmt.NewCount(value)
}

// Max is a wrapper to create a new Max expression.
func Max(value string) stmt.Max {
	return stmt.NewMax(value)
}

// Min is a wrapper to create a new Min expression.
func Min(value string) stmt.Min {
	return stmt.NewMin(value)
}

// Sum is a wrapper to create a new Sum expression.
func Sum(value string) stmt.Sum {
	return stmt.NewSum(value)
}

// With is a wrapper to create a new WithQuery statement.
func With(name string, value interface{}) stmt.WithQuery {
	return stmt.NewWithQuery(name, value)
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

// DoNothing is a wrapper to create a new ConflictNoAction statement.
func DoNothing() stmt.ConflictNoAction {
	return stmt.NewConflictNoAction()
}

// DoUpdate is a wrapper to create a new ConflictUpdateAction statement.
func DoUpdate(args ...interface{}) stmt.ConflictUpdateAction {
	return stmt.NewConflictUpdateAction(builder.ToSet(args))
}
