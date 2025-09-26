package gameservice

import (
	"context"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
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

func (m *mockGameReposiroy) GetTopGames(ctx context.Context, filters dto.GameFilters, limit uint32) ([]model.ShortGame, error) {
	args := m.Called(ctx, filters, limit)
	return args.Get(0).([]model.ShortGame), args.Error(1)
}

func (m *mockGameReposiroy) SaveGame(ctx context.Context, game model.Game) (uint64, error) {
	args := m.Called(ctx, game)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *mockGameReposiroy) DaleteGame(ctx context.Context, gameID uint64) (*dto.DeletedGame, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DeletedGame), args.Error(1)
}

func (m *mockGameReposiroy) UpdateGameStatus(ctx context.Context, gameID uint64, newStatus gamev4.GameStatusType) error {
	args := m.Called(ctx, gameID, newStatus)
	return args.Error(0)
}
