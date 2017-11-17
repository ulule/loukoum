package types

import (
	"bytes"
	"fmt"
)

type Context struct {
	buffer bytes.Buffer
	values map[string]interface{}
}

func (ctx *Context) Write(query string) {
	ctx.buffer.WriteString(query)
}

func (ctx *Context) Bind(value interface{}) {
	idx := len(ctx.values) + 1
	key := fmt.Sprintf(":arg_%d", idx)
	ctx.values[key] = value
	ctx.Write(key)
}

func (ctx *Context) Query() string {
	return ctx.buffer.String()
}

func (ctx *Context) Values() map[string]interface{} {
	return ctx.values
}

func NewContext() *Context {
	return &Context{
		values: make(map[string]interface{}),
	}
}
