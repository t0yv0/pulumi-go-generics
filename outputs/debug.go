package outputs

import (
	"fmt"
)

func Debug[T any](x Output[T]) {
	v, err := toPromise(x).Await()
	if err != nil {
		fmt.Printf("outputs.Debug: failed with %v\n", err)
		return
	}
	fmt.Printf("outputs.Debug: succeeded\n")
	fmt.Printf("    value: %v\n", v.value)
	fmt.Printf("    isKnown: %v\n", !v.info.isUnknown)
	fmt.Printf("    isSecret: %v\n", v.info.isSecret)
	fmt.Printf("    deps: %d\n", len(v.info.deps))
	for _, d := range v.info.deps {
		fmt.Printf("        %v", d)
	}
}
