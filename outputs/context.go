package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type Context struct {
	p.Waiter
}

var _ p.Observer = &Context{}

func NewContext() *Context {
	w := p.NewWaiter()
	return &Context{*w}
}

// helpers

func (ctx *Context) Int(n int) Output[int] {
	return Resolved(ctx, n)
}

func (ctx *Context) String(s string) Output[string] {
	return Resolved(ctx, s)
}
