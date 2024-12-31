package utils

func Ternary[T comparable](cond bool, then T, or T) T {
	if cond {
		return then
	}
	return or
}
