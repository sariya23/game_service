package grpchandlers

import (
	"context"
	"testing"
	"time"

	"github.com/sariya23/game_service/internal/model"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mockGameServicer struct {
	mock.Mock
}

func (m *mockGameServicer) AddGame(ctx context.Context, game *gamev4.Game) (uint64, error) {
	args := m.Called(ctx, game)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *mockGameServicer) GetGame(ctx context.Context, gameTitle string) (*gamev4.Game, error) {
	args := m.Called(ctx, gameTitle)
	return args.Get(0).(*gamev4.Game), args.Error(1)
}

func (m *mockGameServicer) GetTopGames(ctx context.Context, gameFilters model.GameFilters, limit uint32) ([]gamev4.Game, error) {
	args := m.Called(ctx, gameFilters, limit)
	return args.Get(0).([]gamev4.Game), args.Error(1)

}

func TestAddGame(t *testing.T) {
	mockGameService := new(mockGameServicer)
	srv := serverAPI{gameServicer: mockGameService}
	t.Run("Успешное добавление игры", func(t *testing.T) {
		game := gamev4.Game{
			Title:       "Dark Souls 3",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			Description: "test",
			ReleaseYear: timestamppb.New(time.Date(2016, time.March, 16, 0, 0, 0, 0, time.UTC)),
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		expcetedGameId := uint64(1)
		req := gamev4.AddGameRequest{Game: &game}
		mockGameService.On("AddGame", game).Return(expcetedGameId, nil)
		resp, err := srv.AddGame(context.Background(), &req)
		require.NoError(t, err)
	})
}
