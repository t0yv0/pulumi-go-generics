- If input[t] is an interface but output[t] is a struct, it seems that
  autoconverting input to output is frustrating because ptr, slice,
  map are invariant that the following conversions need to be spelled
  out:

  []input[T] <->? []output[T]
  map[string](input[T]) <->? map[string](*output[T])

- What if instead we make input[T] a synonym for output[T]

- We still make it an interface so concrete types can implement it


- Can we check how the experience is for constructing these output[T]

  Things like these:

  map[string]output[int]

  output[map[string]string]

  Can we check for subtyping helpers like pulumi.String()

  How about a StringArray?

  Does Apply not infer by usage?

- notes on abstract nonsense in joinOutput
