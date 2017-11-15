package stmt

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

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
		return NewValue(value)
	case int:
		return NewValue(strconv.FormatInt(int64(value), 10))
	case int8:
		return NewValue(strconv.FormatInt(int64(value), 10))
	case int16:
		return NewValue(strconv.FormatInt(int64(value), 10))
	case int32:
		return NewValue(strconv.FormatInt(int64(value), 10))
	case int64:
		return NewValue(strconv.FormatInt(value, 10))
	case uint:
		return NewValue(strconv.FormatUint(uint64(value), 10))
	case uint8:
		return NewValue(strconv.FormatUint(uint64(value), 10))
	case uint16:
		return NewValue(strconv.FormatUint(uint64(value), 10))
	case uint32:
		return NewValue(strconv.FormatUint(uint64(value), 10))
	case uint64:
		return NewValue(strconv.FormatUint(value, 10))
	case []string, []int, []uint, []int8, []uint8, []int16, []uint16, []int32, []uint32, []int64, []uint64:
		return NewValue(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(value)), ", "), "[]"))
	default:
		return NewValue(fmt.Sprint(value))
	}
}

// Identifier is an identifier.
// TODO Refacto ?
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
func (identifier Identifier) In(value interface{}) In {
	return NewIn(identifier, NewExpression(value))
}

// NotIn performs a "not in" condition.
func (identifier Identifier) NotIn(value interface{}) In {
	return NewNotIn(identifier, NewExpression(value))
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

// Value is an expression value.
// TODO Refacto ?
type Value struct {
	Value string
}

// NewValue returns a an expression value.
func NewValue(value string) Value {
	return Value{
		Value: value,
	}
}

func (Value) expression() {}

func (value Value) Write(buffer *bytes.Buffer) {
	buffer.WriteString(value.Value)
}

// IsEmpty return true if statement is undefined.
func (value Value) IsEmpty() bool {
	return value.Value == ""
}
