package model

import "time"

// Tag...
type Tag struct {
	TagID   uint64
	TagName string
}

// Genre...
type Genre struct {
	GenreID   uint64
	GenreName string
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
