package types

// OrderType represents an order type.
type OrderType string

func (e OrderType) String() string {
	return string(e)
}

// Order types.
const (
	// Asc indicates forward order.
	Asc = OrderType("ASC")
	// Desc indicates reverse order.
	Desc = OrderType("DESC")
)
