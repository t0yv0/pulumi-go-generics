package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type Context struct {
	p.Waiter
}

var _ p.Observer = &Context{}
