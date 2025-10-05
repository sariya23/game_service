package random

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/model"
)

func NewRandomShortGame() *model.ShortGame {
	var res model.ShortGame
	fakeit := gofakeit.New(0)
	res.GameID = gofakeit.Int64()
	res.Title = fakeit.Book().Title
	res.Description = fakeit.LetterN(20)
	res.ReleaseDate = fakeit.Date()
	res.ImageURL = fakeit.URL()
	return &res
}
