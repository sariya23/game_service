//go:build integrations

package random

func PickMostFrequentValue[T comparable](values []T) T {
	m := make(map[T]int)

	for _, v := range values {
		m[v]++
	}
	max := struct {
		cnt int
		val T
	}{}
	for k, v := range m {
		if v > max.cnt {
			max.cnt = v
			max.val = k
		}
	}
	return max.val
}
