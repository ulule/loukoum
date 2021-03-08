package stmt

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/ulule/loukoum/v3/types"
)

// Expression is a SQL expression.
type Expression interface {
	Statement
	expression()
}

// NewExpression returns a new Expression instance from arg.
func NewExpression(arg interface{}) Expression { // nolint: gocyclo
	if arg == nil {
		return NewValue(nil)
	}
	switch value := arg.(type) {
	case Expression:
		return value
	case string, bool, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, float32, float64, []byte:
		return NewValue(value)
	case time.Time:
		return NewValue(value)
	case *time.Time:
		return NewValue(*value)
	case driver.Valuer:
		return NewValue(value)
	case StatementEncoder:
		stmt := value.Statement()
		expression, ok := stmt.(Expression)
		if !ok {
			panic(fmt.Sprintf("cannot use {%+v}[%T] as loukoum Expression", value, value))
		}
		return expression
	case Int64Encoder:
		return NewValue(value.Int64())
	case BoolEncoder:
		return NewValue(value.Bool())
	case TimeEncoder:
		return NewValue(value.Time())
	case StringEncoder:
		return NewValue(value.String())
	default:
		panic(fmt.Sprintf("cannot use {%+v}[%T] as loukoum Expression", arg, arg))
	}
}

// ----------------------------------------------------------------------------
// Identifier
// ----------------------------------------------------------------------------

// Identifier is an identifier.
type Identifier struct {
	Identifier string
}

// NewIdentifier returns a new Identifier.
func NewIdentifier(identifier string) Identifier {
	return Identifier{
		Identifier: identifier,
	}
}

func (Identifier) expression() {}

// Write exposes statement as a SQL query.
func (identifier Identifier) Write(ctx types.Context) {
	ctx.Write(identifier.Identifier)
}

// IsEmpty returns true if statement is undefined.
func (identifier Identifier) IsEmpty() bool {
	return identifier.Identifier == ""
}

// Contains performs a "contains" comparison.
func (identifier Identifier) Contains(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Contains)
	return NewInfixExpression(identifier, operator, NewWrapper(NewExpression(value)))
}

// IsContainedBy performs a "is contained by" comparison.
func (identifier Identifier) IsContainedBy(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.IsContainedBy)
	return NewInfixExpression(identifier, operator, NewWrapper(NewExpression(value)))
}

// Overlap performs an "overlap" comparison.
func (identifier Identifier) Overlap(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Overlap)
	return NewInfixExpression(identifier, operator, NewWrapper(NewExpression(value)))
}

// Equal performs an "equal" comparison.
func (identifier Identifier) Equal(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Equal)
	return NewInfixExpression(identifier, operator, NewWrapper(NewExpression(value)))
}

// NotEqual performs a "not equal" comparison.
func (identifier Identifier) NotEqual(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.NotEqual)
	return NewInfixExpression(identifier, operator, NewWrapper(NewExpression(value)))
}

// Is performs a "is" comparison.
func (identifier Identifier) Is(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Is)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// IsNot performs a "is not" comparison.
func (identifier Identifier) IsNot(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.IsNot)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// IsNull performs a "is null" comparison.
func (identifier Identifier) IsNull(value bool) InfixExpression {
	if value {
		return identifier.Is(nil)
	}
	return identifier.IsNot(nil)
}

// GreaterThan performs a "greater than" comparison.
func (identifier Identifier) GreaterThan(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.GreaterThan)
	return NewInfixExpression(identifier, operator, NewWrapper(NewExpression(value)))
}

// GreaterThanOrEqual performs a "greater than or equal to" comparison.
func (identifier Identifier) GreaterThanOrEqual(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.GreaterThanOrEqual)
	return NewInfixExpression(identifier, operator, NewWrapper(NewExpression(value)))
}

// LessThan performs a "less than" comparison.
func (identifier Identifier) LessThan(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.LessThan)
	return NewInfixExpression(identifier, operator, NewWrapper(NewExpression(value)))
}

