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
