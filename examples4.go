package main

func checks(
	outInt Out[int],
	outIntSlice Out[[]int],
	mapOutInt map[string]Out[int],
) In[int] {

	useIn(outInt)
	useInSlice(outIntSlice)

	// If the call site is generic, type inferences does not work.
	useInT[int](outInt)
	useInT[[]int](outIntSlice)

	// If equivalent In and Out types are nested inside a complex
	// type, the complex unification does not work *at all*,
	// requiring unsafe casts or runtime code to repack types.
	//
	// useMapIn(mapOutInt)
	//
	// --> cannot convert mapOutInt (variable of type map[string]Out[int])
	//     to type map[string]In[int]

	// Generic type aliases do not work.
	//
	// type In[T] = Out[T]
	// --> generic type cannot be alias

	panic("?")
}

type Out[T any] interface {
	isOut() bool
}

type In[T any] interface {
	Out[T]
}

func useIn(x In[int]) {
}

func useInSlice(x In[[]int]) {
}

func useInT[T any](x In[T]) {
}

func useMapIn(x map[string]In[int]) {
}
