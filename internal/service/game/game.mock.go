package gameservice

import (
	"context"

	"github.com/sariya23/game_service/internal/model"
	"github.com/stretchr/testify/mock"
)

type mockGameReposiroy struct {
	mock.Mock
}

func (m *mockGameReposiroy) GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (*model.Game, error) {
	args := m.Called(ctx, title, releaseYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}

func (m *mockGameReposiroy) GetGameByID(ctx context.Context, gameID uint64) (*model.Game, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}

func (m *mockGameReposiroy) GetTopGames(ctx context.Context, releaseYear int32, tags []string, genres []string, limit uint32) ([]model.Game, error) {
	args := m.Called(ctx, releaseYear, tags, genres, limit)
	return args.Get(0).([]model.Game), args.Error(1)
}

func (m *mockGameReposiroy) SaveGame(ctx context.Context, game model.Game) (*model.Game, error) {
	args := m.Called(ctx, game)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}

func (m *mockGameReposiroy) DaleteGame(ctx context.Context, gameID uint64) (*model.Game, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}
