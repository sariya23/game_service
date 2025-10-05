//go:build integrations

package random

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/model/dto"
)

func GameToAddService(genresIDs, tagIDs []int64) dto.AddGameService {
	var game dto.AddGameService

	game.Title = gofakeit.LetterN(20)
	game.Description = gofakeit.Sentence(50)
	game.ReleaseDate = gofakeit.Date()
	game.ImageURL = gofakeit.URL()
	game.GenreIDs = Sample(genresIDs, 2)
	game.TagIDs = Sample(tagIDs, 3)
	return game
}
