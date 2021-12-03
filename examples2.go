// Inspired by:
//
// https://github.com/pgavlin/go-generics/blob/701df2cf62a492829a23c4fe516b5594d77add82/output/main.go

package main

import (
	//"fmt"
	//"time"

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
	return pulumi.Apply(s.Foo.Context(), func(t *pulumi.T) Struct {
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

	hello := pulumi.Unknown[string](ctx)
	world := pulumi.Unknown[string](ctx)
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

	// // Note that the type of the Map is Map[string, map[string]string], not Map[string, Map[string]]. This is because
	// // Map[K, V]'s underlying type is map[K]Input[V], so V must be the _resolved_ type of the Map element.
	// var in2 Input[map[string]map[string]string] = Map[string, map[string]string]{
	// 	"goodbye": Map[string, string]{
	// 		"cruel": world,
	// 	},
	// }

	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	hello.resolve("hello", true, false, nil)
	// 	world.resolve("world", true, false, nil)
	// }()

	// // Type inference fails here in the same way it fails for Ptr[] above
	// apply := Apply[Struct](in.ToOutput(), func(v Struct) (string, error) {
	// 	return fmt.Sprintf("%#v", v), nil
	// })

	// v, _, _, _, _ := apply.await()
	// fmt.Printf("%v\n", v)

	// v2, _, _, _, _ := in2.ToOutput().await()
	// fmt.Printf("%#v\n", v2)

	// // The compiler cannot infer T as string b/c the parameters to All are Input[T] rather than *Output[T]
	// v3, _, _, _, _ := All[string](hello, world).await()
	// fmt.Printf("%v %v\n", v3[0], v3[1])

	// fmt.Printf("\n\nExamples2\n")
}
