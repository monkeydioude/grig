package dt

func Any[T comparable](slice []T, elm T) bool {
	v, _ := AnyFunc(slice, func(el T) bool {
		return el == elm
	})
	return v
}

func AnyFunc[T comparable](slice []T, fn func(T) bool) (bool, T) {
	for _, item := range slice {
		if fn(item) {
			return true, item
		}
	}
	var zero T
	return false, zero
}

func AppendUnique[T comparable](slice []T, elem T) ([]T, bool) {
	if Any(slice, elem) {
		return slice, false
	}
	return append(slice, elem), true
}
