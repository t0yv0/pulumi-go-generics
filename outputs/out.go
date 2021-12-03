package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type out[T any] struct {
	ctx     *Context
	promise *p.Promise[result[T]]
}

func (s *out[T]) ToOutput() Output[T] {
	return s
}

func normalize[T any](o Output[T]) *out[T] {
	if o == nil {
		panic("nil?")
	}
	for {
		if r, ok := o.(*out[T]); ok {
			return r
		}
		next := o.ToOutput()
		if next == o {
			panic("ToOutput should not return self")
		}
		o = next
	}
}

func toPromise[T any](o Output[T]) *p.Promise[result[T]] {
	return normalize(o).promise
}

func OutputContext[T any](o Output[T]) *Context {
	return normalize(o).ctx
}
