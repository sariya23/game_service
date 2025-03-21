package postgresql

import (
	"context"
	"log/slog"

	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type PostgreSQL struct {
	log *slog.Logger
}

func MustNewConnection(log *slog.Logger) PostgreSQL {
	return PostgreSQL{
		log: log,
	}
}

func (postgresql PostgreSQL) GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (gamev4.GameWithRating, error) {
	panic("impl me")
}
