package helpers

func HasIntersection[T comparable](a, b []T) bool {
	set := make(map[T]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}
	for _, v := range b {
		if _, ok := set[v]; ok {
			return true
		}
	}
	return false
}
