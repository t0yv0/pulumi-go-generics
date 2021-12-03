package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type Output[T any] interface {
	// When implementing the interface on a struct, return an
	// output of another type. When consuming outputs, the
	// framework will continue normalizing until it reaches a
	// terminal internal representation.
	ToOutput() Output[T]
}

func Resolved[T any](ctx *Context, value T) Output[T] {
	return &out[T]{ctx, p.Resolved(ctx, knownResult[T](value))}
}

func Unknown[T any](ctx *Context) Output[T] {
	return &out[T]{ctx, p.Resolved(ctx, unknownResult[T]())}
}

func Rejected[T any](ctx *Context, err error) Output[T] {
	return &out[T]{ctx, p.Rejected[result[T]](ctx, err)}
}

func Secret[A any](o Output[A]) Output[A] {
	return Apply(OutputContext(o), func (t *T) A {
		t.Secret()
		return Eval(t, o)
	})
}

func WithDependencies[A any](o Output[A], deps ...Resource) Output[A] {
	return Apply(OutputContext(o), func (t *T) A {
		t.DependsOn(deps...)
		return Eval(t, o)
	})
}
