package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type Output[T any] interface {
	context()   *Context
	toPromise() *p.Promise[*outputResult[T]]
}

type out[T any] struct {
	ctx     *Context
	promise *p.Promise[*outputResult[T]]
}

func (s *out[T]) context() *Context {
	return s.ctx
}

func (s *out[T]) toPromise() *p.Promise[*outputResult[T]] {
	return s.promise
}

func Resolved[T any](ctx *Context, value T) Output[T] {
	return &out[T]{ctx, p.Resolved(ctx, knownOutputResult(value))}
}

func Secret[T any](o Output[T]) Output[T] {
	asSecretPromise := p.Map(asSecretOutputResult[T])
	return &out[T]{o.context(), asSecretPromise(o.toPromise())}
}

func WithDependencies[T any](o Output[T], deps []Resource) Output[T] {
	withDeps := p.Map(withDepsOutputResult[T](deps))
	return &out[T]{o.context(), withDeps(o.toPromise())}
}

func Rejected[T any](ctx *Context, err error) Output[T] {
	return &out[T]{ctx, p.Rejected[*outputResult[T]](ctx, err)}
}

func MapErr[T any, U any](f func(T) (U, error)) func(Output[T]) Output[U] {
	transform := p.MapErr(mapErrOutputResult(f))
	return func(o Output[T]) Output[U] {
		return &out[U]{o.context(), transform(o.toPromise())}
	}
}

func Map[T any, U any](f func(T) U) func(Output[T]) Output[U] {
	return MapErr(func (x T) (U, error) {
		return f(x), nil
	})
}

func All[T any](ctx *Context, outputs []Output[T]) Output[[]T] {
	promises := []*p.Promise[*outputResult[T]]{}
	for _, o := range outputs {
		promises = append(promises, o.toPromise())
	}
	transform := p.Map(allOutputResults[T])
	return &out[[]T]{ctx, transform(p.All(ctx, promises))}
}

func sequenceOutputResultPromise[T any](
	ctx *Context,
) func(*outputResult[*p.Promise[T]]) *p.Promise[*outputResult[T]] {
	return func(
		res *outputResult[*p.Promise[T]],
	) *p.Promise[*outputResult[T]] {
		if !res.isKnown {
			unk := toUnknownOutputResult[*p.Promise[T],T](res)
			return p.Resolved(ctx, unk)
		}

		f1 := func(value T) *outputResult[T] {
			g := mapOutputResult(func (*p.Promise[T]) T {
				return value
			})
			return g(res)
		}

		f2 := p.Map(f1)
		return f2(res.value)
	}
}

func Join[T any](res Output[Output[T]]) Output[T] {
	ctx := res.context()
	f1 := Map(func (o Output[T]) *p.Promise[*outputResult[T]] {
		return o.toPromise()
	})
	f2 := p.Map(sequenceOutputResultPromise[*outputResult[T]](ctx))
	f3 := p.Map(joinOutputResult[T])
	return &out[T]{ctx, f3(p.Join(f2(f1(res).toPromise())))}
}

func BindErr[T any, U any](
	out Output[T],
	cont func (T)(Output[U], error),
) Output[U] {
	return Join(MapErr(cont)(out))
}

func Bind[T any, U any](
	out Output[T],
	cont func (T)(Output[U]),
) Output[U] {
	return Join(Map(cont)(out))
}
