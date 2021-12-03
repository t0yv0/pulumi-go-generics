# Go Generics experimentation for Pulumi

Ideation on how to utilize Go generics to design a better Go SDK
experience for Pulumi users.

## Building

This requires a go1.18 toolchain that is currently pre-release.

It can be built from source directly or via a helper Nix wrapper:

https://github.com/t0yv0/pulumi-nix-devenv

    $ cd pulumi-nix-devenv
    $ nix-shell devenv-with-go-dev.nix
    $ cd pulumi-go-generics
    $ go build .

## Design

### Problem: Input and Output types

Having separate `Input[T]` and `Output[T]` has issues as Pat Gavlin
pointed out.

1. `map[string]Input[T]` and `map[string]Output[T]` do not safely
   convert (more generally, types do not convert with I/O in nested
   positions)

2. arguments typed `Input[T]` fail to infer `T` at call site when
   passed an `Output[T]`

3. generic aliases `type Input[T] = Output[T]` are prohibited

See `examples4.go`.

### Solution: eliminate Input type entirey

Similar to what @frassle suggested for .NET the idea is to drop
`Input[T]` type entirely. If Input exists for dev experience and not
fundamental needs, it may be removed in the name of dev experience. If
we do need some fundamentals like preserving promptness through
combinators and inspecting for promptness, this can be added to the
Output type implementation.

### Problem: lack of generic methods complicates apply chains

Some languages are happy to infer types and make the syntax of apply
chains easy as in:

    x.apply(y => y.apply(z => z.apply(f)))

With Go Generics, the options available seem to be too baroque. Apply
cannot be a method (needs to be a top-level func).

### Solution: block-style apply with eval

    func Apply[A any](
        ctx *Context,
        body func(t *T) A
    ) Output[A]

    func Eval[A any](
        t *T,
        out Output[A],
    ) A

    Apply(ctx, func(t *T) ZType {
        y := Eval(t, x)
        z := Eval(t, y)
        return f(z)
    })

More in `examples3.go`.

### Problem: issues with All

There is clunkiness around variations on the `All` theme:

    All([]Output[T]) Output[[]T]

Variants like Pair, Triple, Tuple4 have difficulty expressing the
explosion of generics. Boxing/unboxing is unreadable.

### Solution: block-style apply

Dedicated `All` is simply a sugar over Apply. See `all.go`. Tuple
forms are omitted, use block Apply. See `examples.go`

### Problem: no lifted property access

In some languages we support equivalence of the two forms:

    x.apply(xv => xv.foo)
    x.foo

Previously solved by helper types with methods.

### Solution: block-style apply again

Block-form Apply and Eval get the user back to manipulating raw
values, so for complex examples lifted property access loses some of
the motivation. Possibly can consider helpers for common boilerplate
scenarios.
