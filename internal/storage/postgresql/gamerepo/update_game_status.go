package gamerepo

import (
	"context"
	"fmt"
	"log/slog"

	game_api "github.com/sariya23/api_game_service/gen/game"
)

func (gr *GameRepository) UpdateGameStatus(ctx context.Context, gameID int64, newStatus game_api.GameStatusType) error {
	const operationPlace = "postgresql.UpdateGameStatus"
	log := gr.log.With("operationPlace", operationPlace)
	queryUpdateStatusQuery := fmt.Sprintf("update game set %s=$1 where %s=$2", GameGameStatusIDFieldName, GameGameIDFieldName)
	_, err := gr.conn.GetPool().Exec(ctx, queryUpdateStatusQuery, newStatus, gameID)
	if err != nil {
		log.Error("cannot update game status", slog.Int64("gameID", gameID), slog.Any("newStatus", newStatus))
		return fmt.Errorf("%s: %w", operationPlace, err)
	}
	return nil
}
