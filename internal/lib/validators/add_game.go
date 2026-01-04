package validators

import (
	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/outerror"
)

func AddGame(request *game.AddGameRequest) (valid bool, message string) {
	if request.Game == nil {
		return false, outerror.EmptyRequestMessage
	}
	if request.Game.Title == "" {
		return false, outerror.TitleRequiredMessage
	}
	if request.Game.Description == "" {
		return false, outerror.DescriptionRequiredMessage
	}
	if request.Game.ReleaseDate == nil {
		return false, outerror.ReleaseDateRequiredMessage
	}
	if request.Game.ReleaseDate.Year == 0 || request.Game.ReleaseDate.Month == 0 || request.Game.ReleaseDate.Day == 0 {
		return false, outerror.ReleaseDateRequiredMessage
	}
	return true, ""
}
