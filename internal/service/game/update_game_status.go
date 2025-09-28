package gameservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

// UpdateGameStatus меняет статус видимости игры.
// Нельзя:
// DRAFT -> PUBLISH
// PUBLISH -> PENDING
// PUBLISH -> DRAFT
func (gameService *GameService) UpdateGameStatus(ctx context.Context, gameID uint64, newStatus gamev4.GameStatusType) error {
	const operationPlace = "gameservice.UpdateGameStatus"
	log := gameService.log.With("operationPlace", operationPlace)

	if _, ok := gamev4.GameStatusType_name[int32(newStatus)]; !ok {
		log.Warn("pass unknown status", slog.Uint64("gameID", gameID), slog.Any("newStatus", newStatus))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrUnknownGameStatus)
	}

	game, err := gameService.gameRepository.GetGameByID(ctx, gameID)
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			log.Warn("game not found to update status", slog.Uint64("gameID", gameID))
			return fmt.Errorf("%s: %w", operationPlace, err)
		}
		log.Error("cannot get game by id to set new status", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", operationPlace, err)
	}

	if game.GameStatus == int(gamev4.GameStatusType_DRAFT) && newStatus == gamev4.GameStatusType_PUBLISH {
		log.Warn("cannot update status from DRAFT to PUBLISH", slog.Uint64("gameID", gameID))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrInvalidNewGameStatus)
	} else if game.GameStatus == int(gamev4.GameStatusType_PUBLISH) && newStatus == gamev4.GameStatusType_PENDING {
		log.Warn("cannot update status from PUBLISH to PENDING", slog.Uint64("gameID", gameID))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrInvalidNewGameStatus)
	} else if game.GameStatus == int(gamev4.GameStatusType_PUBLISH) && newStatus == gamev4.GameStatusType_DRAFT {
		log.Warn("cannot update status from PUBLISH to DRAFT", slog.Uint64("gameID", gameID))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrInvalidNewGameStatus)
	}

	err = gameService.gameRepository.UpdateGameStatus(ctx, gameID, newStatus)
	if err != nil {
		log.Error("cannot update game status", slog.Any("newStatus", newStatus), slog.Uint64("gameID", gameID))
		return fmt.Errorf("%s: %w", operationPlace, err)
	}

	return nil
}
