package gameservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
)

func (gameService *GameService) AddGame(
	ctx context.Context,
	gameToAdd *dto.AddGame,
) (int64, error) {
	const operationPlace = "gameservice.AddGame"
	log := gameService.log.With("operationPlace", operationPlace)
	_, err := gameService.gameRepository.GetGameByTitleAndReleaseYear(ctx, gameToAdd.Title, gameToAdd.ReleaseDate.Year)
	if err == nil {
		log.Warn(fmt.Sprintf("game with title=%q and release year=%d already exist", gameToAdd.Title, gameToAdd.ReleaseDate.Year))
		return 0, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameAlreadyExist)
	} else if !errors.Is(err, outerror.ErrGameNotFound) {
		log.Error(fmt.Sprintf("cannot get game by title=%q and release year=%d", gameToAdd.Title, gameToAdd.ReleaseDate.Year))
		return 0, fmt.Errorf("%s:%w", operationPlace, err)
	}
	var errSaveImage error
	var imageURL string
	if len(gameToAdd.CoverImage) != 0 {
		gameKey := minioclient.GameKey(gameToAdd.Title, int(gameToAdd.ReleaseDate.GetYear()))
		imageURL, err = gameService.s3Storager.SaveObject(
			ctx,
			gameKey,
			bytes.NewReader(gameToAdd.CoverImage),
		)
		if err != nil {
			log.Error(fmt.Sprintf("cannot save game cover image (title=%s) in s3; err = %v", gameKey, err))
			errSaveImage = outerror.ErrCannotSaveGameImage
		} else {
			log.Info(fmt.Sprintf("image successfully saved in s3 with key=%s", gameKey))
		}
	} else {
		log.Info("no image data in game")
	}
	var tags []*model.Tag
	if t := gameToAdd.Tags; len(t) != 0 {
		tags, err = gameService.tagReposetory.GetTagByNames(ctx, t)
		if err != nil {
			if errors.Is(err, outerror.ErrTagNotFound) {
				log.Warn("tags with this names not found", slog.String("tags", fmt.Sprintf("%#v", t)))
				return 0, fmt.Errorf("%s: %w", operationPlace, outerror.ErrTagNotFound)
			}
			log.Error(fmt.Sprintf("cannot get tags, err=%v", err))
			return 0, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}
	var genres []*model.Genre
	if g := gameToAdd.Genres; len(g) != 0 {
		genres, err = gameService.genreReposetory.GetGenreByNames(ctx, g)
		if err != nil {
			if errors.Is(err, outerror.ErrGenreNotFound) {
				log.Warn("genres with this names not found", slog.String("genres", fmt.Sprintf("%#v", g)))
				return 0, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGenreNotFound)
			}
			log.Error(fmt.Sprintf("cannot get genres, err=%v", err))
			return 0, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}
	game := model.Game{
		Title:       gameToAdd.Title,
		Description: gameToAdd.Description,
		ReleaseDate: time.Date(int(gameToAdd.ReleaseDate.Year), time.Month(gameToAdd.ReleaseDate.Month), int(gameToAdd.ReleaseDate.Day), 0, 0, 0, 0, time.UTC),
		Tags:        tags,
		Genres:      genres,
		ImageURL:    imageURL,
	}
	gameID, err := gameService.gameRepository.SaveGame(ctx, &game)
	if err != nil {
		log.Error(fmt.Sprintf("cannot save game: err = %v", fmt.Errorf("%w: %w", errSaveImage, err)))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}
	// Отправка сообщения в кафку
	log.Info("game save successfully")
	return gameID, errSaveImage
}
