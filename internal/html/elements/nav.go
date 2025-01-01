package elements

type Nav struct {
	Current string
	Links   []Link
}

// WithCurrent set `Current` in a new `Nav` and returns it.
// Some may call this an immutable builder method.
func (n Nav) WithCurent(current string) Nav {
	return Nav{
		Current: current,
		Links:   n.Links,
	}
}
