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
			return nil, outerror.ErrGameNotFound
		} else {
			log.Error(fmt.Sprintf("Uncaught error: %v", err))
			return nil, err
		}
	}
	genreRows, err := postgresql.connection.Query(ctx, getGameGenresQuery, gameID)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get genres, uncaught error: %v", err), slog.Uint64("gameID", gameID))
		return nil, err
	}
	defer genreRows.Close()
	for genreRows.Next() {
		var gameGenreName string
		err = genreRows.Scan(&gameGenreName)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game genres, uncaught error: %v", err), slog.Uint64("gameID", gameID))
			return nil, err
		}
		game.Genres = append(game.Genres, gameGenreName)
	}
	if genreRows.Err() != nil {
		log.Error(fmt.Sprintf("Uncaught error: %v", err), slog.Uint64("gameID", gameID))
		return nil, err
	}

	tagRows, err := postgresql.connection.Query(ctx, getGameTagsQuery, gameID)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot get tags, uncaught error: %v", err), slog.Uint64("gameID", gameID))
		return nil, err
	}
	defer tagRows.Close()
	for tagRows.Next() {
		var gameTagName string
		err = tagRows.Scan(&gameTagName)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot scan game tags, uncaught error: %v", err), slog.Uint64("gameID", gameID))
			return nil, err
		}
		game.Tags = append(game.Tags, gameTagName)
	}
	if tagRows.Err() != nil {
		log.Error(fmt.Sprintf("Uncaught error: %v", err), slog.Uint64("gameID", gameID))
		return nil, err
	}

	return &game, nil
}

func (postgresql PostgreSQL) SaveGame(ctx context.Context, game *gamev4.DomainGame) (*gamev4.DomainGame, error) {
	panic("impl me")
}

func (postgresql PostgreSQL) GetTopGames(ctx context.Context, releaseYear int32, tags []string, genres []string, limit uint32) (games []*gamev4.DomainGame, err error) {
	panic("impl me")
}

func (postgresql PostgreSQL) DaleteGame(ctx context.Context, gameID uint64) (*gamev4.DomainGame, error) {
	panic("impl me")
}
