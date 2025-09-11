package model

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// Tag...
type Tag struct {
	TagID   uint64
	TagName string
}

// TagNames возвращает имена тэгов.
func TagNames(t []Tag) []string {
	res := make([]string, 0, len(t))
	for _, v := range t {
		res = append(res, v.TagName)
	}
	return res
}

func NewRandomTag() *Tag {
	var tag Tag
	fakeit := gofakeit.New(0)
	tag.TagID = fakeit.Uint64()
	tag.TagName = fakeit.Book().Genre
	return &tag
}

// Genre...
type Genre struct {
	GenreID   uint64
	GenreName string
}

// GenreNames возвращает имена жанров.
func GenreNames(g []Genre) []string {
	res := make([]string, 0, len(g))
	for _, v := range g {
		res = append(res, v.GenreName)
	}
	return res
}

func NewRandomGenre() *Genre {
	var genre Genre
	fakeit := gofakeit.New(0)
	genre.GenreID = fakeit.Uint64()
	genre.GenreName = fakeit.Book().Genre
	return &genre
}

// Game...
type Game struct {
	GameID      uint64
	Title       string
	Description string
	ReleaseDate time.Time
	ImageURL    string
	Tags        []Tag
	Genres      []Genre
}

func NewRandomGame() *Game {
	var res Game
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
