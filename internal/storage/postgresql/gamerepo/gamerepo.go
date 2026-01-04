package gamerepo

import (
	"log/slog"

	"github.com/sariya23/game_service/internal/storage/db"
)

const (
	GameGameIDFieldName       = "game_id"
	GameTitleFieldName        = "title"
	GameDescriptionFieldName  = "description"
	GameReleaseDateFieldName  = "release_date"
	GameImageURLFieldName     = "image_key"
	GameGameStatusIDFieldName = "game_status_id"
)

type GameRepository struct {
	conn *db.Database
	log  *slog.Logger
}

func NewGameRepository(conn *db.Database, log *slog.Logger) *GameRepository {
	return &GameRepository{
		conn: conn,
		log:  log,
	}
}
