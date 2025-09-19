package grpchandlers

import (
	"context"
	"errors"
	"testing"

	"github.com/sariya23/game_service/internal/converters"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockGameServicer struct {
	mock.Mock
}

func (m *mockGameServicer) AddGame(ctx context.Context, game *gamev4.GameRequest) (uint64, error) {
	args := m.Called(ctx, game)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *mockGameServicer) GetGame(ctx context.Context, gameID uint64) (*model.Game, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}

func (m *mockGameServicer) GetTopGames(ctx context.Context, gameFilters dto.GameFilters, limit uint32) ([]model.ShortGame, error) {
	args := m.Called(ctx, gameFilters, limit)
	return args.Get(0).([]model.ShortGame), args.Error(1)

}

func (m *mockGameServicer) DeleteGame(ctx context.Context, gameID uint64) (uint64, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(uint64), args.Error(1)
}

func TestAddGameHandler(t *testing.T) {
	t.Run("Успешное добавление игры", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		game := random.RandomAddGameRequest()
		gameID := uint64(222)
		req := gamev4.AddGameRequest{Game: game}
		mockGameService.On("AddGame", mock.Anything, game).Return(gameID, nil).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		require.NoError(t, err)
		require.Equal(t, gameID, resp.GetGameId())
	})
	t.Run("Не указано поле Title", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		game := random.RandomAddGameRequest()
		game.Title = ""
		req := &gamev4.AddGameRequest{Game: game}
		resp, err := srv.AddGame(context.Background(), req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.TitleRequiredMessage, s.Message())
		require.Equal(t, uint64(0), resp.GetGameId())
	})
	t.Run("Не указано поле Description", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		game := random.RandomAddGameRequest()
		game.Description = ""
		req := gamev4.AddGameRequest{Game: game}
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.DescriptionRequiredMessage, s.Message())
		require.Equal(t, uint64(0), resp.GetGameId())
	})
	t.Run("Не указано поле Release Year", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		game := random.RandomAddGameRequest()
		game.ReleaseDate = nil
		req := gamev4.AddGameRequest{Game: game}
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.ReleaseYearRequiredMessage, s.Message())
		require.Equal(t, uint64(0), resp.GetGameId())
	})
	t.Run("Игра уже существует", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		game := random.RandomAddGameRequest()
		req := gamev4.AddGameRequest{Game: game}
		mockGameService.On("AddGame", mock.Anything, game).Return(uint64(0), outerror.ErrGameAlreadyExist).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.AlreadyExists, s.Code())
		require.Equal(t, outerror.GameAlreadyExistMessage, s.Message())
		require.Equal(t, uint64(0), resp.GetGameId())
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		game := random.RandomAddGameRequest()
		req := gamev4.AddGameRequest{Game: game}
		mockGameService.On("AddGame", mock.Anything, game).Return(uint64(0), errors.New("some error")).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.Internal, s.Code())
		require.Equal(t, uint64(0), resp.GetGameId())
	})
	t.Run("Игра сохранена, но без обложки", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		game := random.RandomAddGameRequest()
		gameID := uint64(228)
		req := gamev4.AddGameRequest{Game: game}
		mockGameService.On("AddGame", mock.Anything, game).Return(gameID, outerror.ErrCannotSaveGameImage)
		resp, err := srv.AddGame(context.Background(), &req)
		require.NoError(t, err)
		require.Equal(t, gameID, resp.GetGameId())
	})
	t.Run("Нельзя создать игру с несуществующим жанром", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		game := random.RandomAddGameRequest()
		req := gamev4.AddGameRequest{Game: game}
		mockGameService.On("AddGame", mock.Anything, game).Return(uint64(0), outerror.ErrGenreNotFound).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.GenreNotFoundMessage, s.Message())
		require.Equal(t, uint64(0), resp.GetGameId())
	})
	t.Run("Нельзя создать игру с несуществующим тэгом", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		game := random.RandomAddGameRequest()
		req := gamev4.AddGameRequest{Game: game}
		mockGameService.On("AddGame", mock.Anything, game).Return(uint64(0), outerror.ErrTagNotFound).Once()
		resp, err := srv.AddGame(context.Background(), &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.TagNotFoundMessage, s.Message())
		require.Equal(t, uint64(0), resp.GetGameId())
	})
}

