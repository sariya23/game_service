package gamestatusrepo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sariya23/game_service/internal/model"
)

func (gs *GameStatusRepository) GetGameStatusByName(ctx context.Context, statusName string) (*model.GameStatus, error) {
	const operationPlace = "postgresql.GetGameStatusByName"
	log := gs.log.With("operationPlace", operationPlace)

	getStatusQuery := fmt.Sprintf("select %s, %s from %s where name=$1",
		GameStatusGameStatusIDFieldName,
		GameStatusGameNameFieldName,
		GameStatusTable,
	)
	var gameStatus model.GameStatus
	statusRow := gs.conn.GetPool().QueryRow(ctx, getStatusQuery, statusName)
	err := statusRow.Scan(
		&gameStatus.ID,
		&gameStatus.Name,
	)
	if err != nil {
		log.Error("cannot get game status by name", slog.String("gameStatusName", statusName))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return &gameStatus, nil
}
