package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type PostgreSQL struct {
	log        *slog.Logger
	connection *pgxpool.Pool
}

func MustNewConnection(ctx context.Context, log *slog.Logger, dbURL string) PostgreSQL {
	const opearationPlace = "postgresql.MustNewConnection"
	log = log.With("operationPlace", opearationPlace)
	ctx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Error(fmt.Sprintf("%s: cannot connect to db with URL: %s, with error: %v", opearationPlace, dbURL, err))
		panic(fmt.Sprintf("%s: cannot connect to db with URL: %s, with error: %v", opearationPlace, dbURL, err))
	}
	err = conn.Ping(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("%s: db is unreachable: %v", opearationPlace, err))
		panic(fmt.Sprintf("%s: db is unreachable: %v", opearationPlace, err))
	}
	log.Info("Postgres ready to get connections")
	return PostgreSQL{log: log, connection: conn}
}

func (postgresql PostgreSQL) GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (*gamev4.DomainGame, error) {
	panic("impl me")

}

func (postgresql PostgreSQL) GetGameByID(ctx context.Context, gameID uint64) (*gamev4.DomainGame, error) {
	const operationPlace = "postgresql.GetGameByID"
	log := postgresql.log.With("operationPlace", operationPlace)
	getGameMainInfoQuery := fmt.Sprintf(
		"select %s, %s, %s, %s, %s from game where %s=$1",
		gameGameIDFieldName,
		gameTitleFieldName,
		gameDescriptionFieldName,
		gameReleaseDateFieldName,
		gameImageURLFieldName,
		gameGameIDFieldName,
	)
	getGameGenresQuery := fmt.Sprintf(`
	select %s 
	from game join game_genre using(%s) join genre using(%s)
	where %s=$1`,
		genreGenreNameFieldName,
		gameGenreGameIDFieldName,
		genreGenreIDFieldName,
		gameGameIDFieldName,
	)
	getGameTagsQuery := fmt.Sprintf(`
	select %s 
	from game join game_tag using(%s) join tag using(%s)
	where %s=$1`,
		tagTagNameFieldName,
		gameTagGameIDFieldName,
		tagTagIDFieldName,
		gameGameIDFieldName,
	)
	var game gamev4.DomainGame
	gameRow := postgresql.connection.QueryRow(ctx, getGameMainInfoQuery, gameID)
	err := gameRow.Scan(
		&game.ID,
		&game.Title,
		&game.Description,
		&game.ReleaseDate,
		&game.CoverImageUrl,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("cannot get game", slog.Uint64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameNotFound)
		} else {
			log.Error(fmt.Sprintf("Uncaught error: %v", err))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}
	genreRows, err := postgresql.connection.Query(ctx, getGameGenresQuery, gameID)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get genres, uncaught error: %v", err), slog.Uint64("gameID", gameID))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer genreRows.Close()
	for genreRows.Next() {
		var gameGenreName string
		err = genreRows.Scan(&gameGenreName)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game genres, uncaught error: %v", err), slog.Uint64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		game.Genres = append(game.Genres, gameGenreName)
	}
	if genreRows.Err() != nil {
		log.Error(fmt.Sprintf("Uncaught error: %v", err), slog.Uint64("gameID", gameID))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}

	tagRows, err := postgresql.connection.Query(ctx, getGameTagsQuery, gameID)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get tags, uncaught error: %v", err), slog.Uint64("gameID", gameID))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tagRows.Close()
	for tagRows.Next() {
		var gameTagName string
		err = tagRows.Scan(&gameTagName)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game tags, uncaught error: %v", err), slog.Uint64("gameID", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		game.Tags = append(game.Tags, gameTagName)
	}
	if tagRows.Err() != nil {
		log.Error(fmt.Sprintf("Uncaught error: %v", err), slog.Uint64("gameID", gameID))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}

	return &game, nil
}

