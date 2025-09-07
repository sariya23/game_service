package gameservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type GameReposetory interface {
	GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (*model.Game, error)
	GetGameByID(ctx context.Context, gameID uint64) (*model.Game, error)
	GetTopGames(ctx context.Context, filters model.GameFilters, limit uint32) ([]model.Game, error)
	SaveGame(ctx context.Context, game model.Game) error
	DaleteGame(ctx context.Context, gameID uint64) (*model.Game, error)
}

type TagRepository interface {
	GetTags(ctx context.Context, tags []string) ([]model.Tag, error)
}

type GenreRepository interface {
	GetGenres(ctx context.Context, genres []string) ([]model.Genre, error)
}

type S3Storager interface {
	SaveObject(ctx context.Context, name string, data io.Reader) (string, error)
	GetObject(ctx context.Context, name string) (io.Reader, error)
	DeleteObject(ctx context.Context, name string) error
}

type EmailAlerter interface {
	SendMessage(subject string, body string) error
}

type GameService struct {
	log             *slog.Logger
	gameRepository  GameReposetory
	tagReposetory   TagRepository
	genreReposetory GenreRepository
	s3Storager      S3Storager
	mailer          EmailAlerter
}

func NewGameService(
	log *slog.Logger,
	gameReposiroy GameReposetory,
	tagReposetory TagRepository,
	genreReposetory GenreRepository,
	s3Storager S3Storager,
	mailer EmailAlerter,

) *GameService {
	return &GameService{
		log:             log,
		s3Storager:      s3Storager,
		tagReposetory:   tagReposetory,
		genreReposetory: genreReposetory,
		gameRepository:  gameReposiroy,
		mailer:          mailer,
	}
}

func (gameService *GameService) AddGame(
	ctx context.Context,
	gameToAdd *gamev4.GameRequest,
) error {
	const operationPlace = "gameservice.AddGame"
	log := gameService.log.With("operationPlace", operationPlace)
	_, err := gameService.gameRepository.GetGameByTitleAndReleaseYear(ctx, gameToAdd.GetTitle(), gameToAdd.GetReleaseDate().Year)
	if err == nil {
		log.Warn(fmt.Sprintf("game with title=%q and release year=%d already exist", gameToAdd.GetTitle(), gameToAdd.GetReleaseDate().Year))
		return fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameAlreadyExist)
	} else if !errors.Is(err, outerror.ErrGameNotFound) {
		log.Error(fmt.Sprintf("cannot get game by title=%q and release year=%d", gameToAdd.GetTitle(), gameToAdd.GetReleaseDate().Year))
		return fmt.Errorf("%s:%w", operationPlace, err)
	}
	var errSaveImage error
	var imageURL string
	if len(gameToAdd.GetCoverImage()) != 0 {
		gameKey := minioclient.GameKey(gameToAdd.GetTitle(), int(gameToAdd.ReleaseDate.GetYear()))
		imageURL, err = gameService.s3Storager.SaveObject(
			ctx,
			gameKey,
			bytes.NewReader(gameToAdd.GetCoverImage()),
		)
		if err != nil {
			log.Error(fmt.Sprintf("cannot save game cover image (title=%s) in s3; err = %v", gameKey, err))
			errSaveImage = outerror.ErrCannotSaveGameImage
		} else {
			log.Info(fmt.Sprintf("image successfully saved in s3 with key=%s", gameKey))
		}
	}
	log.Info("no image data in game")
	var tags []model.Tag
	if t := gameToAdd.GetTags(); len(t) != 0 {
		tags, err = gameService.tagReposetory.GetTags(ctx, t)
		if err != nil {
			if errors.Is(err, outerror.ErrTagNotFound) {
				log.Warn("tags with this names not found", slog.String("tags", fmt.Sprintf("%#v", t)))
				return fmt.Errorf("%s: %w", operationPlace, outerror.ErrTagNotFound)
			}
		}
	}
	var genres []model.Genre
	if g := gameToAdd.GetGenres(); len(g) != 0 {
		genres, err = gameService.genreReposetory.GetGenres(ctx, g)
		if err != nil {
			if errors.Is(err, outerror.ErrGenreNotFound) {
				log.Warn("genres with this names not found", slog.String("genres", fmt.Sprintf("%#v", g)))
				return fmt.Errorf("%s: %w", operationPlace, outerror.ErrGenreNotFound)
			}
		}
	}
	game := model.Game{
		Title:       gameToAdd.GetTitle(),
		Description: gameToAdd.GetDescription(),
		ReleaseDate: time.Date(int(gameToAdd.ReleaseDate.Year), time.Month(gameToAdd.ReleaseDate.Month), int(gameToAdd.ReleaseDate.Day), 0, 0, 0, 0, time.UTC),
		Tags:        tags,
		Genres:      genres,
		ImageURL:    imageURL,
	}
	err = gameService.gameRepository.SaveGame(ctx, game)
	if err != nil {
		log.Error(fmt.Sprintf("cannot save game: err = %v", fmt.Errorf("%w: %w", errSaveImage, err)))
		return fmt.Errorf("%s: %w", operationPlace, err)
	}
	log.Info("game save successfully")
	err = gameService.mailer.SendMessage(
		"Добавлена игра",
		fmt.Sprintf("Добавлена игра %s %d года", game.Title, game.ReleaseDate.Year()),
	)
	if err != nil {
		log.Warn(fmt.Sprintf("cannot send alert; err = %v", err))
	}
	return errSaveImage
}

