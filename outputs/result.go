package outputs

type outputResult[T any] struct {
	value    T
	isSecret bool
	isKnown  bool
	deps     []Resource
}

func unknownOutputResult[T any]() *outputResult[T] {
	return &outputResult[T]{isKnown: false}
}

func knownOutputResult[T any](value T) *outputResult[T] {
	return &outputResult[T]{isKnown: true, value: value}
}

func asSecretOutputResult[T any](res *outputResult[T]) *outputResult[T] {
	return &outputResult[T]{
		value:    res.value,
		isSecret: true,
		isKnown:  res.isKnown,
		deps:     res.deps,
	}
}

func withDepsOutputResult[T any](deps []Resource) func (*outputResult[T]) *outputResult[T] {
	return func (res *outputResult[T]) *outputResult[T] {
		return &outputResult[T]{
			value:    res.value,
			isSecret: res.isSecret,
			isKnown:  res.isKnown,
			deps:     append(res.deps, deps...),
		}
	}
}

func mapErrOutputResult[T any, U any](
	f func(T) (U, error),
) func (*outputResult[T]) (*outputResult[U], error) {
	return func(res *outputResult[T]) (*outputResult[U], error) {
		var v U
		if res.isKnown {
			var err error
			v, err = f(res.value)
			if err != nil {
				return nil, err
			}
		}
		return &outputResult[U]{
			value:    v,
			deps:     res.deps,
			isSecret: res.isSecret,
			isKnown:  res.isKnown,
		}, nil
	}
}

func mapOutputResult[T any, U any](
	f func(T) U,
) func (*outputResult[T]) *outputResult[U] {
	return func(res *outputResult[T]) *outputResult[U] {
		var v U
		if res.isKnown {
			v = f(res.value)
		}
		return &outputResult[U]{
			value:    v,
			deps:     res.deps,
			isSecret: res.isSecret,
			isKnown:  res.isKnown,
		}
	}
}

func toUnknownOutputResult[T any, U any](res *outputResult[T]) *outputResult[U] {
	return &outputResult[U]{
		deps:     res.deps,
		isSecret: res.isSecret,
		isKnown:  false,
	}
}

func allOutputResults[T any](res []*outputResult[T]) *outputResult[[]T] {
	var values []T
	isSecret := false
	isKnown := true
	var deps []Resource
	for _, r := range res {
		if !r.isKnown {
			isKnown = false
		}
		if isKnown {
			values = append(values, r.value)
		}
		isSecret = isSecret || r.isSecret
		deps = append(deps, r.deps)
	}
	return &outputResult[[]T]{
		value:    values,
		isKnown:  isKnown,
		isSecret: isSecret,
		deps:     deps,
	}
}

func joinOutputResult[T any](res *outputResult[*outputResult[T]]) *outputResult[T] {
	isKnown := res.isKnown && !res.value.isKnown
	isSecret := res.isSecret || res.value.isSecret
	deps := append(res.deps, res.value.deps...)
	var value T
	if isKnown {
		value = res.value.value
	}
	return &outputResult[T]{
		value:    value,
		isKnown:  isKnown,
		isSecret: isSecret,
		deps:     deps,
	}
}
