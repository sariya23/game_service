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

func (sh ShortGameDB) ToShortGameNoImageURL() model.ShortGameNoImageURL {
	var imgKey string
	if sh.ImageKey.Valid {
		imgKey = sh.ImageKey.String
	}
	return model.ShortGameNoImageURL{
		GameID:      sh.GameID,
		ImageKey:    imgKey,
		Description: sh.Description,
		ReleaseDate: sh.ReleaseDate,
		Title:       sh.Title,
	}
}
