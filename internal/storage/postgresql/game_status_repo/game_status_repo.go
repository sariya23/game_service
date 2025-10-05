package gamestatusrepo

import (
	"log/slog"

	"github.com/sariya23/game_service/internal/storage/db"
)

const (
	GameStatusTable                 = "game_status"
	GameStatusGameStatusIDFieldName = "game_status_id"
	GameStatusGameNameFieldName     = "name"
)

type GameStatusRepository struct {
	conn *db.Database
	log  *slog.Logger
}

func NewGameStatusRepository(conn *db.Database, log *slog.Logger) *GameStatusRepository {
	return &GameStatusRepository{
		conn: conn,
		log:  log,
	}
}
