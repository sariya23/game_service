package gameservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
)

func (gameService *GameService) GetGame(
	ctx context.Context,
	gameID uint64,
) (*model.Game, error) {
	const operationPlace = "gameservice.GetGame"
	log := gameService.log.With("operationPlace", operationPlace)
	game, err := gameService.gameRepository.GetGameByID(ctx, gameID)
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			log.Warn(fmt.Sprintf("game with id=%d not found", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameNotFound)
		}
		log.Error(fmt.Sprintf("unexpected error; err=%v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return game, nil
}
