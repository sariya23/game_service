//go:build integrations

package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCasePickMostFreq[T comparable] struct {
	name     string
	values   []T
	expected T
}

func TestPickMostFrequentValue(t *testing.T) {
	cases := []testCasePickMostFreq[string]{
		{
			name:     "success",
			values:   []string{"A", "A", "C", "Q"},
			expected: "A",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := PickMostFrequentValue(tc.values)
			assert.Equal(t, tc.expected, res)
		})
	}
}
