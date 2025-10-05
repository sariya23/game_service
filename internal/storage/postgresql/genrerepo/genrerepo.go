package genrerepo

import (
	"log/slog"

	"github.com/sariya23/game_service/internal/storage/db"
)

const (
	GenreGenreIDFieldName   = "genre_id"
	GenreGenreNameFieldName = "genre_name"
)

type GenreRepository struct {
	conn *db.Database
	log  *slog.Logger
}

func NewGenreRepository(conn *db.Database, log *slog.Logger) *GenreRepository {
	return &GenreRepository{
		conn: conn,
		log:  log,
	}
}
