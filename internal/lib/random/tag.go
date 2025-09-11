package random

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/model"
)

func NewRandomTag() *model.Tag {
	var tag model.Tag
	fakeit := gofakeit.New(0)
	tag.TagID = fakeit.Uint64()
	tag.TagName = fakeit.Book().Genre
	return &tag
}
