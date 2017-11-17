package stmt

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/ulule/loukoum/types"
)

// Expression is a SQL expression.
type Expression interface {
	Statement
	expression()
}

// NewExpression returns a new Expression instance from arg.
func NewExpression(arg interface{}) Expression {
	if arg == nil {
		return NewValue("NULL")
	}

	switch value := arg.(type) {
	case Expression:
		return value
	case string:
		return NewValueString(value)
	case int:
		return NewValueInt(value)
	case int8:
		return NewValueInt8(value)
	case int16:
		return NewValueInt16(value)
	case int32:
		return NewValueInt32(value)
	case int64:
		return NewValueInt64(value)
	case uint:
		return NewValueUint(value)
	case uint8:
		return NewValueUint8(value)
	case uint16:
		return NewValueUint16(value)
	case uint32:
		return NewValueUint32(value)
	case uint64:
		return NewValueUint64(value)
	case bool:
		return NewValueBool(value)
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

func (identifier Identifier) Write(buffer *bytes.Buffer) {
	buffer.WriteString(identifier.Identifier)
}

// IsEmpty return true if statement is undefined.
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
	return NewIn(identifier, NewInExpression(value...))
}

// NotIn performs a "not in" condition.
func (identifier Identifier) NotIn(value ...interface{}) In {
	return NewNotIn(identifier, NewInExpression(value...))
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
	Value string
}

// NewValue returns an expression value.
func NewValue(value string) Value {
	return Value{
		Value: value,
	}
}

// NewValueString returns an expression value for "string" type.
func NewValueString(value string) Value {
	return NewValue(value)
}

// NewValueInt returns an expression value for "int" type.
func NewValueInt(value int) Value {
	return NewValueInt64(int64(value))
}

// NewValueInt8 returns an expression value for "int8" type.
func NewValueInt8(value int8) Value {
	return NewValueInt64(int64(value))
}

// NewValueInt16 returns an expression value for "int16" type.
func NewValueInt16(value int16) Value {
	return NewValueInt64(int64(value))
}

// NewValueInt32 returns an expression value for "int32" type.
func NewValueInt32(value int32) Value {
	return NewValueInt64(int64(value))
}

// NewValueInt64 returns an expression value for "int64" type.
func NewValueInt64(value int64) Value {
	return NewValue(strconv.FormatInt(value, 10))
}

// NewValueUint returns an expression value for "uint" type.
func NewValueUint(value uint) Value {
	return NewValueUint64(uint64(value))
}

// NewValueUint8 returns an expression value for "uint8" type.
func NewValueUint8(value uint8) Value {
	return NewValueUint64(uint64(value))
}

// NewValueUint16 returns an expression value for "uint16" type.
func NewValueUint16(value uint16) Value {
	return NewValueUint64(uint64(value))
}

// NewValueUint32 returns an expression value for "uint32" type.
func NewValueUint32(value uint32) Value {
	return NewValueUint64(uint64(value))
}

// NewValueUint64 returns an expression value for "uint64" type.
func NewValueUint64(value uint64) Value {
	return NewValue(strconv.FormatUint(value, 10))
}

// NewValueBool returns an expression value for "bool" type.
func NewValueBool(value bool) Value {
	return NewValue(strconv.FormatBool(value))
}

func (Value) expression() {}

func (value Value) Write(buffer *bytes.Buffer) {
	buffer.WriteString(value.Value)
}

// IsEmpty return true if statement is undefined.
func (value Value) IsEmpty() bool {
	return value.Value == ""
}

// Array contains a list of expression values.
type Array struct {
	Values []Value
}

// NewArray returns a an expression array.
func NewArray() Array {
	return Array{}
}

// NewArrayString returns an expression array for "string" type.
func NewArrayString(values []string) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueString(values[i]))
	}
	return array
}

// NewArrayInt returns an expression array for "int" type.
func NewArrayInt(values []int) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueInt(values[i]))
	}
	return array
}

// NewArrayInt8 returns an expression array for "int8" type.
func NewArrayInt8(values []int8) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueInt8(values[i]))
	}
	return array
}

// NewArrayInt16 returns an expression array for "int16" type.
func NewArrayInt16(values []int16) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueInt16(values[i]))
	}
	return array
}

// NewArrayInt32 returns an expression array for "int32" type.
func NewArrayInt32(values []int32) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueInt32(values[i]))
	}
	return array
}

// NewArrayInt64 returns an expression array for "int64" type.
func NewArrayInt64(values []int64) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueInt64(values[i]))
	}
	return array
}

// NewArrayUint returns an expression array for "uint" type.
func NewArrayUint(values []uint) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueUint(values[i]))
	}
	return array
}

// NewArrayUint8 returns an expression array for "uint8" type.
func NewArrayUint8(values []uint8) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueUint8(values[i]))
	}
	return array
}

// NewArrayUint16 returns an expression array for "uint16" type.
func NewArrayUint16(values []uint16) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueUint16(values[i]))
	}
	return array
}

// NewArrayUint32 returns an expression array for "uint32" type.
func NewArrayUint32(values []uint32) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueUint32(values[i]))
	}
	return array
}

// NewArrayUint64 returns an expression array for "uint64" type.
func NewArrayUint64(values []uint64) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueUint64(values[i]))
	}
	return array
}

// NewArrayBool returns an expression array for "bool" type.
func NewArrayBool(values []bool) Array {
	array := NewArray()
	for i := range values {
		array.Add(NewValueBool(values[i]))
	}
	return array
}

func (Array) expression() {}

func (array Array) Write(buffer *bytes.Buffer) {
	for i := range array.Values {
		if i != 0 {
			buffer.WriteString(", ")
		}
		array.Values[i].Write(buffer)
	}
}

// Add append a value to given array.
func (array *Array) Add(value Value) {
	array.Values = append(array.Values, value)
}

// IsEmpty return true if statement is undefined.
func (array Array) IsEmpty() bool {
	return len(array.Values) == 0
}

// NewInExpression creates a new Expression for IN value.
func NewInExpression(values ...interface{}) Expression {
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
			array.Add(NewValueString(value))
		case int:
			array.Add(NewValueInt(value))
		case int8:
			array.Add(NewValueInt8(value))
		case int16:
			array.Add(NewValueInt16(value))
		case int32:
			array.Add(NewValueInt32(value))
		case int64:
			array.Add(NewValueInt64(value))
		case uint:
			array.Add(NewValueUint(value))
		case uint8:
			array.Add(NewValueUint8(value))
		case uint16:
			array.Add(NewValueUint16(value))
		case uint32:
			array.Add(NewValueUint32(value))
		case uint64:
			array.Add(NewValueUint64(value))
		case bool:
			array.Add(NewValueBool(value))
		default:
			panic(fmt.Sprintf("cannot use {%+v}[%T] as loukoum Value", value, value))
		}
	}

	return array
}
