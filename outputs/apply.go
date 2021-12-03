package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type T struct {
	result *outputResult[interface{}]
}

func ApplyErr[A any](ctx *Context, body func(t *T) (A, error)) Output[A] {
	t := &T{knownOutputResult[interface{}](nil)}

	promise, resolve, reject := p.NewPromise[*outputResult[A]](ctx)

	go func() {
		v, err := body(t)
		if err != nil {
			reject(err)
		}
		res := firstOutputResult(knownOutputResult(v), t.result)
		resolve(res)
	}()
	return &out[A]{ctx, promise}
}

func Apply[A any](ctx *Context, body func(t *T) A) Output[A] {
	return ApplyErr(ctx, func (t *T) (A, error) {
		return body(t), nil
	})
}

func Eval[A any](t *T, out Output[A]) A {
	p := out.toPromise()
	v, err := p.Await()
	if err != nil {
		panic(err) // TODO recover this in ApplyErr
	}

	t.result = firstOutputResult(t.result, v)
	return v.value
}

func EvalErr[A any](t *T, out Output[A]) (A, error) {
	p := out.toPromise()
	v, err := p.Await()
	if err != nil {
		return v.value, err
	}

	t.result = firstOutputResult(t.result, v)
	return v.value, nil
}
