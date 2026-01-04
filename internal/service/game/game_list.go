package gameservice

import (
	"context"
	"fmt"

	"github.com/sariya23/game_service/internal/lib/logger"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
)

func (gameService *GameService) GameList(
	ctx context.Context,
	gameFilters dto.GameFilters,
	limit uint32,
) ([]model.ShortGame, error) {
	const operationPlace = "gameservice.GetTopGames"
	log := gameService.log.With("operationPlace", operationPlace)
	log = logger.EnrichRequestID(ctx, log)
	if limit == 0 {
		limit = 10
	}
	gamesNoImageURL, err := gameService.gameRepository.GameList(ctx, gameFilters, limit)
	if err != nil {
		log.Error(fmt.Sprintf("unexpected error; err=%v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}

	games := make([]model.ShortGame, 0, len(gamesNoImageURL))
	for _, g := range gamesNoImageURL {
		imageURL, err := gameService.s3Storager.GeneratePresignedURL(ctx, g.ImageKey)
		if err != nil {
			log.Warn(fmt.Sprintf("unexpected error while generate URL; err=%v", err))
		}
		shortGame := g.ToShortGame(imageURL)
		games = append(games, shortGame)
	}

	return games, nil
}
