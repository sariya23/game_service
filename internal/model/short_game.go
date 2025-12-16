package model

import "time"

type ShortGame struct {
	GameID             int64
	Title, Description string
	ReleaseDate        time.Time
	ImageURL           string
}
