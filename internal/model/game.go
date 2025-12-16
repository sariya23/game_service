package model

import (
	"time"

	"github.com/sariya23/api_game_service/gen/game"
)

type Game struct {
	GameID      int64
	Title       string
	Description string
	ReleaseDate time.Time
	ImageURL    string
	Tags        []Tag
	Genres      []Genre
	GameStatus  game.GameStatusType
}
