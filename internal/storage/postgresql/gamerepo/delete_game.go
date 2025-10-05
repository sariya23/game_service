package gamerepo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
)

func (gr *GameRepository) DaleteGame(ctx context.Context, gameID int64) (*dto.DeletedGame, error) {
	const operationPlace = "postgresql.DeleteGame"
	log := gr.log.With("operationPlace", operationPlace)
	deleteGameQuery := fmt.Sprintf(
		"delete from game where %s=$1 returning %s, extract(year from %s), %s",
		GameGameIDFieldName,
		GameGameIDFieldName,
		GameReleaseDateFieldName,
		GameTitleFieldName,
	)
	var deltedGameInfo dto.DeletedGame
	deleteGameRow := gr.conn.GetPool().QueryRow(ctx, deleteGameQuery, gameID)
	err := deleteGameRow.Scan(&deltedGameInfo.GameID, &deltedGameInfo.ReleaseYear, &deltedGameInfo.Title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("cannot delete game because it is not found", slog.Int("gameID", int(gameID)))
			return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameNotFound)
		}
		log.Error("cannot delete game", slog.Any("gameID", gameID), slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	log.Info("game deleted successfully", slog.Int("gameID", int(gameID)))
	return &deltedGameInfo, nil
}
