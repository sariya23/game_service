package gamerepo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/game_service/internal/lib/logger"
	"github.com/sariya23/game_service/internal/model/dto"
)

func (gr *GameRepository) SaveGame(ctx context.Context, game dto.AddGameService) (int64, error) {
	const operationPlace = "postgresql.gamerepo.SaveGame"
	log := gr.log.With("operationPlace", operationPlace)
	log = logger.EnrichRequestID(ctx, log)
	saveGameArgs := pgx.NamedArgs{
		"title":        game.Title,
		"description":  game.Description,
		"release_date": game.ReleaseDate,
		"image_key":    game.ImageKey,
	}
	saveMainGameInfoQuery := fmt.Sprintf(`
		insert into game (%s, %s, %s, %s) values 
		(@title, @description, @release_date, @image_key)
		returning %s
	`,
		GameTitleFieldName,
		GameDescriptionFieldName,
		GameReleaseDateFieldName,
		GameImageKeyFieldName,
		GameGameIDFieldName)
	addTagsForGameQuery := "insert into game_tag values ($1, $2)"
	addGenresForGameQuery := "insert into game_genre values ($1, $2)"

	tx, err := gr.conn.GetPool().Begin(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("cannot start transaction, unexpected error = %v", err))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Error("cannot rollback transaction, unexpected error", slog.String("error", err.Error()))
		}
	}()

	var savedGameID int64
	saveGameRow := tx.QueryRow(ctx, saveMainGameInfoQuery, saveGameArgs)
	err = saveGameRow.Scan(&savedGameID)
	if err != nil {
		log.Error(fmt.Sprintf("cannot save game, unexpected error = %v", err))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}

	for _, tagID := range game.TagIDs {
		_, err = tx.Exec(ctx, addTagsForGameQuery, savedGameID, tagID)
		if err != nil {
			log.Error(fmt.Sprintf("cannot link tag with game, unexpected error = %v", err), slog.Int64("tagID", tagID), slog.Int("gameID", int(savedGameID)))
			return 0, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}

	for _, genreID := range game.GenreIDs {
		_, err = tx.Exec(ctx, addGenresForGameQuery, savedGameID, genreID)
		if err != nil {
			log.Error(fmt.Sprintf("cannot link tag with game, unexpected error = %v", err), slog.Int64("genreID", genreID), slog.Int("gameID", int(savedGameID)))
			return 0, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("cannot commit, err = %v", err))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return savedGameID, nil
}
