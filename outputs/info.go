package outputs

type info struct {
	isSecret bool
	isKnown  bool
	deps     []Resource
}

func infos(infos ...info) info {
	i := info{}
	for _, x := range infos {
		i.isSecret = i.isSecret || x.isSecret
		i.isKnown = i.isKnown && x.isKnown
		i.deps = append(i.deps, x.deps...)
	}
	return i
}
