package gameservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type KafkaProducer interface {
	SendMessage(message string) error
}

type GameProvider interface {
	GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (game gamev4.GameWithRating, err error)
}

type GameService struct {
	log           *slog.Logger
	kafkaProducer KafkaProducer
	gameProvider  GameProvider
}

func NewGameService(log *slog.Logger, kafkaProducer KafkaProducer, gameProvider GameProvider) *GameService {
	return &GameService{
		log:           log,
		kafkaProducer: kafkaProducer,
		gameProvider:  gameProvider,
	}
}

func (gameService *GameService) AddGame(
	ctx context.Context,
	gameToAdd *gamev4.Game,
) (uint64, error) {
	const operationPlace = "gameservice.AddGame"
	log := gameService.log.With("operationPlace", operationPlace)
	_, err := gameService.gameProvider.GetGameByTitleAndReleaseYear(ctx, gameToAdd.GetTitle(), gameToAdd.GetReleaseYear().Year)
	if err != nil {
		if errors.Is(err, outerror.ErrGameAlreadyExist) {
			log.Warn(fmt.Sprintf("game with title=%q and release year=%d already exist", gameToAdd.GetTitle(), gameToAdd.GetReleaseYear().Year))
			return 0, outerror.ErrGameAlreadyExist
		}
	}
	return 0, nil
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
