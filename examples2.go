// Inspired by:
//
// https://github.com/pgavlin/go-generics/blob/701df2cf62a492829a23c4fe516b5594d77add82/output/main.go

package main

import (
	"fmt"
	"time"

	pulumi "github.com/t0yv0/pulumi-go-generics/outputs"
)

type Struct struct {
	Foo int               `pulumi:"foo"`
	Bar string            `pulumi:"bar"`
	Baz []string          `pulumi:"baz"`
	Qux map[string]string `pulumi:"qux"`
	Zed *Struct           `pulumi:"zed"`
}

type StructArgs struct {
	Foo pulumi.Output[int]
	Bar pulumi.Output[string]
	Baz pulumi.Output[[]string]
	Qux pulumi.Output[map[string]string]
	Zed pulumi.Output[*Struct]
}

// This is presents direct non-reflective approach, though reflection
// may be preferable to cut the boilerplate.
func (s StructArgs) ToOutput() pulumi.Output[Struct] {
	return pulumi.Apply(pulumi.OutputContext(s.Foo), func(t *pulumi.T) Struct {
		return Struct{
			Foo: pulumi.Eval(t, s.Foo),
			Bar: pulumi.Eval(t, s.Bar),
			Baz: pulumi.Eval(t, s.Baz),
			Qux: pulumi.Eval(t, s.Qux),
			Zed: pulumi.Eval(t, s.Zed),
		}
	})
}


func Examples2(ctx *pulumi.Context) {
	fmt.Printf("\n\nExamples2\n")

	hello := pulumi.Apply(ctx, func(*pulumi.T) string {
		time.Sleep(1 * time.Second)
		return "hello"
	})

	world := pulumi.Apply(ctx, func(*pulumi.T) string {
		time.Sleep(1 * time.Second)
		return "world"
	})

	monde := ctx.String("monde cruel")

	var in pulumi.Output[Struct] = StructArgs{
		Foo: ctx.Int(42),
		Bar: hello,
		Baz: pulumi.Array(ctx.String("wonderful"), world),
		Qux: pulumi.Map[string, string]{
			"goodbye": ctx.String("cruel world"),
			"adieu":   monde,
		},
		Zed: pulumi.Resolved(ctx, &Struct{
			Foo: 24,
		}),
	}

	pulumi.Debug(in)

	// Note that the type of the Map is Map[string, map[string]string], not Map[string, Map[string]]. This is because
	// Map[K, V]'s underlying type is map[K]Input[V], so V must be the _resolved_ type of the Map element.
	var in2 pulumi.Output[map[string]map[string]string] = pulumi.Map[string, map[string]string]{
		"goodbye": pulumi.Map[string, string]{
			"cruel": world,
		},
	}
	pulumi.Debug(in2)

	// Type inference fails here in the same way it fails for Ptr[] above
	applied := pulumi.Apply(ctx, func(t *pulumi.T) string {
		v := pulumi.Eval(t, in)
		return fmt.Sprintf("%#v", v)
	})
	pulumi.Debug(applied)

	hw := pulumi.Apply(ctx, func(t *pulumi.T) string {
		return fmt.Sprintf("%v %v\n",
			pulumi.Eval(t, hello),
			pulumi.Eval(t, world))
	})
	pulumi.Debug(hw)
}
