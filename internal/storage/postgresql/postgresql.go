package postgresql

import (
	"context"
	"log/slog"

	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type PostgreSQL struct {
	log *slog.Logger
}

type GameTransaction struct {
}

func (gt GameTransaction) Approve() {
	panic("impl me")
}

func (gt GameTransaction) Reject() {
	panic("impl me")
}

func MustNewConnection(log *slog.Logger) PostgreSQL {
	return PostgreSQL{
		log: log,
	}
}

func (postgresql PostgreSQL) GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (gamev4.Game, error) {
	panic("impl me")
}

func (postgresql PostgreSQL) SaveGame(ctx context.Context, game *gamev4.Game) (*GameTransaction, error) {
	panic("impl me")
}
