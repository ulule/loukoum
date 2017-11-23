package stmt

import (
	"fmt"

	"github.com/ulule/loukoum/types"
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
		uint, uint8, uint16, uint32, uint64:
		return NewValue(value)
	case []string:
		return NewArrayString(value)
	case []int:
		return NewArrayInt(value)
	case []uint:
		return NewArrayUint(value)
	case []int8:
		return NewArrayInt8(value)
	case []uint8:
		return NewArrayUint8(value)
	case []int16:
		return NewArrayInt16(value)
	case []uint16:
		return NewArrayUint16(value)
	case []int32:
		return NewArrayInt32(value)
	case []uint32:
		return NewArrayUint32(value)
	case []int64:
		return NewArrayInt64(value)
	case []uint64:
		return NewArrayUint64(value)
	case []bool:
		return NewArrayBool(value)
	default:
		panic(fmt.Sprintf("cannot use {%+v}[%T] as loukoum Expression", value, value))
	}
}

// NewArrayExpression creates a new Expression using a list of values.
func NewArrayExpression(values ...interface{}) Expression { // nolint: gocyclo
	// We pass only one argument and it's a slice or an expression.
	if len(values) == 1 {
		value := values[0]
		switch value.(type) {
		case []string, []int, []uint, []int8, []uint8, []int16, []uint16,
			[]int32, []uint32, []int64, []uint64, []bool:
			return NewExpression(value)

		case string, int, uint, int8, uint8, int16, uint16,
			int32, uint32, int64, uint64, bool:
			return NewExpression(value)

		case Select:
			return NewExpression(value)

		default:
			panic(fmt.Sprintf("cannot use {%+v}[%T] as loukoum Expression", value, value))
		}
	}

	array := NewArray()
	for i := range values {
		switch value := values[i].(type) {
		case string:
			array.AddValue(NewValue(value))
		case int:
			array.AddValue(NewValue(value))
		case int8:
			array.AddValue(NewValue(value))
		case int16:
			array.AddValue(NewValue(value))
		case int32:
			array.AddValue(NewValue(value))
		case int64:
			array.AddValue(NewValue(value))
		case uint:
			array.AddValue(NewValue(value))
		case uint8:
			array.AddValue(NewValue(value))
		case uint16:
			array.AddValue(NewValue(value))
		case uint32:
			array.AddValue(NewValue(value))
		case uint64:
			array.AddValue(NewValue(value))
		case bool:
			array.AddValue(NewValue(value))
		case Raw:
			array.AddRaw(value)
		default:
			panic(fmt.Sprintf("cannot use {%+v}[%T] as loukoum Value", value, value))
		}
	}

	return array
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
func (identifier Identifier) Write(ctx *types.Context) {
	ctx.Write(identifier.Identifier)
}

// IsEmpty returns true if statement is undefined.
func (identifier Identifier) IsEmpty() bool {
	return identifier.Identifier == ""
}

// Equal performs an "equal" comparison.
func (identifier Identifier) Equal(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Equal)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// NotEqual performs a "not equal" comparison.
func (identifier Identifier) NotEqual(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.NotEqual)
	return NewInfixExpression(identifier, operator, NewExpression(value))
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
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// GreaterThanOrEqual performs a "greater than or equal to" comparison.
func (identifier Identifier) GreaterThanOrEqual(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.GreaterThanOrEqual)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// LessThan performs a "less than" comparison.
func (identifier Identifier) LessThan(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.LessThan)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

// LessThanOrEqual performs a "less than or equal to" comparison.
func (identifier Identifier) LessThanOrEqual(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.LessThanOrEqual)
	return NewInfixExpression(identifier, operator, NewExpression(value))
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
func (value Value) Write(ctx *types.Context) {
	if value.Value == nil {
		ctx.Write("NULL")
	} else {
		ctx.Bind(value.Value)
	}
}

// IsEmpty returns true if statement is undefined.
func (value Value) IsEmpty() bool {
	return false
}

// ----------------------------------------------------------------------------
// Array
// ----------------------------------------------------------------------------

// Array contains a list of expression values.
type Array struct {
	Values []Expression
}

// NewArray returns a an expression array.
func NewArray() Array {
	return Array{}
}

// NewArrayString returns an expression array for "string" type.
func NewArrayString(values []string) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayInt returns an expression array for "int" type.
func NewArrayInt(values []int) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayInt8 returns an expression array for "int8" type.
func NewArrayInt8(values []int8) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayInt16 returns an expression array for "int16" type.
func NewArrayInt16(values []int16) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayInt32 returns an expression array for "int32" type.
func NewArrayInt32(values []int32) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayInt64 returns an expression array for "int64" type.
func NewArrayInt64(values []int64) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayUint returns an expression array for "uint" type.
func NewArrayUint(values []uint) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayUint8 returns an expression array for "uint8" type.
func NewArrayUint8(values []uint8) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayUint16 returns an expression array for "uint16" type.
func NewArrayUint16(values []uint16) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayUint32 returns an expression array for "uint32" type.
func NewArrayUint32(values []uint32) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayUint64 returns an expression array for "uint64" type.
func NewArrayUint64(values []uint64) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

// NewArrayBool returns an expression array for "bool" type.
func NewArrayBool(values []bool) Array {
	array := NewArray()
	for i := range values {
		array.AddValue(NewValue(values[i]))
	}
	return array
}

func (Array) expression() {}

// Write exposes statement as a SQL query.
func (array Array) Write(ctx *types.Context) {
	for i := range array.Values {
		if i != 0 {
			ctx.Write(", ")
		}
		array.Values[i].Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (array Array) IsEmpty() bool {
	return len(array.Values) == 0
}

// AddValue appends a value to given array.
func (array *Array) AddValue(value Value) {
	array.Values = append(array.Values, value)
}

// AddRaw appends a raw value to given array.
func (array *Array) AddRaw(value Raw) {
	array.Values = append(array.Values, value)
}

// ----------------------------------------------------------------------------
// Raw
// ----------------------------------------------------------------------------

// Raw is an raw expression value.
type Raw struct {
	Value string
}

// NewRaw returns an raw expression value.
func NewRaw(value string) Raw {
	return Raw{
		Value: value,
	}
}

func (Raw) expression() {}

// Write exposes statement as a SQL query.
func (raw Raw) Write(ctx *types.Context) {
	ctx.Write(raw.Value)
}

// IsEmpty returns true if statement is undefined.
func (raw Raw) IsEmpty() bool {
	return false
}
