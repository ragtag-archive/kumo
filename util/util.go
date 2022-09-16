package util

// Set[T] is a set of items of type T
type Set[T comparable] map[T]struct{}

// SetDifference returns the difference between two sets.
func SetDifference[T comparable](a, b Set[T]) Set[T] {
	c := make(Set[T])
	for k := range a {
		if _, ok := b[k]; !ok {
			c[k] = struct{}{}
		}
	}
	return c
}

// SetUnion returns the union between two sets.
func SetUnion[T comparable](a, b Set[T]) Set[T] {
	c := make(Set[T])
	for k := range a {
		c[k] = struct{}{}
	}
	for k := range b {
		c[k] = struct{}{}
	}
	return c
}

// SetToSlice converts a set to a slice.
func SetToSlice[T comparable](s Set[T]) []T {
	res := make([]T, 0, len(s))
	for k := range s {
		res = append(res, k)
	}
	return res
}
