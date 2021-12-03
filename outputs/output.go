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
	return &out[T]{ctx, p.Resolved(ctx, knownOutputResult(value))}
}

func Unknown[T any](ctx *Context) Output[T] {
	return &out[T]{ctx, p.Resolved(ctx, unknownOutputResult[T]())}
}

func Rejected[T any](ctx *Context, err error) Output[T] {
	return &out[T]{ctx, p.Rejected[*outputResult[T]](ctx, err)}
}

func Secret[T any](o Output[T]) Output[T] {
	asSecretPromise := p.Map(asSecretOutputResult[T])
	return &out[T]{context(o), asSecretPromise(toPromise(o))}
}

func WithDependencies[T any](o Output[T], deps ...Resource) Output[T] {
	withDeps := p.Map(withDepsOutputResult[T](deps))
	return &out[T]{context(o), withDeps(toPromise(o))}
}
