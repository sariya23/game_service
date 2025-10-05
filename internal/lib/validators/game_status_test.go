package validators

import (
	"testing"

	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"github.com/stretchr/testify/assert"
)

func TestGameStatus(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name                     string
		currentStatus, newStatus gamev2.GameStatusType
		expected                 bool
	}{
		{
			name:          "Invalid, from DRAFT to PUBLISH",
			currentStatus: gamev2.GameStatusType_DRAFT,
			newStatus:     gamev2.GameStatusType_PUBLISH,
			expected:      false,
		},
		{
			name:          "Invalid, from PUBLISH to PENDING",
			currentStatus: gamev2.GameStatusType_PUBLISH,
			newStatus:     gamev2.GameStatusType_PENDING,
			expected:      false,
		},
		{
			name:          "Invalid, from PUBLISH to DRAFT",
			currentStatus: gamev2.GameStatusType_PUBLISH,
			newStatus:     gamev2.GameStatusType_DRAFT,
			expected:      false,
		},
		{
			name:          "Valid, from DRAFT to PENDING",
			currentStatus: gamev2.GameStatusType_DRAFT,
			newStatus:     gamev2.GameStatusType_PENDING,
			expected:      true,
		},
		{
			name:          "Valid, from PENDING to PUBLISH",
			currentStatus: gamev2.GameStatusType_PENDING,
			newStatus:     gamev2.GameStatusType_PUBLISH,
			expected:      true,
		},
		{
			name:          "Valid, from PENDING to DRAFT",
			currentStatus: gamev2.GameStatusType_PENDING,
			newStatus:     gamev2.GameStatusType_DRAFT,
			expected:      true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, GameStatus(tc.currentStatus, tc.newStatus))
		})
	}
}
