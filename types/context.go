package types

import (
	"bytes"
	"fmt"
)

// A Context is passed to a root stmt.Statement to generate a query.
type Context struct {

	// Prefix added to query placeholders
	Prefix string

	buffer bytes.Buffer
	values map[string]interface{}
	args   []interface{}
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
	idx := len(ctx.args) + 1
	name := fmt.Sprintf("arg_%d", idx)
	ctx.values[name] = value
	ctx.args = append(ctx.args, value)
	ctx.Write(fmt.Sprintf("%s%d", ctx.Prefix, idx))
}

// Query returns the underlaying query.
func (ctx *Context) Query() string {
	return ctx.buffer.String()
}

// Values returns the underlaying values of the query.
func (ctx *Context) Values() map[string]interface{} {
	return ctx.values
}

// Args return the underlying values of the query as a slice.
func (ctx *Context) Args() []interface{} {
	return ctx.args
}

// Placeholder prefixes
const (
	NamedPrefix    = ":arg_"
	PostgresPrefix = "$"
)

// NewContext creates a new Context instance.
func NewContext() *Context {
	return &Context{
		Prefix: NamedPrefix,
		values: make(map[string]interface{}),
	}
}
