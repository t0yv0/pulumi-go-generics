package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type Output[T any] interface {
	Context()   *Context
	toPromise() *p.Promise[*outputResult[T]]
}

type out[T any] struct {
	ctx     *Context
	promise *p.Promise[*outputResult[T]]
}

func (s *out[T]) Context() *Context {
	return s.ctx
}

func (s *out[T]) toPromise() *p.Promise[*outputResult[T]] {
	return s.promise
}

func Unknown[T any](ctx *Context) Output[T] {
	return &out[T]{ctx, p.Resolved(ctx, unknownOutputResult[T]())}
}

func Resolved[T any](ctx *Context, value T) Output[T] {
	return &out[T]{ctx, p.Resolved(ctx, knownOutputResult(value))}
}

func Secret[T any](o Output[T]) Output[T] {
	asSecretPromise := p.Map(asSecretOutputResult[T])
	return &out[T]{o.Context(), asSecretPromise(o.toPromise())}
}

func WithDependencies[T any](o Output[T], deps []Resource) Output[T] {
	withDeps := p.Map(withDepsOutputResult[T](deps))
	return &out[T]{o.Context(), withDeps(o.toPromise())}
}

func Rejected[T any](ctx *Context, err error) Output[T] {
	return &out[T]{ctx, p.Rejected[*outputResult[T]](ctx, err)}
}
