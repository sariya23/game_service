package pointer

func OfSlice[T any](sl []T) []*T {
	res := make([]*T, len(sl))
	for i, v := range sl {
		res[i] = &v
	}
	return res
}
