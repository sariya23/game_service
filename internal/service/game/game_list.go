package gameservice

import (
	"context"
	"fmt"

	"github.com/sariya23/game_service/internal/interceptors"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
)

func (gameService *GameService) GameList(
	ctx context.Context,
	gameFilters dto.GameFilters,
	limit uint32,
) ([]model.ShortGame, error) {
	const operationPlace = "gameservice.GetTopGames"
	requestID, _ := ctx.Value(interceptors.RequestIDKey).(string)
	log := gameService.log.With("operationPlace", operationPlace)
	log = log.With("request_id", requestID)
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
			continue
		}
		shortGame := g.ToShortGame(imageURL)
		games = append(games, shortGame)
	}

	return games, nil
}
