package model

import (
	"time"

	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

// Game...
type Game struct {
	GameID      int64
	Title       string
	Description string
	ReleaseDate time.Time
	ImageURL    string
	Tags        []Tag
	Genres      []Genre
	GameStatus  gamev2.GameStatusType
}
