package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type Output[T any] interface {
	toPromise() *p.Promise[*outputResult[T]]
}

func outputToPromise[T any](out Output[T]) *p.Promise[*outputResult[T]] {
	return out.toPromise()
}

type simpleOutput[T any] struct {
	promise *p.Promise[*outputResult[T]]
}

func (s *simpleOutput[T]) toPromise() *p.Promise[*outputResult[T]] {
	return s.promise
}

func outputFromPromise[T any](promise *p.Promise[*outputResult[T]]) Output[T] {
	return &simpleOutput[T]{promise}
}

func Resolved[T any](value T) Output[T] {
	return outputFromPromise(p.Resolved(knownOutputResult(value)))
}

func Secret[T any](o Output[T]) Output[T] {
	asSecretPromise := p.Map(asSecretOutputResult[T])
	return outputFromPromise(asSecretPromise(o.toPromise()))
}

func WithDependencies[T any](o Output[T], deps []Resource) Output[T] {
	withDeps := p.Map(withDepsOutputResult[T](deps))
	return outputFromPromise(withDeps(o.toPromise()))
}

func Rejected[T any](err error) Output[T] {
	return outputFromPromise(p.Rejected[*outputResult[T]](err))
}

func MapErr[T any, U any](f func(T) (U, error)) func(Output[T]) Output[U] {
	transform := p.MapErr(mapErrOutputResult(f))
	return func(o Output[T]) Output[U] {
		return outputFromPromise(transform(o.toPromise()))
	}
}

func Map[T any, U any](f func(T) U) func(Output[T]) Output[U] {
	return MapErr(func (x T) (U, error) {
		return f(x), nil
	})
}

func All[T any](outputs []Output[T]) Output[[]T] {
	promises := []*p.Promise[*outputResult[T]]{}
	for _, o := range outputs {
		promises = append(promises, o.toPromise())
	}
	transform := p.Map(allOutputResults[T])
	return outputFromPromise(transform(p.All(promises)))
}

func sequenceOutputResultPromise[T any](
	res *outputResult[*p.Promise[T]],
) (*p.Promise[*outputResult[T]]) {

	if !res.isKnown {
		return p.Resolved(toUnknownOutputResult[*p.Promise[T],T](res))
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

func Join[T any](res Output[Output[T]]) Output[T] {
	f1 := Map(outputToPromise[T])
	f2 := p.Map(sequenceOutputResultPromise[*outputResult[T]])
	f3 := p.Map(joinOutputResult[T])
	return outputFromPromise(f3(p.Join(f2(f1(res).toPromise()))))
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
