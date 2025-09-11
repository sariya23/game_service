package model

import (
	"time"
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
