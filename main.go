package main

import (
	pulumi "github.com/t0yv0/pulumi-go-generics/outputs"
)

func main() {
	ctx := pulumi.NewContext()

	Examples(ctx)
	Examples2(ctx)
	Examples3(ctx)
}
