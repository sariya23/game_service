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

func (m *mockGameServicer) AddGame(ctx context.Context, game *gamev4.GameRequest) (*gamev4.DomainGame, error) {
	args := m.Called(ctx, game)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gamev4.DomainGame), args.Error(1)
}

func (m *mockGameServicer) GetGame(ctx context.Context, gameID uint64) (*gamev4.DomainGame, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gamev4.DomainGame), args.Error(1)
}

func (m *mockGameServicer) GetTopGames(ctx context.Context, gameFilters model.GameFilters, limit uint32) ([]*gamev4.DomainGame, error) {
	args := m.Called(ctx, gameFilters, limit)
	return args.Get(0).([]*gamev4.DomainGame), args.Error(1)

}

func (m *mockGameServicer) DeleteGame(ctx context.Context, gameID uint64) (*gamev4.DomainGame, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gamev4.DomainGame), args.Error(1)
}

func TestAddGameHandler(t *testing.T) {
	mockGameService := new(mockGameServicer)
	srv := serverAPI{gameServicer: mockGameService}
	t.Run("Успешное добавление игры", func(t *testing.T) {
		game := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		expectedGame := &gamev4.DomainGame{
			Title:         "Dark Souls 3",
			Genres:        []string{"Action RPG", "Dark Fantasy"},
			Description:   "test",
			ReleaseDate:   &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImageUrl: "http://",
			Tags:          []string{"Hard"},
		}
		req := gamev4.AddGameRequest{Game: game}
		mockGameService.On("AddGame", mock.Anything, game).Return(expectedGame, nil).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		require.NoError(t, err)
		require.Equal(t, expectedGame, resp.GetGame())
	})
	t.Run("Не указано поле Title", func(t *testing.T) {
		game := &gamev4.GameRequest{
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		req := &gamev4.AddGameRequest{Game: game}
		resp, err := srv.AddGame(context.Background(), req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.TitleRequiredMessage, s.Message())
		require.Nil(t, resp.GetGame())
	})
	t.Run("Не указано поле Description", func(t *testing.T) {
		game := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		req := gamev4.AddGameRequest{Game: game}
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.DescriptionRequiredMessage, s.Message())
		require.Nil(t, resp.GetGame())
	})
	t.Run("Не указано поле Release Year", func(t *testing.T) {
		game := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		req := gamev4.AddGameRequest{Game: game}
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.ReleaseYearRequiredMessage, s.Message())
		require.Nil(t, resp.GetGame())
	})
	t.Run("Игра уже существует", func(t *testing.T) {
		game := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		req := gamev4.AddGameRequest{Game: game}
		mockGameService.On("AddGame", mock.Anything, game).Return(nil, outerror.ErrGameAlreadyExist).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.AlreadyExists, s.Code())
		require.Equal(t, outerror.GameAlreadyExistMessage, s.Message())
		require.Nil(t, resp.GetGame())
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		game := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		req := gamev4.AddGameRequest{Game: game}
		mockGameService.On("AddGame", mock.Anything, game).Return(nil, errors.New("some error")).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.Internal, s.Code())
		require.Nil(t, resp.GetGame())
	})
	t.Run("Игра сохранена, но без обложки", func(t *testing.T) {
		game := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Genres:      []string{"Action RPG", "Dark Fantasy"},
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"Hard"},
		}
		expectedGame := &gamev4.DomainGame{
			Title:         "Dark Souls 3",
			Genres:        []string{"Action RPG", "Dark Fantasy"},
			Description:   "test",
			ReleaseDate:   &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImageUrl: "http://",
			Tags:          []string{"Hard"},
		}
		req := gamev4.AddGameRequest{Game: game}
		mockGameService.On("AddGame", mock.Anything, game).Return(expectedGame, outerror.ErrCannotSaveGameImage)
		resp, err := srv.AddGame(context.Background(), &req)
		require.Equal(t, expectedGame, resp.GetGame())
		require.NoError(t, err)
	})
}

func TestGetGameHandler(t *testing.T) {
	mockGameService := new(mockGameServicer)
	srv := serverAPI{gameServicer: mockGameService}
	t.Run("Игра не найдена", func(t *testing.T) {
		gameID := uint64(2)
		req := gamev4.GetGameRequest{GameId: gameID}

		mockGameService.On("GetGame", mock.Anything, gameID).Return(nil, outerror.ErrGameNotFound).Once()
		resp, err := srv.GetGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.NotFound, s.Code())
		require.Equal(t, outerror.GameNotFoundMessage, s.Message())
		require.Nil(t, resp.GetGame())
	})
	t.Run("Успешное получение игры", func(t *testing.T) {
		gameID := uint64(2)
		expectedGame := &gamev4.DomainGame{
			Title:         "Dark Souls 3",
			Genres:        []string{"Action RPG", "Dark Fantasy"},
			Description:   "test",
			ReleaseDate:   &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImageUrl: "https://",
			Tags:          []string{"Hard"},
			Rating:        5.0,
		}
		req := gamev4.GetGameRequest{GameId: gameID}

		mockGameService.On("GetGame", mock.Anything, gameID).Return(expectedGame, nil).Once()
		resp, err := srv.GetGame(context.Background(), &req)
		require.NoError(t, err)
		require.Equal(t, expectedGame, resp.GetGame())
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		gameID := uint64(2)
		req := gamev4.GetGameRequest{GameId: gameID}

		mockGameService.On("GetGame", mock.Anything, gameID).Return(nil, errors.New("some error")).Once()
		resp, err := srv.GetGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.Internal, s.Code())
		require.Equal(t, outerror.InternalMessage, s.Message())
		require.Nil(t, resp.GetGame())
	})
}
