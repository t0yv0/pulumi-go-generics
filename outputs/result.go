package outputs

type result[T any] struct {
	value T
	info  info
}

func knownResult[T any](v T) result[T] {
	return result[T]{v, info{isKnown: true}}
}

func unknownResult[T any]() result[T] {
	return result[T]{info: info{isKnown: false}}
}

func secretResult[T any](res result[T]) result[T] {
	res.info.isSecret = true
	return res
}

func withDepsResult[T any](res result[T], deps []Resource) result[T] {
	res.info.deps = append(res.info.deps, deps...)
	return res
}
