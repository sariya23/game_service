package model

import (
	"time"
)

// Tag...
type Tag struct {
	TagID   uint64
	TagName string
}

// GetTagNames возвращает имена тэгов.
func GetTagNames(t []*Tag) []string {
	if t == nil {
		return nil
	}
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

// GetGenreNames возвращает имена жанров.
func GetGenreNames(g []*Genre) []string {
	if g == nil {
		return nil
	}
	res := make([]string, 0, len(g))
	for _, v := range g {
		res = append(res, v.GenreName)
	}
	return res
}

// Game...
type Game struct {
	GameID      int64
	Title       string
	Description string
	ReleaseDate time.Time
	ImageURL    string
	Tags        []*Tag
	Genres      []*Genre
	GameStatus  int
}
