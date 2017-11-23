package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// OrderBy is a ORDER BY clause.
type OrderBy struct {
	Orders []Order
}

// NewOrderBy returns a new OrderBy instance.
func NewOrderBy(orders []Order) OrderBy {
	return OrderBy{
		Orders: orders,
	}
}

// Write exposes statement as a SQL query.
func (order OrderBy) Write(ctx *types.Context) {
	if order.IsEmpty() {
		return
	}
	ctx.Write(token.Order.String())
	ctx.Write(" ")
	ctx.Write(token.By.String())
	ctx.Write(" ")
	for i := range order.Orders {
		if i != 0 {
			ctx.Write(", ")
		}
		order.Orders[i].Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (order OrderBy) IsEmpty() bool {
	return len(order.Orders) == 0
}

// Ensure that OrderBy is a Statement
var _ Statement = OrderBy{}
