package main

import (
	"fmt"
	"time"

	pulumi "github.com/t0yv0/pulumi-go-generics/outputs"
)

func Examples3(ctx *pulumi.Context) {
	fmt.Printf("\n\nExamples3\n")

	// Prompt values lift to Output with generic Resolved
	var stringO pulumi.Output[string] = pulumi.Resolved(ctx, "foo")
	pulumi.Debug(stringO)

	// Note that Resolved infers the type param
	var int1 pulumi.Output[int] = pulumi.Resolved(ctx, 1)
	pulumi.Debug(int1)

	// Mapping over an output uses Apply builder form.
	// Inside the body of the builder, use Eval to peek into an Output.
	// To enforce the convention Eval requires a *T that Apply provides.
	// Its name is short, similar to the conventional *testing.T.
	var int2 pulumi.Output[int] = pulumi.Apply(ctx, func(t *pulumi.T) int {
		return pulumi.Eval(t, int1) + 1
	})
	pulumi.Debug(int2)

	// We can combine more than one output, similar to Map2, Map3
	// etc with the exact same form. Note that the form evaluates
	// sequentially but it does not necessarily matter since
	// output values allocated outside the form already have
	// active goroutines trying to resolve their promises in
	// parallel. So awaiting them sequentially is not impeding
	// parallelism.
	joinO := pulumi.Apply(ctx, func(t *pulumi.T) string {
		s := pulumi.Eval(t, stringO)
		i := pulumi.Eval(t, int1)
		res := fmt.Sprintf("stringO=%s intO=%d", s, i)
		return res
	})
	pulumi.Debug(joinO)

	// We can combine dependent chains of outputs using the same
	// exact form. Simply Eval intermediate values in the builder
	// sequentially. To make it curious let's introduce a function
	// that makes the output as secret.
	f := func (x int) pulumi.Output[int] {
		return pulumi.Apply(ctx, func(t *pulumi.T) int {
			secret := pulumi.Eval(t, pulumi.Secret(pulumi.Resolved(ctx, 1)))
			time.Sleep(100)
			return secret + x + 1
		})
	}

	// Now `chain0` becomes secret through the use of `f`.
	chainO := pulumi.Apply(ctx, func(t *pulumi.T) int {
		i1 := pulumi.Eval(t, int1)
		i2 := pulumi.Eval(t, f(i1))
		i3 := pulumi.Eval(t, f(i2))
		return i3
	})
	pulumi.Debug(chainO)

	// Curiously Resolved is expressible in the form of Apply also:
	just1 := pulumi.Apply(ctx, func (*pulumi.T) int {
		return 1
	})
	pulumi.Debug(just1)

	// All is expressible as Apply alspulumi. So is All/Map such as
	// taking a sum. It's just as parallel using the same logic
	// as Map2/Map3 etc.
	xs := []pulumi.Output[int]{int1, int1, f(2)}
	var xsO pulumi.Output[[]int] = pulumi.Apply(ctx, func(t *pulumi.T) []int {
		var results []int
		for _, x := range xs {
			results = append(results, pulumi.Eval(t, x))
		}
		return results
	})
	pulumi.Debug(xsO)
}
