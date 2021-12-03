package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type Map[K comparable, V any] map[K]Output[V]

func (m Map[K,V]) Context() *Context {
	panic("TODO")
}

func (m Map[K,V]) toPromise() *p.Promise[*outputResult[map[K]V]] {
	panic("TODO")
}
