package grpchandlers

import (
	"context"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/domain"
	"github.com/stretchr/testify/mock"
)

type mockAuther struct {
	mock.Mock
}

func (m *mockAuther) AddGame(ctx context.Context, game domain.Game) (int64, error) {
	args := m.Called(ctx, game)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockAuther) GetGame(ctx context.Context, gameTitle string) (domain.Game, error) {
	args := m.Called(ctx, gameTitle)
	return args.Get(0).(domain.Game), args.Error(1)
}

func (m *mockAuther) GetTopGames(ctx context.Context, gameFilters model.GameFilters, limit int32) ([]domain.Game, error) {
	args := m.Called(ctx, gameFilters, limit)
	return args.Get(0).([]domain.Game), args.Error(1)

}
