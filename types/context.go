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
	key1 := fmt.Sprintf("arg_%d", idx)
	key2 := fmt.Sprint(":", key1)
	ctx.values[key1] = value
	ctx.Write(key2)
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
