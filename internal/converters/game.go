package converters

import (
	"github.com/sariya23/game_service/internal/model"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/genproto/googleapis/type/date"
)

// ToGame...
func ToProtoGame(modelGame model.Game) *gamev4.DomainGame {
	game := gamev4.DomainGame{}
	game.Title = modelGame.Title
	game.Description = modelGame.Description
	game.CoverImageUrl = modelGame.ImageURL
	game.ID = modelGame.GameID
	if len(modelGame.Genres) > 0 {
		genres := make([]string, 0, len(modelGame.Genres))
		for _, g := range modelGame.Genres {
			genres = append(genres, g.GenreName)
		}
		game.Genres = genres
	}

	if len(modelGame.Tags) > 0 {
		tags := make([]string, 0, len(modelGame.Genres))
		for _, t := range modelGame.Tags {
			tags = append(tags, t.TagName)
		}
		game.Tags = tags
	}
	game.ReleaseDate = &date.Date{
		Year:  int32(modelGame.ReleaseDate.Year()),
		Month: int32(modelGame.ReleaseDate.Month()),
		Day:   int32(modelGame.ReleaseDate.Day()),
	}
	return &game
}
