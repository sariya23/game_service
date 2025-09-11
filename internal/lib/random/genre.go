package random

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/model"
)

func NewRandomGenre() *model.Genre {
	var genre model.Genre
	fakeit := gofakeit.New(0)
	genre.GenreID = fakeit.Uint64()
	genre.GenreName = fakeit.Book().Genre
	return &genre
}
