package grpchandlers

import (
	"context"
	"errors"
	"testing"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockGameServicer struct {
	mock.Mock
}

func (m *mockGameServicer) AddGame(ctx context.Context, game *gamev4.Game) (uint64, error) {
	args := m.Called(ctx, game)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *mockGameServicer) GetGame(ctx context.Context, gameID uint64) (*gamev4.Game, error) {
	args := m.Called(ctx, gameID)
	return args.Get(0).(*gamev4.Game), args.Error(1)
}

func (m *mockGameServicer) GetTopGames(ctx context.Context, gameFilters model.GameFilters, limit uint32) ([]gamev4.Game, error) {
	args := m.Called(ctx, gameFilters, limit)
	return args.Get(0).([]gamev4.Game), args.Error(1)

}

func (m *mockGameServicer) DeleteGame(ctx context.Context, gameID uint64) (*gamev4.Game, error) {
	args := m.Called(ctx, gameID)
	return args.Get(0).(*gamev4.Game), args.Error(1)
}

func TestAddGame(t *testing.T) {
	mockGameService := new(mockGameServicer)
	srv := serverAPI{gameServicer: mockGameService}
	t.Run("Успешное добавление игры", func(t *testing.T) {
		game := gamev4.Game{
			Title:       "Dark Souls 3",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		expectedGameID := uint64(1)
		req := gamev4.AddGameRequest{Game: &game}
		mockGameService.On("AddGame", mock.Anything, &game).Return(expectedGameID, nil).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		require.NoError(t, err)
		require.GreaterOrEqual(t, resp.GetGameId(), uint64(0))
	})
	t.Run("Не указано поле Title", func(t *testing.T) {
		game := gamev4.Game{
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		expectedGameID := uint64(0)
		req := gamev4.AddGameRequest{Game: &game}
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.TitleRequiredMessage, s.Message())
		require.Equal(t, expectedGameID, resp.GetGameId())
	})
	t.Run("Не указано поле Description", func(t *testing.T) {
		game := gamev4.Game{
			Title:       "Dark Souls 3",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		expectedGameID := uint64(0)
		req := gamev4.AddGameRequest{Game: &game}
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.DescriptionRequiredMessage, s.Message())
		require.Equal(t, expectedGameID, resp.GetGameId())
	})
	t.Run("Не указано поле Release Year", func(t *testing.T) {
		game := gamev4.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		expectedGameID := uint64(0)
		req := gamev4.AddGameRequest{Game: &game}
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.ReleaseYearRequiredMessage, s.Message())
		require.Equal(t, expectedGameID, resp.GetGameId())
	})
	t.Run("Игра уже существует", func(t *testing.T) {
		game := gamev4.Game{
			Title:       "Dark Souls 3",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		expectedGameID := uint64(0)
		req := gamev4.AddGameRequest{Game: &game}
		mockGameService.On("AddGame", mock.Anything, &game).Return(expectedGameID, outerror.ErrGameAlreadyExist).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.AlreadyExists, s.Code())
		require.Equal(t, outerror.GameAlreadyExistMessage, s.Message())
		require.Equal(t, expectedGameID, resp.GetGameId())
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		game := gamev4.Game{
			Title:       "Dark Souls 3",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		expectedGameID := uint64(0)
		req := gamev4.AddGameRequest{Game: &game}
		mockGameService.On("AddGame", mock.Anything, &game).Return(expectedGameID, errors.New("some error")).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.Internal, s.Code())
		require.Equal(t, expectedGameID, resp.GetGameId())
	})
}

func TestGetGame(t *testing.T) {
	mockGameService := new(mockGameServicer)
	srv := serverAPI{gameServicer: mockGameService}
	t.Run("Игра не найдена", func(t *testing.T) {
		gameID := uint64(2)
		expectedGame := &gamev4.Game{}
		req := gamev4.GetGameRequest{GameId: gameID}

		mockGameService.On("GetGame", mock.Anything, gameID).Return(expectedGame, outerror.ErrGameNotFound).Once()
		resp, err := srv.GetGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.NotFound, s.Code())
		require.Equal(t, outerror.GameNotFoundMessage, s.Message())
		require.Nil(t, resp.GetGame())
	})
}
