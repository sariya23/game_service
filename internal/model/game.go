package model

import "time"

type Game struct {
	Title       string
	Genre       string
	Description string
	ReleaseYear time.Time
	CoverImage  []byte
	Tags        []string
}
