package converters

import (
	game_api "github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/model"
)

func ToShortGameResponse(game model.ShortGame) *game_api.GameListResponse_ShortGame {
	return &game_api.GameListResponse_ShortGame{
		ID:            game.GameID,
		Title:         game.Title,
		Description:   game.Description,
		CoverImageUrl: game.ImageURL,
		ReleaseDate:   ToProtoDate(game.ReleaseDate),
	}
}