func (gameService *GameService) GetGame(
	ctx context.Context,
	gameID uint64,
) (*model.Game, error) {
	const operationPlace = "gameservice.GetGame"
	log := gameService.log.With("operationPlace", operationPlace)
	game, err := gameService.gameRepository.GetGameByID(ctx, gameID)
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			log.Warn(fmt.Sprintf("game with id=%d not found", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameNotFound)
		}
		log.Error(fmt.Sprintf("unexpected error; err=%v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return game, nil
}

func (gameService *GameService) GetTopGames(
	ctx context.Context,
	gameFilters model.GameFilters,
	limit uint32,
) ([]model.Game, error) {
	const operationPlace = "gameservice.GetTopGames"
	log := gameService.log.With("operationPlace", operationPlace)
	if gameFilters.ReleaseYear == 0 {
		gameFilters.ReleaseYear = int32(time.Now().Year())
	}
	games, err := gameService.gameRepository.GetTopGames(ctx, gameFilters, limit)
	if err != nil {
		log.Error(fmt.Sprintf("unexcpected error; err=%v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return games, nil
}

func (gameService *GameService) DeleteGame(
	ctx context.Context,
	gameID uint64,
) (*model.Game, error) {
	const operationPlace = "gameservice.DeleteGame"
	log := gameService.log.With("operationPlace", operationPlace)
	deletedGame, err := gameService.gameRepository.DaleteGame(ctx, gameID)
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			log.Warn(fmt.Sprintf("game with id=%v not found", gameID))
			return nil, fmt.Errorf("%s: %w", operationPlace, err)
		}
		log.Error(fmt.Sprintf("unexpected error; err=%v", err))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	log.Info(fmt.Sprintf("game with id=%v deleted from DB", gameID))
	gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseDate.Year()))
	err = gameService.s3Storager.DeleteObject(ctx, gameKey)
	var errDeleteImage error
	if err != nil {
		if errors.Is(err, outerror.ErrImageNotFoundS3) {
			log.Info("there is not image for game")
		} else {
			log.Error("cannot delete image from s3")
			log.Info(fmt.Sprintf("game: %+v", deletedGame))
		}
		errDeleteImage = err
	} else {
		log.Info("image cover deleted from s3")
	}
	return deletedGame, errDeleteImage
}
