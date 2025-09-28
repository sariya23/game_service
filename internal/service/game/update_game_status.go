package gameservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

func (gameService *GameService) UpdateGameStatus(ctx context.Context, gameID uint64, newStatus gamev4.GameStatusType) error {
	const operationPlace = "gameservice.UpdateGameStatus"
	log := gameService.log.With("operationPlace", operationPlace)
	game, err := gameService.gameRepository.GetGameByID(ctx, gameID)
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			log.Warn("game not found to update status", slog.Uint64("gameID", gameID))
			return fmt.Errorf("%s: %w", operationPlace, err)
		}
		log.Error("cannot get game by id to set new status", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", operationPlace, err)
	}

	if int(gamev4.GameStatusType_DRAFT) == game.GameStatus && newStatus == gamev4.GameStatusType_PUBLISH {
		log.Warn("cannot update status from DRAFT to PUBLISH", slog.Uint64("gameID", gameID))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrInvalidNewGameStatus)
	}
	return nil
}
