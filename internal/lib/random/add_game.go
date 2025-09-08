package random

import (
	"github.com/brianvoe/gofakeit/v7"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/genproto/googleapis/type/date"
)

// RandomAddGameRequest случайная игра для запроса на создание
func RandomAddGameRequest() *gamev4.GameRequest {
	var res gamev4.GameRequest
	fakeit := gofakeit.New(0)
	res.Title = fakeit.Book().Title
	res.Description = fakeit.LetterN(20)
	randomDate := fakeit.Date()
	res.ReleaseDate = &date.Date{
		Year:  int32(randomDate.Year()),
		Month: int32(randomDate.Month()),
		Day:   int32(randomDate.Day()),
	}
	res.CoverImage = []byte(fakeit.URL())
	var tags []string
	var genres []string
	fakeit.Slice(&tags)
	fakeit.Slice(&genres)
	res.Tags = tags
	res.Genres = genres
	return &res
}
