package converters

import (
	"github.com/sariya23/game_service/internal/model"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

func ToShortGameResponse(game model.ShortGame) *gamev4.GetTopGamesResponse_ShortGame {
	return &gamev4.GetTopGamesResponse_ShortGame{
		ID:            game.GameID,
		Title:         game.Title,
		Description:   game.Description,
		CoverImageUrl: game.ImageURL,
		ReleaseDate:   ToProtoDate(game.ReleaseDate),
	}
}
