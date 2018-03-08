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
	Query() string
}

// StringContext embeds values directly in the query.
type StringContext struct {
	buffer bytes.Buffer
	values map[string]interface{}
}

// Write appends given subquery in context's buffer.
func (ctx *StringContext) Write(query string) {
	_, err := ctx.buffer.WriteString(query)
	if err != nil {
		panic("loukoum: cannot write on buffer")
	}
}

// Bind adds given value in context's values.
func (ctx *StringContext) Bind(value interface{}) {
	ctx.Write(format.Value(value))
}

// Query returns the underlaying query.
func (ctx *StringContext) Query() string {
	return ctx.buffer.String()
}

// NamedContext uses named query placeholders
type NamedContext struct {
	StringContext
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

// StdContext uses positional query placeholders
type StdContext struct {
	StringContext
	values []interface{}
}

// Bind adds give value in context's values.
func (ctx *StdContext) Bind(value interface{}) {
	idx := len(ctx.values) + 1
	ctx.values = append(ctx.values, value)
	ctx.Write(fmt.Sprintf("$%d", idx))
}
