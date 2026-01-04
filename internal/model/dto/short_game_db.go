package dto

import (
	"database/sql"
	"time"

	"github.com/sariya23/game_service/internal/model"
)

type ShortGameDB struct {
	GameID             int64
	Title, Description string
	ReleaseDate        time.Time
	ImageKey           sql.NullString
}

func (s *ShortGameDB) ToDomain() model.ShortGame {
	var imgURL string
	if s.ImageKey.Valid {
		imgURL = s.ImageKey.String
	}
	return model.ShortGame{
		GameID:      s.GameID,
		Title:       s.Title,
		Description: s.Description,
		ReleaseDate: s.ReleaseDate,
		ImageURL:    imgURL,
	}
}
