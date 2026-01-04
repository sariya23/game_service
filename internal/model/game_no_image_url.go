package model

import (
	"time"

	"github.com/sariya23/api_game_service/gen/game"
)

type GameNoImageURL struct {
	GameID      int64
	Title       string
	Description string
	ReleaseDate time.Time
	ImageKey    string
	Tags        []Tag
	Genres      []Genre
	GameStatus  game.GameStatusType
}

func (g GameNoImageURL) ToDomain(imageURL string) Game {
	return Game{
		GameID:      g.GameID,
		Title:       g.Title,
		Description: g.Description,
		ReleaseDate: g.ReleaseDate,
		ImageURL:    imageURL,
		Tags:        g.Tags,
		Genres:      g.Genres,
		GameStatus:  g.GameStatus,
	}
}
