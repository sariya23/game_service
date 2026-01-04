package dto

import (
	"time"
)

type AddGameHandler struct {
	Title       string
	Genres      []string
	Description string
	ReleaseDate time.Time
	CoverImage  []byte
	Tags        []string
}

type AddGameService struct {
	Title       string
	GenreIDs    []int64
	Description string
	ReleaseDate time.Time
	ImageKey    string
	TagIDs      []int64
}
