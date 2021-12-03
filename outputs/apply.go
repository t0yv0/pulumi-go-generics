package outputs

import (
	p "github.com/t0yv0/pulumi-go-generics/promises"
)

type T struct {
	info info
}

func ApplyErr[A any](ctx *Context, body func(t *T) (A, error)) Output[A] {
	t := &T{}

	promise, resolve, reject := p.NewPromise[result[A]](ctx)

	go func() {
		v, err := body(t)
		if err != nil {
			reject(err)
		}
		resolve(result[A]{
			value: v,
			info:  t.info,
		})
	}()
	return &out[A]{ctx, promise}
}

func Apply[A any](ctx *Context, body func(t *T) A) Output[A] {
	return ApplyErr(ctx, func (t *T) (A, error) {
		return body(t), nil
	})
}

func Eval[A any](t *T, out Output[A]) A {
	p := toPromise(out)
	v, err := p.Await()
	if err != nil {
		panic(err) // TODO recover this in ApplyErr
	}

	t.info = infos(t.info, v.info)
	return v.value
}

func EvalErr[A any](t *T, out Output[A]) (A, error) {
	p := toPromise(out)
	v, err := p.Await()
	if err != nil {
		return v.value, err
	}

	t.info = infos(t.info, v.info)
	return v.value, nil
}
