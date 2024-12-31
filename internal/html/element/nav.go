package element

type Nav struct {
	Current string
	Links   []Link
}

func (n Nav) WithCurent(current string) Nav {
	return Nav{
		Current: current,
		Links:   n.Links,
	}
}
