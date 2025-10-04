package dto

import (
	"google.golang.org/genproto/googleapis/type/date"
)

type AddGame struct {
	Title       string
	Genres      []string
	Description string
	ReleaseDate *date.Date
	CoverImage  []byte
	Tags        []string
}
