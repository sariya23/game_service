package validators

import (
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

func AddGame(request *gamev2.AddGameRequest) (valid bool, message string) {
	if request.Game.Title == "" {
		return false, outerror.TitleRequiredMessage
	}
	if request.Game.Description == "" {
		return false, outerror.DescriptionRequiredMessage
	}
	if request.Game.ReleaseDate == nil {
		return false, outerror.ReleaseYearRequiredMessage
	}
	if request.Game.ReleaseDate.Year == 0 || request.Game.ReleaseDate.Month == 0 || request.Game.ReleaseDate.Day == 0 {
		return false, outerror.ReleaseYearRequiredMessage
	}
	return true, ""
}
