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
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type KafkaProducer interface {
	SendMessage(message string) error
}

type GameProvider interface {
	GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (game *gamev4.DomainGame, err error)
}

type GameSaver interface {
	SaveGame(ctx context.Context, game *gamev4.DomainGame) (savedGame *gamev4.DomainGame, err error)
}

type S3Storager interface {
	Save(ctx context.Context, data io.Reader, key string) (url string, err error)
	Get(ctx context.Context, bucket, key string) io.Reader
}

type GameService struct {
	log           *slog.Logger
	kafkaProducer KafkaProducer
	gameProvider  GameProvider
	s3Storager    S3Storager
	gameSaver     GameSaver
}

func NewGameService(
	log *slog.Logger,
	kafkaProducer KafkaProducer,
	gameProvider GameProvider,
	s3Storager S3Storager,
	gameSaver GameSaver,
) *GameService {
	return &GameService{
		log:           log,
		kafkaProducer: kafkaProducer,
		gameProvider:  gameProvider,
		s3Storager:    s3Storager,
		gameSaver:     gameSaver,
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
		imageURL, err = gameService.s3Storager.Save(
			ctx,
			bytes.NewReader(gameToAdd.CoverImage),
			fmt.Sprintf("%s_%d_image", gameToAdd.Title, gameToAdd.ReleaseYear.Year),
		)
		if err != nil {
			log.Error(fmt.Sprintf("cannot save game cover image in s3; err = %v", err))
			errSaveImage = outerror.ErrCannotSaveGameImage
		} else {
			log.Info("image successfully saved in s3")
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

	return savedGame, nil
}

func (gameService *GameService) GetGame(
	ctx context.Context,
	gameID uint64,
) (*gamev4.DomainGame, error) {
	panic("impl me")
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
