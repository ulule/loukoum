package types

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ulule/loukoum/v3/format"
)

// A Context is passed to a root stmt.Statement to generate a query.
type Context interface {
	Write(query string)
	Bind(value interface{})
}

// RawContext embeds values directly in the query.
type RawContext struct {
	buffer strings.Builder
}

// NewRawContext returns a new RawContext instance.
func NewRawContext() *RawContext {
	ctx := poolRawContext.Get().(*RawContext)
	return ctx
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

// Close returns instance to memory pool to reduce pressure on the garbage collector.
func (ctx *RawContext) Close() {
	if ctx != nil && ctx.buffer.Cap() < (1<<16) {
		ctx.buffer.Reset()
		poolRawContext.Put(ctx)
	}
}

// NamedContext uses named query placeholders.
type NamedContext struct {
	RawContext
	values map[string]interface{}
}

// NewNamedContext returns a new NamedContext instance.
func NewNamedContext() *NamedContext {
	ctx := poolNamedContext.Get().(*NamedContext)
	return ctx
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

// Values returns the named argument values.
func (ctx *NamedContext) Values() map[string]interface{} {
	return ctx.values
}

// Close returns instance to memory pool to reduce pressure on the garbage collector.
func (ctx *NamedContext) Close() {
	if ctx != nil && len(ctx.values) < (1<<12) {
		for k := range ctx.values {
			delete(ctx.values, k)
		}
		poolNamedContext.Put(ctx)
	}
}

// StdContext uses positional query placeholders.
type StdContext struct {
	RawContext
	values []interface{}
}

// NewStdContext returns a new StdContext instance.
func NewStdContext() *StdContext {
	ctx := poolStdContext.Get().(*StdContext)
	return ctx
}

// Bind adds given value in context's values.
func (ctx *StdContext) Bind(value interface{}) {
	idx := len(ctx.values) + 1
	ctx.values = append(ctx.values, value)
	ctx.Write(fmt.Sprintf("$%d", idx))
}

// Values returns the positional argument values.
func (ctx *StdContext) Values() []interface{} {
	return ctx.values
}

// Close returns instance to memory pool to reduce pressure on the garbage collector.
func (ctx *StdContext) Close() {
	if ctx != nil && cap(ctx.values) < (1<<12) {
		ctx.values = ctx.values[:0]
		poolStdContext.Put(ctx)
	}
}

var poolRawContext = sync.Pool{
	New: func() interface{} {
		return &RawContext{}
	},
}

var poolNamedContext = sync.Pool{
	New: func() interface{} {
		return &NamedContext{
			values: make(map[string]interface{}, 256),
		}
	},
}

var poolStdContext = sync.Pool{
	New: func() interface{} {
		return &StdContext{
			values: make([]interface{}, 0, 256),
		}
	},
}
