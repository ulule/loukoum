package stmt

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ulule/loukoum/types"
)

// Where is a WHERE clause.
type Where struct {
	Statement
	Operator   types.LogicalOperator
	Conditions []Condition
}

// NewWhere returns a new WHERE clause.
func NewWhere(operator types.LogicalOperator, conditions ...Condition) Where {
	return Where{
		Operator:   operator,
		Conditions: conditions,
	}
}

// Write writes WHERE clause into the given buffer.
func (where Where) Write(buffer *bytes.Buffer) {
	if where.IsEmpty() {
		panic("loukoum: a where clause expects at least one condition")
	}

	buffer.WriteString(" (")
	for i, condition := range where.Conditions {
		if i > 0 && condition.operator != "" {
			fmt.Fprintf(buffer, " %s ", condition.operator.String())
		}
		condition.Write(buffer)
	}
	buffer.WriteString(")")
}

// IsEmpty return true if statement is undefined.
func (where Where) IsEmpty() bool {
	return len(where.Conditions) == 0
}

// Condition is a WHERE condition
type Condition struct {
	Statement
	operator   types.LogicalOperator
	comparator types.ComparisonOperator
	column     Column
	values     []interface{}
}

// NewCondition returns a new WHERE condition.
func NewCondition(column Column, operator types.LogicalOperator) Condition {
	return Condition{
		operator: operator,
		column:   column,
	}
}

// Equal performs a WHERE equal.
func (c Condition) Equal(value interface{}) Condition {
	c.comparator = types.Equal
	c.values = append(c.values, value)
	return c
}

// NotEqual performs a WHERE not equal.
func (c Condition) NotEqual(value interface{}) Condition {
	c.comparator = types.NotEqual
	c.values = append(c.values, value)
	return c
}

// Is performs a WHERE is.
func (c Condition) Is(value interface{}) Condition {
	c.comparator = types.Is
	c.values = append(c.values, value)
	return c
}

// IsNot performs a WHERE is not.
func (c Condition) IsNot(value interface{}) Condition {
	c.comparator = types.IsNot
	c.values = append(c.values, value)
	return c
}

// GreaterThan performs a WHERE greater than.
func (c Condition) GreaterThan(value interface{}) Condition {
	c.comparator = types.GreaterThan
	c.values = append(c.values, value)
	return c
}

// GreaterThanOrEqual performs a WHERE greater than or equal.
func (c Condition) GreaterThanOrEqual(value interface{}) Condition {
	c.comparator = types.GreaterThanOrEqual
	c.values = append(c.values, value)
	return c
}

// LessThan performs a WHERE less than.
func (c Condition) LessThan(value interface{}) Condition {
	c.comparator = types.LessThan
	c.values = append(c.values, value)
	return c
}

// LessThanOrEqual performs a WHERE less than or equal.
func (c Condition) LessThanOrEqual(value interface{}) Condition {
	c.comparator = types.LessThanOrEqual
	c.values = append(c.values, value)
	return c
}

// In performs a WHERE in.
func (c Condition) In(values []interface{}) Condition {
	c.comparator = types.In
	c.values = append(c.values, values...)
	return c
}

// NotIn performs a WHERE not in.
func (c Condition) NotIn(values []interface{}) Condition {
	c.comparator = types.NotIn
	c.values = append(c.values, values...)
	return c
}

// Like performs a WHERE like.
func (c Condition) Like(value interface{}) Condition {
	c.comparator = types.Like
	c.values = append(c.values, value)
	return c
}

// NotLike performs a WHERE not like.
func (c Condition) NotLike(value interface{}) Condition {
	c.comparator = types.NotLike
	c.values = append(c.values, value)
	return c
}

// ILike performs a WHERE ilike.
func (c Condition) ILike(value interface{}) Condition {
	c.comparator = types.ILike
	c.values = append(c.values, value)
	return c
}

// NotILike performs a WHERE not ilike.
func (c Condition) NotILike(value interface{}) Condition {
	c.comparator = types.NotILike
	c.values = append(c.values, value)
	return c
}

// Between performs a WHERE between.
func (c Condition) Between(values []interface{}) Condition {
	c.comparator = types.Between
	c.values = append(c.values, values...)
	return c
}

// NotBetween performs a WHERE not between.
func (c Condition) NotBetween(values []interface{}) Condition {
	c.comparator = types.NotBetween
	c.values = append(c.values, values...)
	return c
}

// Write writes WHERE clause into the given buffer.
func (c Condition) Write(buffer *bytes.Buffer) {
	if len(c.values) == 0 {
		panic("loukoum: a where clause expects at least a value")
	}

	fmt.Fprintf(buffer, "%s %s ", c.column.Name, c.comparator)

	switch c.comparator {
	case types.Between, types.NotBetween:
		if len(c.values) != 2 {
			panic("loukoum: a between condition expects at least two values")
		}
		fmt.Fprintf(buffer, "%v AND %v", c.values[0], c.values[1])
	case types.In, types.NotIn:
		var values []string
		for _, v := range c.values {
			values = append(values, fmt.Sprintf("%v", v))
		}
		fmt.Fprintf(buffer, "(%s)", strings.Join(values, ","))
	default:
		if len(c.values) > 1 {
			panic("loukoum: a between condition expects at least two values")
		}
		fmt.Fprintf(buffer, "%v", c.values[0])
	}
}

// IsEmpty return true if statement is undefined.
func (c Condition) IsEmpty() bool {
	return c.column.IsEmpty() || len(c.values) == 0
}
