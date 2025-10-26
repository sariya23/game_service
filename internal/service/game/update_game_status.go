package gameservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	game_api "github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/lib/validators"
	"github.com/sariya23/game_service/internal/outerror"
)

// UpdateGameStatus меняет статус видимости игры.
func (gameService *GameService) UpdateGameStatus(ctx context.Context, gameID int64, newStatus game_api.GameStatusType) error {
	const operationPlace = "gameservice.UpdateGameStatus"
	log := gameService.log.With("operationPlace", operationPlace)

	if _, ok := game_api.GameStatusType_name[int32(newStatus)]; !ok {
		log.Warn("pass unknown status", slog.Int64("gameID", gameID), slog.Any("newStatus", newStatus))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrUnknownGameStatus)
	}

	game, err := gameService.gameRepository.GetGameByID(ctx, gameID)
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			log.Warn("game not found to update status", slog.Int64("gameID", gameID))
			return fmt.Errorf("%s: %w", operationPlace, err)
		}
		log.Error("cannot get game by id to set new status", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", operationPlace, err)
	}
	if valid := validators.GameStatus(game.GameStatus, newStatus); !valid {
		log.Warn("cannot update status", slog.Int64("gameID", gameID), slog.Any("currentStatus", game.GameStatus), slog.Any("newStatus", newStatus))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrInvalidNewGameStatus)
	}

	err = gameService.gameRepository.UpdateGameStatus(ctx, gameID, newStatus)
	if err != nil {
		log.Error("cannot update game status", slog.Any("newStatus", newStatus), slog.Int64("gameID", gameID))
		return fmt.Errorf("%s: %w", operationPlace, err)
	}

	return nil
}
