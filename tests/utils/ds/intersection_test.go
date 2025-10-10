//go:build integrations

package ds

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseIntersection[T comparable] struct {
	name       string
	arr1, arr2 []T
	expected   []T
}

func TestIntersection(t *testing.T) {
	cases := []testCaseIntersection[string]{
		{
			name:     "has intersection",
			arr1:     []string{"A", "B", "C"},
			arr2:     []string{"A", "C", "Q"},
			expected: []string{"A", "C"},
		},
		{
			name:     "no intersection",
			arr1:     []string{"A", "B", "C"},
			arr2:     []string{"Z", "X", "Q"},
			expected: []string{},
		},
		{
			name:     "one of empty",
			arr1:     []string{},
			arr2:     []string{"Z", "X", "Q"},
			expected: []string{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := Intersection(tc.arr1, tc.arr2)
			sort.Slice(res, func(i, j int) bool {
				return res[i] < res[j]
			})
			assert.Equal(t, tc.expected, res)
		})
	}
}
