package gameservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/sariya23/game_service/internal/outerror"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
)

func (gameService *GameService) DeleteGame(
	ctx context.Context,
	gameID int64,
) (int64, error) {
	const operationPlace = "gameservice.DeleteGame"
	log := gameService.log.With("operationPlace", operationPlace)
	deletedGame, err := gameService.gameRepository.DaleteGame(ctx, gameID)
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			log.Warn(fmt.Sprintf("game with id=%v not found", gameID))
			return 0, fmt.Errorf("%s: %w", operationPlace, err)
		}
		log.Error(fmt.Sprintf("unexpected error; err=%v", err))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}
	log.Info(fmt.Sprintf("game with id=%v deleted from DB", gameID))
	gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseYear))
	err = gameService.s3Storager.DeleteObject(ctx, gameKey)
	var errDeleteImage error
	if err != nil {
		log.Error("cannot delete image from s3")
		log.Info(fmt.Sprintf("game: %+v", deletedGame))
		errDeleteImage = err
	} else {
		log.Info("image cover deleted from s3 or is not present")
	}
	return deletedGame.GameID, errDeleteImage
}
