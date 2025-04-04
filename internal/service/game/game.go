package gameservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/game_service/internal/storage/s3"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type GameProvider interface {
	GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (game *gamev4.DomainGame, err error)
	GetGameByID(ctx context.Context, gameID uint64) (game *gamev4.DomainGame, err error)
}

type GameSaver interface {
	SaveGame(ctx context.Context, game *gamev4.DomainGame) (savedGame *gamev4.DomainGame, err error)
}

type S3Storager interface {
	Save(ctx context.Context, data io.Reader, key string) (url string, err error)
	Get(ctx context.Context, bucket, key string) io.Reader
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
		return nil, outerror.ErrGameAlreadyExist
	} else if !errors.Is(err, outerror.ErrGameNotFound) {
		log.Error(fmt.Sprintf("cannot get game by title=%q and release year=%d", gameToAdd.GetTitle(), gameToAdd.GetReleaseYear().Year))
		return nil, err
	}
	var errSaveImage error
	var imageURL string
	if len(gameToAdd.GetCoverImage()) != 0 {
		gameKey := s3.CreateGameKey(gameToAdd.GetTitle(), int(gameToAdd.GetReleaseYear().Year))
		imageURL, err = gameService.s3Storager.Save(
			ctx,
			bytes.NewReader(gameToAdd.GetCoverImage()),
			gameKey,
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
		return nil, err
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
			return nil, outerror.ErrGameNotFound
		}
		log.Error(fmt.Sprintf("unexpected error; err=%v", err))
		return nil, err
	}
	return game, nil
}

func (gameService *GameService) GetTopGames(
	ctx context.Context,
	gameFilters model.GameFilters,
	limit uint32,
) ([]*gamev4.DomainGame, error) {
	panic("impl me")
}

func (gameService *GameService) DeleteGame(
	ctx context.Context,
	gameID uint64,
) (*gamev4.DomainGame, error) {
	panic("empl me")
}
