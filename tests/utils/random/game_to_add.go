//go:build integrations

package random

import (
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	game_api "github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/lib/converters"
	"github.com/sariya23/game_service/internal/model/dto"
)

type GameToAddFields int

const (
	All GameToAddFields = iota
	OnlyRequired
	NoTitle
	NoDescription
	NoReleaseDate
	NoImage
	NoTags
	NoGenres
)

func GameToAddService(genresIDs, tagIDs []int64) dto.AddGameService {
	var game dto.AddGameService

	game.Title = strings.ToLower(gofakeit.LetterN(20))
	game.Description = gofakeit.Sentence(50)
	game.ReleaseDate = time.Date(gofakeit.Year(), time.Month(gofakeit.Month()), gofakeit.Day(), 0, 0, 0, 0, time.UTC)
	game.ImageKey = gofakeit.UUID()
	game.GenreIDs = Sample(genresIDs, 2)
	game.TagIDs = Sample(tagIDs, 3)
	return game
}

func GameToAddRequest(genres, tags []string) *game_api.GameRequest {
	var game game_api.GameRequest
	game.Title = strings.ToLower(gofakeit.LetterN(20))
	game.Description = gofakeit.Sentence(50)
	game.ReleaseDate = converters.ToProtoDate(time.Date(gofakeit.Year(), time.Month(gofakeit.Month()), gofakeit.Day(), 0, 0, 0, 0, time.UTC))
	img, err := Image()
	if err != nil {
		panic(err)
	}
	game.CoverImage = img
	game.Genres = Sample(genres, 5)
	game.Tags = Sample(tags, 5)
	return &game
}
