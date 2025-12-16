package validators

import (
	"github.com/sariya23/api_game_service/gen/game"
)

// GameStatus ...
// Запрещено:
// DRAFT -> PUBLISH
// PUBLISH -> PENDING
// PUBLISH -> DRAFT
func GameStatus(currentStatus, newStatus game.GameStatusType) bool {
	if currentStatus == game.GameStatusType_DRAFT && newStatus == game.GameStatusType_PUBLISH {
		return false
	} else if currentStatus == game.GameStatusType_PUBLISH && newStatus == game.GameStatusType_PENDING {
		return false
	} else if currentStatus == game.GameStatusType_PUBLISH && newStatus == game.GameStatusType_DRAFT {
		return false
	}
	return true
}
