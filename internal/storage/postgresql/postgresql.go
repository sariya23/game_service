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

func (postgresql PostgreSQL) GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (*gamev4.DomainGame, error) {
	panic("impl me")
}

func (postgresql PostgreSQL) GetGameByID(ctx context.Context, gameID uint64) (*gamev4.DomainGame, error) {
	panic("impl me")
}

func (postgresql PostgreSQL) SaveGame(ctx context.Context, game *gamev4.DomainGame) (*gamev4.DomainGame, error) {
	panic("impl me")
}
