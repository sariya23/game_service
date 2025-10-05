package validators

import (
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

// GameStatus. Запрещено:
// DRAFT -> PUBLISH
// PUBLISH -> PENDING
// PUBLISH -> DRAFT
func GameStatus(currentStatus, newStatus gamev2.GameStatusType) bool {
	if currentStatus == gamev2.GameStatusType_DRAFT && newStatus == gamev2.GameStatusType_PUBLISH {
		return false
	} else if currentStatus == gamev2.GameStatusType_PUBLISH && newStatus == gamev2.GameStatusType_PENDING {
		return false
	} else if currentStatus == gamev2.GameStatusType_PUBLISH && newStatus == gamev2.GameStatusType_DRAFT {
		return false
	}
	return true
}
