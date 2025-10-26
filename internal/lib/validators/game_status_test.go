package validators

import (
	"testing"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/stretchr/testify/assert"
)

func TestGameStatus(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name                     string
		currentStatus, newStatus game.GameStatusType
		expected                 bool
	}{
		{
			name:          "Invalid, from DRAFT to PUBLISH",
			currentStatus: game.GameStatusType_DRAFT,
			newStatus:     game.GameStatusType_PUBLISH,
			expected:      false,
		},
		{
			name:          "Invalid, from PUBLISH to PENDING",
			currentStatus: game.GameStatusType_PUBLISH,
			newStatus:     game.GameStatusType_PENDING,
			expected:      false,
		},
		{
			name:          "Invalid, from PUBLISH to DRAFT",
			currentStatus: game.GameStatusType_PUBLISH,
			newStatus:     game.GameStatusType_DRAFT,
			expected:      false,
		},
		{
			name:          "Valid, from DRAFT to PENDING",
			currentStatus: game.GameStatusType_DRAFT,
			newStatus:     game.GameStatusType_PENDING,
			expected:      true,
		},
		{
			name:          "Valid, from PENDING to PUBLISH",
			currentStatus: game.GameStatusType_PENDING,
			newStatus:     game.GameStatusType_PUBLISH,
			expected:      true,
		},
		{
			name:          "Valid, from PENDING to DRAFT",
			currentStatus: game.GameStatusType_PENDING,
			newStatus:     game.GameStatusType_DRAFT,
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
