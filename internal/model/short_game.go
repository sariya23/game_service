package model

import "time"

// ShortGame...
type ShortGame struct {
	GameID             int64
	Title, Description string
	ReleaseDate        time.Time
	ImageURL           string
}
