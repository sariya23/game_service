package postgresql

import (
	"context"
	"log/slog"

	"github.com/sariya23/game_service/internal/model/domain"
)

type PostgreSQL struct {
	log *slog.Logger
}

type GameTransaction struct {
}

func (gt GameTransaction) Approve() error {
	panic("impl me")
}

func (gt GameTransaction) Reject() error {
	panic("impl me")
}

func MustNewConnection(log *slog.Logger) PostgreSQL {
	return PostgreSQL{
		log: log,
	}
}

func (postgresql PostgreSQL) GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (domain.Game, error) {
	panic("impl me")
}

func (postgresql PostgreSQL) SaveGame(ctx context.Context, game domain.Game) (*GameTransaction, error) {
	panic("impl me")
}