// LessThanOrEqual performs a "less than or equal to" comparison.
func (identifier Identifier) LessThanOrEqual(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.LessThanOrEqual)
	return NewInfixExpression(identifier, operator, NewWrapper(NewExpression(value)))
}

// In performs a "in" condition.
func (identifier Identifier) In(value ...interface{}) In {
	return NewIn(identifier, NewArrayExpression(value...))
}

// NotIn performs a "not in" condition.
func (identifier Identifier) NotIn(value ...interface{}) In {
	return NewNotIn(identifier, NewArrayExpression(value...))
}

// Like performs a "like" condition.
func (identifier Identifier) Like(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Like)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// NotLike performs a "not like" condition.
func (identifier Identifier) NotLike(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.NotLike)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// ILike performs a "ilike" condition.
func (identifier Identifier) ILike(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.ILike)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// NotILike performs a "not ilike" condition.
func (identifier Identifier) NotILike(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.NotILike)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// Between performs a "between" condition.
func (identifier Identifier) Between(from, to interface{}) Between {
	return NewBetween(identifier, NewExpression(from), NewExpression(to))
}

// NotBetween performs a "not between" condition.
func (identifier Identifier) NotBetween(from, to interface{}) Between {
	return NewNotBetween(identifier, NewExpression(from), NewExpression(to))
}

// IsDistinctFrom performs an "is distinct from" comparison.
func (identifier Identifier) IsDistinctFrom(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.IsDistinctFrom)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// IsNotDistinctFrom performs an "is not distinct from" comparison.
func (identifier Identifier) IsNotDistinctFrom(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.IsNotDistinctFrom)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// Ensure that Identifier is an Expression
var _ Expression = Identifier{}

// ----------------------------------------------------------------------------
// Value
// ----------------------------------------------------------------------------

// Value is an expression value.
type Value struct {
	Value interface{}
}

// NewValue returns an expression value.
func NewValue(value interface{}) Value {
	return Value{
		Value: value,
	}
}

func (Value) expression() {}

// Write exposes statement as a SQL query.
func (value Value) Write(ctx types.Context) {
	if value.Value == nil {
		ctx.Write("NULL")
	} else {
		ctx.Bind(value.Value)
	}
}

// Overlap performs an "overlap" comparison.
func (value Value) Overlap(what interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Overlap)
	return NewInfixExpression(value, operator, NewWrapper(NewExpression(what)))
}

// Equal performs an "equal" comparison.
func (value Value) Equal(what interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Equal)
	return NewInfixExpression(value, operator, NewWrapper(NewExpression(what)))
}

// NotEqual performs a "not equal" comparison.
func (value Value) NotEqual(what interface{}) InfixExpression {
	operator := NewComparisonOperator(types.NotEqual)
	return NewInfixExpression(value, operator, NewWrapper(NewExpression(what)))
}

// Contains performs a "contains" comparison.
func (value Value) Contains(what interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Contains)
	return NewInfixExpression(value, operator, NewWrapper(NewExpression(what)))
}

// IsContainedBy performs a "is contained by" comparison.
func (value Value) IsContainedBy(what interface{}) InfixExpression {
	operator := NewComparisonOperator(types.IsContainedBy)
	return NewInfixExpression(value, operator, NewWrapper(NewExpression(what)))
}

// IsEmpty returns true if statement is undefined.
func (value Value) IsEmpty() bool {
	return false
}

// Ensure that Value is an Expression
var _ Expression = Value{}

// ----------------------------------------------------------------------------
// Array
// ----------------------------------------------------------------------------

// Array contains a list of expression values.
type Array struct {
	Values []Expression
}

