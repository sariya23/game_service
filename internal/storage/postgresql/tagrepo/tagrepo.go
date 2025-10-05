package tagrepo

import (
	"log/slog"

	"github.com/sariya23/game_service/internal/storage/db"
)

const (
	TagTagIDFieldName   = "tag_id"
	TagTagNameFieldName = "tag_name"
)

type TagRepository struct {
	conn *db.Database
	log  *slog.Logger
}

func NewTagRepository(conn *db.Database, log *slog.Logger) *TagRepository {
	return &TagRepository{
		conn: conn,
		log:  log,
	}
}
