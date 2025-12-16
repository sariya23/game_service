package converters

import (
	game_api "github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/model"
)

func ToProtoGame(modelGame *model.Game) *game_api.DomainGame {
	game := game_api.DomainGame{}
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
	game.ReleaseDate = ToProtoDate(modelGame.ReleaseDate)
	return &game
}
