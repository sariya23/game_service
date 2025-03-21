package gameservice

import (
	"context"
	"log/slog"

	"github.com/sariya23/game_service/internal/model"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type KafkaProducer interface {
	SendMessage(message string) error
}

type GameProvider interface {
	GetGame(ctx context.Context, gameID uint64) (game gamev4.GameWithRating, err error)
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
	game *gamev4.Game,
) (uint64, error) {
	panic("impl me")
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
