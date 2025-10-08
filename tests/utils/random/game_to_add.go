//go:build integrations

package random

import (
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/lib/converters"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

func GameToAddService(genresIDs, tagIDs []int64) dto.AddGameService {
	var game dto.AddGameService

	game.Title = strings.ToLower(gofakeit.LetterN(20))
	game.Description = gofakeit.Sentence(50)
	game.ReleaseDate = time.Date(gofakeit.Year(), time.Month(gofakeit.Month()), gofakeit.Day(), 0, 0, 0, 0, time.UTC)
	game.ImageURL = gofakeit.URL()
	game.GenreIDs = Sample(genresIDs, 2)
	game.TagIDs = Sample(tagIDs, 3)
	return game
}

func GameToAddRequest(genres, tags []string) gamev2.GameRequest {
	var game gamev2.GameRequest
	game.Title = strings.ToLower(gofakeit.LetterN(20))
	game.Description = gofakeit.Sentence(50)
	game.ReleaseDate = converters.ToProtoDate(time.Date(gofakeit.Year(), time.Month(gofakeit.Month()), gofakeit.Day(), 0, 0, 0, 0, time.UTC))
	img, err := Image()
	if err != nil {
		panic(err)
	}
	game.CoverImage = img
	game.Genres = Sample(genres, 2)
	game.Tags = Sample(tags, 3)
	return game
}
