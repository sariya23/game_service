package gameservice

import (
	"context"
	"log/slog"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/domain"
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
	game domain.Game,
) (int64, error) {
	panic("impl me")
}

func (gameService *GameService) GetGame(
	ctx context.Context,
	gameTitle string,
) (domain.Game, error) {
	panic("impl me")
}

func (gameService *GameService) GetTopGames(
	ctx context.Context,
	gameFilters model.GameFilters, limit int32,
) ([]domain.Game, error) {
	panic("impl me")
}
