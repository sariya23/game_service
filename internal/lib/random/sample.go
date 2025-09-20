package random

import "math/rand/v2"

func Sample[T any](slice []T, k int) []T {
	if k > len(slice) {
		k = len(slice)
	}
	// копируем, чтобы не портить исходный
	cp := append([]T(nil), slice...)
	rand.Shuffle(len(cp), func(i, j int) {
		cp[i], cp[j] = cp[j], cp[i]
	})
	return cp[:k]
}