// NewArrayExpression creates a new Expression using a list of values.
func NewArrayExpression(values ...interface{}) Expression { // nolint: gocyclo
	// We pass only one argument and it's a slice or an expression.
	if len(values) == 1 {
		return toArray(values[0])
	}
	array := Array{}
	for _, value := range values {
		array.Append(value)
	}
	return array
}

// toArray tries to cast the value to a slice.
// It returns a single element Array otherwise.
func toArray(value interface{}) Array { // nolint: gocyclo
	array := Array{}
	if value == nil {
		return array
	}

	switch values := value.(type) {
	case []string:
		for _, v := range values {
			array.Append(v)
		}
	case []int:
		for _, v := range values {
			array.Append(v)
		}
	case []uint:
		for _, v := range values {
			array.Append(v)
		}
	case []int8:
		for _, v := range values {
			array.Append(v)
		}
	case []int16:
		for _, v := range values {
			array.Append(v)
		}
	case []uint16:
		for _, v := range values {
			array.Append(v)
		}
	case []int32:
		for _, v := range values {
			array.Append(v)
		}
	case []uint32:
		for _, v := range values {
			array.Append(v)
		}
	case []int64:
		for _, v := range values {
			array.Append(v)
		}
	case []uint64:
		for _, v := range values {
			array.Append(v)
		}
	case []bool:
		for _, v := range values {
			array.Append(v)
		}
	case []float32:
		for _, v := range values {
			array.Append(v)
		}
	case []float64:
		for _, v := range values {
			array.Append(v)
		}
	case [][]byte:
		for _, v := range values {
			array.Append(v)
		}
	case []Expression:
		for _, v := range values {
			array.Append(v)
		}
	case []interface{}:
		for _, v := range values {
			array.Append(v)
		}
	default:
		array.Append(value)
	}
	return array
}

func (Array) expression() {}

// Write exposes statement as a SQL query.
func (array Array) Write(ctx types.Context) {
	for i, value := range array.Values {
		if i > 0 {
			ctx.Write(", ")
		}
		value.Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (array Array) IsEmpty() bool {
	return len(array.Values) == 0
}

// Append an expression to the given array.
func (array *Array) Append(value interface{}) {
	array.Values = append(array.Values, NewExpression(value))
}

// Ensure that Array is an Expression
var _ Expression = Array{}

// ----------------------------------------------------------------------------
// Raw
// ----------------------------------------------------------------------------

// Raw is a raw expression value.
type Raw struct {
	Value string
}

// NewRaw returns a raw expression value.
func NewRaw(value string) Raw {
	return Raw{
		Value: value,
	}
}

func (Raw) expression() {}

// Write exposes statement as a SQL query.
func (raw Raw) Write(ctx types.Context) {
	ctx.Write(raw.Value)
}

// IsEmpty returns true if statement is undefined.
func (raw Raw) IsEmpty() bool {
	return raw.Value == ""
}

// Ensure that Raw is an Expression
var _ Expression = Raw{}

// ----------------------------------------------------------------------------
// Wrapper
// ----------------------------------------------------------------------------

// Wrapper encapsulates an expression between parenthesis.
type Wrapper struct {
	Value Expression
}

// NewWrapper returns a new Wrapper expression when it's required.
func NewWrapper(arg Expression) Expression {
	switch value := arg.(type) {
	case Select:
		return &Wrapper{
			Value: value,
		}
	case Exists:
		return &Wrapper{
			Value: value,
		}
	case NotExists:
		return &Wrapper{
			Value: value,
		}
	default:
		return arg
	}
}

func (Wrapper) expression() {}

// Write exposes statement as a SQL query.
func (wrapper Wrapper) Write(ctx types.Context) {
	ctx.Write("(")
	wrapper.Value.Write(ctx)
	ctx.Write(")")
}

// IsEmpty returns true if statement is undefined.
func (wrapper Wrapper) IsEmpty() bool {
	return wrapper.Value.IsEmpty()
}

// Ensure that Wrapper is an Expression
var _ Expression = Wrapper{}
