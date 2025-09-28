package postgresql

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sariya23/game_service/internal/model"
)

func (postgresql PostgreSQL) GetGameStatusByName(ctx context.Context, statusName string) (*model.GameStatus, error) {
	const operationPlace = "postgresql.GetGameStatusByName"
	log := postgresql.log.With("operationPlace", operationPlace)

	getStatusQuery := fmt.Sprintf("select %s, %s from %s where name=$1",
		gameStatusGameStatusIDFieldName,
		gameStatusGameNameFieldName,
		gameStatusTable,
	)
	var gameStatus model.GameStatus
	statusRow := postgresql.connection.QueryRow(ctx, getStatusQuery, statusName)
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
