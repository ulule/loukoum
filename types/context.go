package types

import (
	"bytes"
	"fmt"
)

// A Context is passed to a root stmt.Statement to generate a query.
type Context struct {
	buffer bytes.Buffer
	values map[string]interface{}
}

// Write appends given subquery in context's buffer.
func (ctx *Context) Write(query string) {
	_, err := ctx.buffer.WriteString(query)
	if err != nil {
		panic("loukoum: cannot write on buffer")
	}
}

// Bind adds given value in context's values.
func (ctx *Context) Bind(value interface{}) {
	idx := len(ctx.values) + 1
	key := fmt.Sprintf("arg_%d", idx)
	ctx.values[key] = value
	ctx.Write(":")
	ctx.Write(key)
}

// Query returns the underlaying query.
func (ctx *Context) Query() string {
	return ctx.buffer.String()
}

// Values returns the underlaying values of the query.
func (ctx *Context) Values() map[string]interface{} {
	return ctx.values
}

// NewContext creates a new Context instance.
func NewContext() *Context {
	return &Context{
		values: make(map[string]interface{}),
	}
}
