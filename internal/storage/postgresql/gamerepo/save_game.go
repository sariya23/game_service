package gamerepo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/game_service/internal/model"
)

func (r GameRepository) SaveGame(ctx context.Context, game model.Game) (int64, error) {
	const operationPlace = "postgresql.gamerepo.SaveGame"
	log := r.log.With("operationPlace", operationPlace)
	saveGameArgs := pgx.NamedArgs{
		"title":        game.Title,
		"description":  game.Description,
		"release_date": game.ReleaseDate,
		"image_url":    game.ImageURL,
	}
	saveMainGameInfoQuery := fmt.Sprintf(`
		insert into game (%s, %s, %s, %s) values 
		(@title, @description, @release_date, @image_url)
		returning %s
	`,
		GameTitleFieldName,
		GameDescriptionFieldName,
		GameReleaseDateFieldName,
		GameImageURLFieldName,
		GameGameIDFieldName)
	addTagsForGameQuery := "insert into game_tag values ($1, $2)"
	addGenresForGameQuery := "insert into game_genre values ($1, $2)"
	genreIDs := make([]int, 0, len(game.Genres))
	tagIDs := make([]int, 0, len(game.Tags))

	if len(game.Genres) != 0 {
		for _, g := range game.Genres {
			genreIDs = append(genreIDs, int(g.GenreID))
		}
	}

	if len(game.Tags) != 0 {
		for _, t := range game.Tags {
			tagIDs = append(tagIDs, int(t.TagID))
		}
	}

	tx, err := r.conn.StartTransaction(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("cannot start transaction, unexpected error = %v", err))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tx.Rollback(ctx)

	var savedGameID int64
	saveGameRow := tx.QueryRow(ctx, saveMainGameInfoQuery, saveGameArgs)
	err = saveGameRow.Scan(&savedGameID)
	if err != nil {
		log.Error(fmt.Sprintf("cannot save game, unexpected error = %v", err))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}

	for _, tagID := range tagIDs {
		_, err = tx.Exec(ctx, addTagsForGameQuery, savedGameID, tagID)
		if err != nil {
			log.Error(fmt.Sprintf("cannot link tag with game, unexpected error = %v", err), slog.Int("tagID", tagID), slog.Int("gameID", int(savedGameID)))
			return 0, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}

	for _, genreID := range genreIDs {
		_, err = tx.Exec(ctx, addGenresForGameQuery, savedGameID, genreID)
		if err != nil {
			log.Error(fmt.Sprintf("cannot link tag with game, unexpected error = %v", err), slog.Int("genreID", genreID), slog.Int("gameID", int(savedGameID)))
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
