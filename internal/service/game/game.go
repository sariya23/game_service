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
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type GameProvider interface {
	GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (game *gamev4.DomainGame, err error)
	GetGameByID(ctx context.Context, gameID uint64) (game *gamev4.DomainGame, err error)
	GetTopGames(ctx context.Context, releaseYear string, tags []string, genres []string, limit uint32) (games []*gamev4.DomainGame, err error)
}

type GameSaver interface {
	SaveGame(ctx context.Context, game *gamev4.DomainGame) (savedGame *gamev4.DomainGame, err error)
}

type S3Storager interface {
	SaveObject(ctx context.Context, name string, data io.Reader) (string, error)
	GetObject(ctx context.Context, name string) (io.Reader, error)
}

type EmailAlerter interface {
	SendMessage(subject string, body string) error
}

type GameService struct {
	log          *slog.Logger
	gameProvider GameProvider
	s3Storager   S3Storager
	gameSaver    GameSaver
	mailer       EmailAlerter
}

func NewGameService(
	log *slog.Logger,
	gameProvider GameProvider,
	s3Storager S3Storager,
	gameSaver GameSaver,
	mailer EmailAlerter,

) *GameService {
	return &GameService{
		log:          log,
		gameProvider: gameProvider,
		s3Storager:   s3Storager,
		gameSaver:    gameSaver,
		mailer:       mailer,
	}
}

func (gameService *GameService) AddGame(
	ctx context.Context,
	gameToAdd *gamev4.GameRequest,
) (*gamev4.DomainGame, error) {
	const operationPlace = "gameservice.AddGame"
	log := gameService.log.With("operationPlace", operationPlace)
	_, err := gameService.gameProvider.GetGameByTitleAndReleaseYear(ctx, gameToAdd.GetTitle(), gameToAdd.GetReleaseYear().Year)
	if err == nil {
		log.Warn(fmt.Sprintf("game with title=%q and release year=%d already exist", gameToAdd.GetTitle(), gameToAdd.GetReleaseYear().Year))
		return nil, fmt.Errorf("%s: %w", operationPlace, outerror.ErrGameAlreadyExist)
	} else if !errors.Is(err, outerror.ErrGameNotFound) {
		log.Error(fmt.Sprintf("cannot get game by title=%q and release year=%d", gameToAdd.GetTitle(), gameToAdd.GetReleaseYear().Year))
		return nil, fmt.Errorf("%s:%w", operationPlace, err)
	}
	var errSaveImage error
	var imageURL string
	if len(gameToAdd.GetCoverImage()) != 0 {
		gameKey := fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseYear().Year))
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
	game := gamev4.DomainGame{
		Title:         gameToAdd.GetTitle(),
		Description:   gameToAdd.GetDescription(),
		ReleaseYear:   gameToAdd.GetReleaseYear(),
		Tags:          gameToAdd.GetTags(),
		Genres:        gameToAdd.GetGenres(),
		CoverImageUrl: imageURL,
	}
	savedGame, err := gameService.gameSaver.SaveGame(ctx, &game)
	if err != nil {
		log.Error(fmt.Sprintf("cannot save game: err = %v", fmt.Errorf("%w: %w", errSaveImage, err)))
		return nil, fmt.Errorf("%s: %w", operationPlace, err)
	}
	log.Info("game save successfully")
	err = gameService.mailer.SendMessage(
		"Добавлена игра",
		fmt.Sprintf("Добавлена игра %s %d года", savedGame.Title, savedGame.GetReleaseYear().Year),
	)
	if err != nil {
		log.Warn(fmt.Sprintf("cannot send alert; err = %v", err))
	}
	return savedGame, errSaveImage
}

func (gameService *GameService) GetGame(
	ctx context.Context,
	gameID uint64,
) (*gamev4.DomainGame, error) {
	const operationPlace = "gameservice.GetGame"
	log := gameService.log.With("operationPlace", operationPlace)
	game, err := gameService.gameProvider.GetGameByID(ctx, gameID)
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
) ([]*gamev4.DomainGame, error) {
	const operationPlace = "gameservice.GetTopGames"
	log := gameService.log.With("operationPlace", operationPlace)
	if gameFilters.ReleaseYear == "" {
		gameFilters.ReleaseYear = fmt.Sprint(time.Now().Year())
	}
	games, err := gameService.gameProvider.GetTopGames(ctx, gameFilters.ReleaseYear, gameFilters.Tags, gameFilters.Tags, limit)
	if err != nil {
		log.Error(fmt.Sprintf("unexcpected error; err=%v", err))
		return []*gamev4.DomainGame{}, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return games, nil
}

func (gameService *GameService) DeleteGame(
	ctx context.Context,
	gameID uint64,
) (*gamev4.DomainGame, error) {
	panic("empl me")
}
