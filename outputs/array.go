package outputs

func Array[T any](first Output[T], rest ...Output[T]) Output[[]T] {
	all := append([]Output[T]{first}, rest...)
	return All(OutputContext(first), all)
}
