package model

import "time"

// ShortGame...
type ShortGame struct {
	GameID             uint64
	Title, Description string
	ReleaseDate        time.Time
	ImageURL           string
}
