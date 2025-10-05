package dto

import "github.com/sariya23/game_service/internal/model"

type DeletedGame struct {
	GameID      int64
	Title       string
	ReleaseYear uint64
}

func DeletedGameFromGame(game model.Game) *DeletedGame {
	return &DeletedGame{
		GameID:      game.GameID,
		ReleaseYear: uint64(game.ReleaseDate.Year()),
		Title:       game.Title,
	}
}
