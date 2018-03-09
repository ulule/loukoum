package types

import (
	"bytes"
	"fmt"

	"github.com/ulule/loukoum/format"
)

// A Context is passed to a root stmt.Statement to generate a query.
type Context interface {
	Write(query string)
	Bind(value interface{})
}

// RawContext embeds values directly in the query.
type RawContext struct {
	buffer bytes.Buffer
}

// Write appends given subquery in context's buffer.
func (ctx *RawContext) Write(query string) {
	_, err := ctx.buffer.WriteString(query)
	if err != nil {
		panic("loukoum: cannot write on buffer")
	}
}

// Bind adds given value in context's values.
func (ctx *RawContext) Bind(value interface{}) {
	ctx.Write(format.Value(value))
}

// Query returns the underlaying query.
func (ctx *RawContext) Query() string {
	return ctx.buffer.String()
}

// NamedContext uses named query placeholders
type NamedContext struct {
	RawContext
	values map[string]interface{}
}

// Bind adds given value in context's values.
func (ctx *NamedContext) Bind(value interface{}) {
	if ctx.values == nil {
		ctx.values = make(map[string]interface{})
	}
	idx := len(ctx.values) + 1
	name := fmt.Sprintf("arg_%d", idx)
	ctx.values[name] = value
	ctx.Write(":" + name)
}

// Values returns the named argument values
func (ctx *NamedContext) Values() map[string]interface{} {
	return ctx.values
}

// StdContext uses positional query placeholders
type StdContext struct {
	RawContext
	values []interface{}
}

// Bind adds give value in context's values.
func (ctx *StdContext) Bind(value interface{}) {
	idx := len(ctx.values) + 1
	ctx.values = append(ctx.values, value)
	ctx.Write(fmt.Sprintf("$%d", idx))
}

// Values returns the positional argument values
func (ctx *StdContext) Values() []interface{} {
	return ctx.values
}
