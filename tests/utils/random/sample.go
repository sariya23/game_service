//go:build integrations

package random

import (
	"math/rand"
)

func Sample[T any](items []T, k int) []T {
	if k > len(items) {
		k = len(items)
	}

	cpy := append([]T(nil), items...)
	rand.Shuffle(len(cpy), func(i, j int) { cpy[i], cpy[j] = cpy[j], cpy[i] })

	return cpy[:k]
}
