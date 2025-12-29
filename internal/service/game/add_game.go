package gameservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/game_service/internal/interceptors"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
)

func (gameService *GameService) AddGame(
	ctx context.Context,
	gameToAdd dto.AddGameHandler,
) (int64, error) {
	const operationPlace = "gameservice.AddGame"
	log := gameService.log.With("operationPlace", operationPlace)
	log = log.With("title", gameToAdd.Title)
	log = log.With("release_date", gameToAdd.ReleaseDate.String())
	requestID := ctx.Value(interceptors.RequestIDKey).(string)
	log = log.With("request_id", requestID)
	_, err := gameService.gameRepository.GetGameByTitleAndReleaseYear(ctx, gameToAdd.Title, int32(gameToAdd.ReleaseDate.Year()))
	if err == nil {
		log.Warn("game already exists", slog.String("title", gameToAdd.Title), slog.String("release_date", gameToAdd.ReleaseDate.String()))
		return 0, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameAlreadyExist)
	} else if !errors.Is(err, outerror.ErrGameNotFound) {
		log.Error("unexpected error, cannot check game", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s:%w", operationPlace, err)
	}
	var errSaveImage error
	var imageURL string
	if len(gameToAdd.CoverImage) != 0 {
		gameKey := minioclient.GameKey(gameToAdd.Title, gameToAdd.ReleaseDate.Year())
		imageURL, err = gameService.s3Storager.SaveObject(
			ctx,
			gameKey,
			bytes.NewReader(gameToAdd.CoverImage),
		)
		if err != nil {
			log.Error("failed to save image",
				slog.String("game_key", gameKey),
				slog.String("error", err.Error()),
			)
			errSaveImage = outerror.ErrCannotSaveGameImage
		} else {
			log.Info("image successfully saved in s3", slog.String("game_key", gameKey))
		}
	}
	var tagIDs []int64
	if t := gameToAdd.Tags; len(t) != 0 {
		tags, err := gameService.tagReposetory.GetTagByNames(ctx, t)
		if err != nil {
			if errors.Is(err, outerror.ErrTagNotFound) {
				log.Warn("tag doesnt exists", slog.Any("tags", t))
				return 0, fmt.Errorf("%s: %w", operationPlace, outerror.ErrTagNotFound)
			}
			log.Error("cannot check tags, unexpected error", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s: %w", operationPlace, err)
		}
		tagIDs = model.TagIDs(tags)
	}
	var genreIDs []int64
	if g := gameToAdd.Genres; len(g) != 0 {
		genres, err := gameService.genreReposetory.GetGenreByNames(ctx, g)
		if err != nil {
			if errors.Is(err, outerror.ErrGenreNotFound) {
				log.Warn("genre doesnt exists", slog.Any("tags", g))
				return 0, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGenreNotFound)
			}
			log.Error("cannot check genres, unexpected error", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s: %w", operationPlace, err)
		}
		genreIDs = model.GenreIDs(genres)
	}

	addGameService := dto.AddGameService{
		Title:       gameToAdd.Title,
		ReleaseDate: gameToAdd.ReleaseDate,
		Description: gameToAdd.Description,
		TagIDs:      tagIDs,
		GenreIDs:    genreIDs,
		ImageURL:    imageURL,
	}
	gameID, err := gameService.gameRepository.SaveGame(ctx, addGameService)
	if err != nil {
		log.Error("unexpected error, cannot save game", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}
	// Отправка сообщения в кафку
	return gameID, errSaveImage
}
