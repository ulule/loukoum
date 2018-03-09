package stmt

import (
	"github.com/ulule/loukoum/types"
)

// Order is an expression of a ORDER BY clause.
type Order struct {
	Expression string
	Type       types.OrderType
}

// NewOrder returns a new Order instance.
func NewOrder(expression string, kind types.OrderType) Order {
	return Order{
		Expression: expression,
		Type:       kind,
	}
}

// Write exposes statement as a SQL query.
func (order Order) Write(ctx types.Context) {
	if order.IsEmpty() {
		return
	}
	ctx.Write(order.Expression)
	ctx.Write(" ")
	ctx.Write(order.Type.String())
}

// IsEmpty returns true if statement is undefined.
func (order Order) IsEmpty() bool {
	return order.Expression == ""
}

// Ensure that Order is a Statement
var _ Statement = Order{}
