package model

import "time"

type ShortGameNoImageURL struct {
	GameID             int64
	Title, Description string
	ReleaseDate        time.Time
	ImageKey           string
}

func (sh ShortGameNoImageURL) ToShortGame(imageURL string) ShortGame {
	return ShortGame{
		GameID:      sh.GameID,
		Title:       sh.Title,
		Description: sh.Description,
		ReleaseDate: sh.ReleaseDate,
		ImageURL:    imageURL,
	}
}
