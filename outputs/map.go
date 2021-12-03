package outputs

type Map[K comparable, V any] map[K]Output[V]

func (m Map[K,V]) ToOutput() Output[map[K]V] {
	inner := map[K]Output[V](m)

	var ctx *Context
	for _, v := range inner {
		ctx = OutputContext(v)
		break
	}

	if ctx == nil {
		panic("Map should have at least 1 element")
	}

	return Apply(ctx, func (t *T) map[K]V {
		result := make(map[K]V)
		for k, v := range inner {
			result[k] = Eval(t, v)
		}
		return result
	})
}

var _ Output[map[string]string] = Map[string,string]{}
