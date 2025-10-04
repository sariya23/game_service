package gameservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/game_service/internal/outerror"
	gamev2 "github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

// UpdateGameStatus меняет статус видимости игры.
// Нельзя:
// DRAFT -> PUBLISH
// PUBLISH -> PENDING
// PUBLISH -> DRAFT
func (gameService *GameService) UpdateGameStatus(ctx context.Context, gameID int64, newStatus gamev2.GameStatusType) error {
	const operationPlace = "gameservice.UpdateGameStatus"
	log := gameService.log.With("operationPlace", operationPlace)

	if _, ok := gamev2.GameStatusType_name[int32(newStatus)]; !ok {
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

	if game.GameStatus == int(gamev2.GameStatusType_DRAFT) && newStatus == gamev2.GameStatusType_PUBLISH {
		log.Warn("cannot update status from DRAFT to PUBLISH", slog.Int64("gameID", gameID))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrInvalidNewGameStatus)
	} else if game.GameStatus == int(gamev2.GameStatusType_PUBLISH) && newStatus == gamev2.GameStatusType_PENDING {
		log.Warn("cannot update status from PUBLISH to PENDING", slog.Int64("gameID", gameID))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrInvalidNewGameStatus)
	} else if game.GameStatus == int(gamev2.GameStatusType_PUBLISH) && newStatus == gamev2.GameStatusType_DRAFT {
		log.Warn("cannot update status from PUBLISH to DRAFT", slog.Int64("gameID", gameID))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrInvalidNewGameStatus)
	}

	err = gameService.gameRepository.UpdateGameStatus(ctx, gameID, newStatus)
	if err != nil {
		log.Error("cannot update game status", slog.Any("newStatus", newStatus), slog.Int64("gameID", gameID))
		return fmt.Errorf("%s: %w", operationPlace, err)
	}

	return nil
}
