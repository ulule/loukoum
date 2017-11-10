package stmt

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/ulule/loukoum/types"
)

type Expression interface {
	Statement
	expression()
}

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
	default:
		return NewValue(fmt.Sprint(value))
	}
}

// TODO Refacto ?
type Identifier struct {
	Identifier string
}

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

func (identifier Identifier) Equal(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Equal)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) NotEqual(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.NotEqual)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) Is(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Is)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) IsNot(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.IsNot)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) GreaterThan(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.GreaterThan)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) GreaterThanOrEqual(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.GreaterThanOrEqual)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) LessThan(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.LessThan)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) LessThanOrEqual(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.LessThanOrEqual)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) In(values []interface{}) InfixExpression {
	panic("TODO")
}

func (identifier Identifier) NotIn(values []interface{}) InfixExpression {
	panic("TODO")
}

func (identifier Identifier) Like(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.Like)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) NotLike(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.NotLike)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) ILike(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.ILike)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) NotILike(value interface{}) InfixExpression {
	operator := NewComparisonOperator(types.NotILike)
	return NewInfixExpression(identifier, operator, NewExpression(value))
}

func (identifier Identifier) Between(from, to interface{}) Between {
	return NewBetween(identifier, NewExpression(from), NewExpression(to))
}

func (identifier Identifier) NotBetween(from, to interface{}) Between {
	return NewNotBetween(identifier, NewExpression(from), NewExpression(to))
}

// TODO Refacto ?
type Value struct {
	Value string
}

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
