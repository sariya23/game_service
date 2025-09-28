package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// Game table field names
	gameGameIDFieldName       = "game_id"
	gameTitleFieldName        = "title"
	gameDescriptionFieldName  = "description"
	gameReleaseDateFieldName  = "release_date"
	gameImageURLFieldName     = "image_url"
	gameGameStatusIDFieldName = "game_status_id"

	// GameGenre table field names
	gameGenreGameIDFieldName  = "game_id"
	gameGenreGenreIDFieldName = "genre_id"

	// GameTag table field names
	gameTagGameIDFieldName = "game_id"
	gameTagTagIDFieldName  = "tag_id"

	// Genre table field names
	genreGenreIDFieldName   = "genre_id"
	genreGenreNameFieldName = "genre_name"

	// Tag table field names
	tagTagIDFieldName   = "tag_id"
	tagTagNameFieldName = "tag_name"

	// GameStatus table
	gameStatusTable                 = "game_status"
	gameStatusGameStatusIDFieldName = "game_status_id"
	gameStatusGameNameFieldName     = "name"
)

type PostgreSQL struct {
	log        *slog.Logger
	connection *pgxpool.Pool
}

func MustNewConnection(ctx context.Context, log *slog.Logger, dbURL string) PostgreSQL {
	const opearationPlace = "postgresql.MustNewConnection"
	localLog := log.With("operationPlace", opearationPlace)
	ctx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		localLog.Error(fmt.Sprintf("%s: cannot connect to db with URL: %s, with error: %v", opearationPlace, dbURL, err))
		panic(fmt.Sprintf("%s: cannot connect to db with URL: %s, with error: %v", opearationPlace, dbURL, err))
	}
	err = conn.Ping(ctx)
	if err != nil {
		localLog.Error(fmt.Sprintf("%s: db is unreachable: %v", opearationPlace, err))
		panic(fmt.Sprintf("%s: db is unreachable: %v", opearationPlace, err))
	}
	localLog.Info("Postgres ready to get connections")
	return PostgreSQL{log: log, connection: conn}
}
