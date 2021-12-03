package outputs

func All[A any](ctx *Context, outputs []Output[A]) Output[[]A] {
	return ApplyErr(ctx, func (t *T) ([]A, error) {
		var results []A
		for _, o := range outputs {
			r, err := EvalErr(t, o)
			if err != nil {
				return nil, err
			}
			results = append(results, r)
		}
		return results, nil
	})
}
