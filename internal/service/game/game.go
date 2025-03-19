package gameservice

import (
	"context"
	"log/slog"

	"github.com/sariya23/game_service/internal/model"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
)

type GameService struct {
	log *slog.Logger
}

func NewGameService(log *slog.Logger) *GameService {
	return &GameService{
		log: log,
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
	gameTitle string,
) (*gamev4.Game, error) {
	panic("impl me")
}

func (gameService *GameService) GetTopGames(
	ctx context.Context,
	gameFilters model.GameFilters,
	limit uint32,
) ([]gamev4.Game, error) {
	panic("impl me")
}
