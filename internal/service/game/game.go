package gameservice

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/domain"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	"github.com/sariya23/game_service/internal/storage/s3"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type KafkaProducer interface {
	SendMessage(message string) error
}

type GameProvider interface {
	GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (game domain.Game, err error)
}

type GameSaver interface {
	SaveGame(ctx context.Context, game domain.Game) (*postgresql.GameTransaction, error)
}

type S3Storager interface {
	Save(ctx context.Context, data io.Reader, key string) (string, error)
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
	gameToAdd *gamev4.Game,
) (uint64, error) {
	const operationPlace = "gameservice.AddGame"
	log := gameService.log.With("operationPlace", operationPlace)
	_, err := gameService.gameProvider.GetGameByTitleAndReleaseYear(ctx, gameToAdd.GetTitle(), gameToAdd.GetReleaseYear().Year)
	if err == nil {
		log.Warn(fmt.Sprintf("game with title=%q and release year=%d already exist", gameToAdd.GetTitle(), gameToAdd.GetReleaseYear().Year))
		return 0, outerror.ErrGameAlreadyExist
	} else {
		log.Error(fmt.Sprintf("cannot get game by title=%q and release year=%d", gameToAdd.GetTitle(), gameToAdd.GetReleaseYear().Year))
	}
	game := domain.Game{
		Title:       gameToAdd.GetTitle(),
		Description: gameToAdd.GetDescription(),
		ReleaseYear: gameToAdd.GetReleaseYear(),
		Tags:        gameToAdd.GetTags(),
		Genres:      gameToAdd.GetGenres(),
	}
	imageURL := s3.GetImageURL(game.Title)
	game.ImageCoverURL = imageURL
	_, err = gameService.gameSaver.SaveGame(ctx, game)
	if err != nil {
		log.Error("cannot start transaction to save game")
		return uint64(0), outerror.ErrCannotStartGameTransaction
	}
	log.Info("game save with PENDING status")
	return uint64(1), nil
}

func (gameService *GameService) GetGame(
	ctx context.Context,
	gameID uint64,
) (*gamev4.GameWithRating, error) {
	panic("impl me")
}

func (gameService *GameService) GetTopGames(
	ctx context.Context,
	gameFilters model.GameFilters,
	limit uint32,
) ([]*gamev4.GameWithRating, error) {
	panic("impl me")
}

func (gameService *GameService) DeleteGame(
	ctx context.Context,
	gameID uint64,
) (*gamev4.Game, error) {
	panic("empl me")
}
