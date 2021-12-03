package outputs

type info struct {
	isSecret  bool
	isUnknown bool
	deps      []Resource
}

func infos(infos ...info) info {
	i := info{}
	for _, x := range infos {
		i.isSecret = i.isSecret || x.isSecret
		i.isUnknown = i.isUnknown || x.isUnknown
		i.deps = append(i.deps, x.deps...)
	}
	return i
}
