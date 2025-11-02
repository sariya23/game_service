package gamerepo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	gamegenrerepo "github.com/sariya23/game_service/internal/storage/postgresql/game_genre_repo"
	gametagrepo "github.com/sariya23/game_service/internal/storage/postgresql/game_tag_repo"
	"github.com/sariya23/game_service/internal/storage/postgresql/genrerepo"
	"github.com/sariya23/game_service/internal/storage/postgresql/tagrepo"
)

func (gr *GameRepository) GetGameByID(ctx context.Context, gameID int64) (*model.Game, error) {
	const operationPlace = "postgresql.gamerepo.GetGameByID"
	log := gr.log.With("operationPlace", operationPlace)
	getGameMainInfoQuery := fmt.Sprintf(
		"select %s, %s, %s, %s, %s, %s from game where %s=$1",
		GameGameIDFieldName,
		GameTitleFieldName,
		GameDescriptionFieldName,
		GameReleaseDateFieldName,
		GameImageURLFieldName,
		GameGameStatusIDFieldName,
		GameGameIDFieldName,
	)
	getGameGenresQuery := fmt.Sprintf(`
	select %s, %s
	from game join game_genre using(%s) join genre using(%s)
	where %s=$1`,
		genrerepo.GenreGenreNameFieldName,
		genrerepo.GenreGenreIDFieldName,
		gamegenrerepo.GameGenreGameIDFieldName,
		genrerepo.GenreGenreIDFieldName,
		GameGameIDFieldName,
	)
	getGameTagsQuery := fmt.Sprintf(`
	select %s, %s 
	from game join game_tag using(%s) join tag using(%s)
	where %s=$1`,
		tagrepo.TagTagNameFieldName,
		tagrepo.TagTagIDFieldName,
		gametagrepo.GameTagGameIDFieldName,
		tagrepo.TagTagIDFieldName,
		GameGameIDFieldName,
	)
	var gameDB dto.GameDB
	gameRow := gr.conn.GetPool().QueryRow(ctx, getGameMainInfoQuery, gameID)
	err := gameRow.Scan(
		&gameDB.GameID,
		&gameDB.Title,
		&gameDB.Description,
		&gameDB.ReleaseDate,
		&gameDB.ImageURL,
		&gameDB.GameStatus,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("game does not exists", slog.Int64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameNotFound)
		} else {
			log.Error(fmt.Sprintf("Uncaught error: %v", err))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}
	game := gameDB.ToDomain()
	genreRows, err := gr.conn.GetPool().Query(ctx, getGameGenresQuery, gameID)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get genres, uncaught error: %v", err), slog.Int64("gameID", gameID))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer genreRows.Close()
	for genreRows.Next() {
		var gameGenre model.Genre
		err = genreRows.Scan(&gameGenre.GenreName, &gameGenre.GenreID)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game genres, uncaught error: %v", err), slog.Int64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		if genreRows.Err() != nil {
			log.Error("cannot prepare next row", slog.String("err", err.Error()), slog.Int64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		game.Genres = append(game.Genres, gameGenre)
	}

	tagRows, err := gr.conn.GetPool().Query(ctx, getGameTagsQuery, gameID)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get tags, uncaught error: %v", err), slog.Int64("gameID", gameID))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tagRows.Close()
	for tagRows.Next() {
		var gameTag model.Tag
		err = tagRows.Scan(&gameTag.TagName, &gameTag.TagID)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game tags, uncaught error: %v", err), slog.Int64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		if tagRows.Err() != nil {
			log.Error("cannot prepare next row", slog.String("err", err.Error()), slog.Int64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		game.Tags = append(game.Tags, gameTag)
	}

	return &game, nil
}
