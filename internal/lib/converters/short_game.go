package converters

import (
	"github.com/sariya23/game_service/internal/model"
	gamev2 "github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

func ToShortGameResponse(game model.ShortGame) *gamev2.GameListResponse_ShortGame {
	return &gamev2.GameListResponse_ShortGame{
		ID:            game.GameID,
		Title:         game.Title,
		Description:   game.Description,
		CoverImageUrl: game.ImageURL,
		ReleaseDate:   ToProtoDate(game.ReleaseDate),
	}
}
