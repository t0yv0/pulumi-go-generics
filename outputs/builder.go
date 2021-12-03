package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type T struct {
	result *outputResult[interface{}]
}

func Eval[A any](t *T, out Output[A]) (A, error) {
	p := out.toPromise()
	v, err := p.Await()
	if err != nil {
		return v.value, err
	}

	t.result = firstOutputResult(t.result, v)
	return v.value, nil
}

func Apply[A any](ctx *Context, body func(t *T) (A, error)) Output[A] {
	t := &T{knownOutputResult[interface{}](nil)}

	promise, resolve, reject := p.NewPromise[*outputResult[A]](ctx)

	go func() {
		v, err := body(t)
		if err != nil {
			reject(err)
		}

		resolve(firstOutputResult(knownOutputResult(v), t.result))
	}()
	return &out[A]{ctx, promise}
}
