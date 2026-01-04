package gameservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/game_service/internal/lib/logger"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
)

func (gameService *GameService) GetGame(
	ctx context.Context,
	gameID int64,
) (*model.Game, error) {
	const operationPlace = "gameservice.GetGame"
	log := gameService.log.With("operationPlace", operationPlace)
	log = logger.EnrichRequestID(ctx, log)
	gameNoImageURL, err := gameService.gameRepository.GetGameByID(ctx, gameID)
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			log.Warn("game not found", slog.Int64("game_id", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameNotFound)
		}
		log.Error("unexpected error from repository", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	imageURL, err := gameService.s3Storager.GeneratePresignedURL(ctx, gameNoImageURL.ImageKey)
	if err != nil {
		log.Warn("failed to generate presigned URL for image", slog.String("error", err.Error()))
	}
	game := gameNoImageURL.ToDomain(imageURL)
	return &game, nil
}
