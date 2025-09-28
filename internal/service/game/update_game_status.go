package gameservice

import (
	"context"

	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

func (gameService *GameService) UpdateGameStatus(ctx context.Context, gameID uint64, newStatus gamev4.GameStatusType) error {
	const operationPlace = "gameservice.UpdateGameStatus"
	log := gameService.log.With("operationPlace", operationPlace)
}