func (postgresql PostgreSQL) SaveGame(ctx context.Context, game *gamev4.DomainGame) (*gamev4.DomainGame, error) {
	const operationPlace = "postgresql.SaveGame"
	log := postgresql.log.With("operationPlace", operationPlace)
	saveGameArgs := pgx.NamedArgs{
		"title":        game.GetTitle(),
		"description":  game.GetDescription(),
		"release_date": game.GetReleaseDate(),
		"image_url":    game.GetCoverImageUrl(),
	}
	saveMainGameInfoQuery := fmt.Sprintf(`
		insert into game (%s, %s, %s, %s) values 
		(@title, @description, @release_date, @image_url)
		returning %s
	`, gameTitleFieldName, gameDescriptionFieldName, gameReleaseDateFieldName, gameImageURLFieldName, gameGameIDFieldName)

	getTagIdQuery := fmt.Sprintf("select genre_id from genre where %s=$1", genreGenreNameFieldName)
	getGenreIdQuery := fmt.Sprintf("select tag_id from tag where %s=$1", tagTagNameFieldName)
	addTagsForGameQuery := "insert into game_tag values ($1, $2)"
	addGenresForGameQuery := "insert into game_genre values ($1, $2)"
	genreIDs := make([]int, 0, len(game.GetGenres()))
	tagIDs := make([]int, 0, len(game.GetTags()))

	if len(game.GetGenres()) != 0 {
		for _, genreName := range game.GetGenres() {
			var genreID int
			genreRow := postgresql.connection.QueryRow(ctx, getGenreIdQuery, genreName)
			err := genreRow.Scan(&genreID)
			if err != nil {
				log.Error(fmt.Sprintf("cannot get genre by id, unexpected error = %v", err), slog.String("request Genge", genreName))
				return nil, fmt.Errorf("%s: %w", operationPlace, err)
			}
			genreIDs = append(genreIDs, genreID)
		}
	}

	if len(game.GetTags()) != 0 {
		for _, tagName := range game.GetTags() {
			var tagID int
			tagRow := postgresql.connection.QueryRow(ctx, getTagIdQuery, tagName)
			err := tagRow.Scan(&tagID)
			if err != nil {
				log.Error(fmt.Sprintf("cannot get tag by id, unexpected error = %v", err), slog.String("request Tag", tagName))
				return nil, fmt.Errorf("%s: %w", operationPlace, err)
			}
			tagIDs = append(tagIDs, tagID)
		}
	}

	tx, err := postgresql.connection.Begin(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("cannot start transaction, unexpected error = %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer tx.Rollback(ctx)

	var savedGameID uint64
	saveGameRow := tx.QueryRow(ctx, saveMainGameInfoQuery, saveGameArgs)
	err = saveGameRow.Scan(&savedGameID)
	if err != nil {
		log.Error(fmt.Sprintf("cannot save game, unexpected error = %v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}

	for _, tagID := range tagIDs {
		_, err = tx.Exec(ctx, addTagsForGameQuery, savedGameID, tagID)
		if err != nil {
			log.Error(fmt.Sprintf("cannot link tag with game, unexpected error = %v", err), slog.Int("tagID", tagID), slog.Int("gameID", int(savedGameID)))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}

	for _, genreID := range genreIDs {
		_, err = tx.Exec(ctx, addGenresForGameQuery, savedGameID, genreID)
		if err != nil {
			log.Error(fmt.Sprintf("cannot link tag with game, unexpected error = %v", err), slog.Int("genreID", genreID), slog.Int("gameID", int(savedGameID)))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}
	return game, nil
}

func (postgresql PostgreSQL) GetTopGames(ctx context.Context, releaseYear int32, tags []string, genres []string, limit uint32) (games []*gamev4.DomainGame, err error) {
	panic("impl me")
}

func (postgresql PostgreSQL) DaleteGame(ctx context.Context, gameID uint64) (*gamev4.DomainGame, error) {
	panic("impl me")
}
