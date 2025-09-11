package random

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/model"
)

func NewRandomGame() *model.Game {
	var res model.Game
	fakeit := gofakeit.New(0)
	res.GameID = gofakeit.Uint64()
	res.Title = fakeit.Book().Title
	res.Description = fakeit.LetterN(20)
	res.ReleaseDate = fakeit.Date()
	res.ImageURL = fakeit.URL()
	for i := 0; i < fakeit.IntN(4); i++ {
		res.Tags = append(res.Tags, *NewRandomTag())
	}
	for i := 0; i < fakeit.IntN(4); i++ {
		res.Genres = append(res.Genres, *NewRandomGenre())
	}
	return &res
}
