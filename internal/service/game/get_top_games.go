package gameservice

import (
	"context"
	"fmt"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
)

func (gameService *GameService) GetTopGames(
	ctx context.Context,
	gameFilters dto.GameFilters,
	limit uint32,
) ([]model.ShortGame, error) {
	const operationPlace = "gameservice.GetTopGames"
	log := gameService.log.With("operationPlace", operationPlace)
	if limit == 0 {
		limit = 10
	}
	games, err := gameService.gameRepository.GetTopGames(ctx, gameFilters, limit)
	if err != nil {
		log.Error(fmt.Sprintf("unexcpected error; err=%v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return games, nil
}
