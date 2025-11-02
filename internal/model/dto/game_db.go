package dto

import (
	"database/sql"
	"time"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/model"
)

type GameDB struct {
	GameID      int64
	Title       string
	Description string
	ReleaseDate time.Time
	ImageURL    sql.NullString
	GameStatus  game.GameStatusType
}

func (g GameDB) ToDomain() model.Game {
	var imgURL string
	if g.ImageURL.Valid {
		imgURL = g.ImageURL.String
	}
	return model.Game{
		GameID:      g.GameID,
		Title:       g.Title,
		Description: g.Description,
		ReleaseDate: g.ReleaseDate,
		GameStatus:  g.GameStatus,
		ImageURL:    imgURL,
	}
}