func TestGetGameHandler(t *testing.T) {
	t.Run("Игра не найдена", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
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
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		gameID := uint64(2)
		expectedGame := random.NewRandomGame()
		req := gamev4.GetGameRequest{GameId: gameID}
		mockGameService.On("GetGame", mock.Anything, gameID).Return(expectedGame, nil).Once()
		resp, err := srv.GetGame(context.Background(), &req)
		require.NoError(t, err)
		require.Equal(t, *converters.ToProtoGame(*expectedGame), *resp.GetGame())
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
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

func TestDeleteGame(t *testing.T) {
	t.Run("Игра не найдена", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		ctx := context.Background()
		gameID := uint64(2)
		req := gamev4.DeleteGameRequest{GameId: gameID}
		mockGameService.On("DeleteGame", mock.Anything, gameID).Return(uint64(0), outerror.ErrGameNotFound)
		resp, err := srv.DeleteGame(ctx, &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.NotFound, s.Code())
		require.Equal(t, outerror.GameNotFoundMessage, s.Message())
		require.Zero(t, resp.GetGameId())
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		ctx := context.Background()
		gameID := uint64(2)
		req := gamev4.DeleteGameRequest{GameId: gameID}
		mockGameService.On("DeleteGame", mock.Anything, gameID).Return(uint64(0), errors.New("err"))
		resp, err := srv.DeleteGame(ctx, &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.Internal, s.Code())
		require.Equal(t, outerror.InternalMessage, s.Message())
		require.Zero(t, resp.GetGameId())
	})
	t.Run("Успешное удаление игры", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		ctx := context.Background()
		gameID := uint64(2)
		req := gamev4.DeleteGameRequest{GameId: gameID}
		mockGameService.On("DeleteGame", mock.Anything, gameID).Return(gameID, nil)
		resp, err := srv.DeleteGame(ctx, &req)
		require.NoError(t, err)
		require.Equal(t, gameID, resp.GetGameId())
	})
}

func TestGetTopGames(t *testing.T) {
	t.Run("По фильтрам не нашлось игр", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		ctx := context.Background()
		req := gamev4.GetTopGamesRequest{Limit: 10, Year: 2020}
		filters := dto.GameFilters{ReleaseYear: req.GetYear(), Tags: req.GetTags(), Genres: req.GetGenres()}
		mockGameService.On("GetTopGames", mock.Anything, filters, req.GetLimit()).Return([]model.ShortGame{}, nil)
		resp, err := srv.GetTopGames(ctx, &req)
		require.NoError(t, err)
		require.Empty(t, resp.GetGames())
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		ctx := context.Background()
		req := gamev4.GetTopGamesRequest{Limit: 10, Year: 2020}
		filters := dto.GameFilters{ReleaseYear: req.GetYear(), Tags: req.GetTags(), Genres: req.GetGenres()}
		mockGameService.On("GetTopGames", mock.Anything, filters, req.GetLimit()).Return([]model.ShortGame{}, errors.New("err"))
		resp, err := srv.GetTopGames(ctx, &req)
		s, _ := status.FromError(err)
		require.Equal(t, codes.Internal, s.Code())
		require.Equal(t, outerror.InternalMessage, s.Message())
		require.Nil(t, resp.GetGames())
	})
	t.Run("Успешное получение топа", func(t *testing.T) {
		mockGameService := new(mockGameServicer)
		srv := serverAPI{gameServicer: mockGameService}
		ctx := context.Background()
		req := gamev4.GetTopGamesRequest{Limit: 10, Year: 2020}
		filters := dto.GameFilters{ReleaseYear: req.GetYear(), Tags: req.GetTags(), Genres: req.GetGenres()}
		mockGameService.On("GetTopGames", mock.Anything, filters, req.GetLimit()).Return([]model.ShortGame{*random.NewRandomShortGame()}, nil)
		resp, err := srv.GetTopGames(ctx, &req)
		require.NoError(t, err)
		require.NotEmpty(t, resp.GetGames())
	})
}
