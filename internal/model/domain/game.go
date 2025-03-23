package domain

import "google.golang.org/genproto/googleapis/type/date"

type Game struct {
	Title         string
	Description   string
	ImageCoverURL string
	Genres        []string
	Tags          []string
	ReleaseYear   *date.Date
}
