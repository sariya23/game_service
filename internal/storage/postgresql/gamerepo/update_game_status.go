package gamerepo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

func (gr *GameRepository) UpdateGameStatus(ctx context.Context, gameID int64, newStatus gamev2.GameStatusType) error {
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
