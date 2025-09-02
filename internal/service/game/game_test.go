package gameservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/date"
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

type mockS3Storager struct {
	mock.Mock
}

func (m *mockS3Storager) SaveObject(ctx context.Context, name string, data io.Reader) (string, error) {
	args := m.Called(ctx, name, data)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockS3Storager) GetObject(ctx context.Context, name string) (io.Reader, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(io.Reader), args.Error(1)
}

func (m *mockS3Storager) DeleteObject(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

type mockEmailAlerter struct {
	mock.Mock
}

func (m *mockEmailAlerter) SendMessage(subject string, body string) error {
	args := m.Called(subject, body)

	return args.Error(0)
}

type mockTagRepository struct {
	mock.Mock
}

func (m *mockTagRepository) GetTags(ctx context.Context, tags []string) ([]model.Tag, error) {
	args := m.Called(ctx, tags)
	return args.Get(0).([]model.Tag), args.Error(1)
}

type mockGenreRepository struct {
	mock.Mock
}

func (m *mockGenreRepository) GetGenres(ctx context.Context, genres []string) ([]model.Genre, error) {
	args := m.Called(ctx, genres)
	return args.Get(0).([]model.Genre), args.Error(1)
}

func TestAddGame(t *testing.T) {
	gameMockRepo := new(mockGameReposiroy)
	tagMockRepo := new(mockTagRepository)
	genreMockRepo := new(mockGenreRepository)
	s3Mock := new(mockS3Storager)
	mailerMock := new(mockEmailAlerter)
	gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
	t.Run("Нельзя добавить игру, так как она уже есть в БД", func(t *testing.T) {
		expectedError := outerror.ErrGameAlreadyExist
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		game := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
		}
		gameMockRepo.On("GetGameByTitleAndReleaseYear", mock.Anything, gameToAdd.Title, gameToAdd.GetReleaseDate().Year).Return(game, nil).Once()

		savedGame, err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, savedGame)
	})
	t.Run("Не удалось сохранить игру", func(t *testing.T) {
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		game := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
		}
		expectedErr := errors.New("some error")
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(nil, expectedErr).Once()
		savedGame, err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, savedGame)
	})
	t.Run("Игра сохранена, но не удалось сохранить обложку в S3", func(t *testing.T) {
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
		}
		game := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
		}
		expectedErr := outerror.ErrCannotSaveGameImage
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return("", expectedErr).Once()
		mailerMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil).Once()
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(game, nil).Once()
		savedGame, err := gameService.AddGame(context.Background(), gameToAdd)
		require.Equal(t, game, savedGame)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("Сохранение игры без ошибок", func(t *testing.T) {
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
		}
		game := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			ImageURL:    "qwe",
		}
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return("qwe", nil).Once()
		mailerMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil).Once()
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(game, nil).Once()
		savedGame, err := gameService.AddGame(context.Background(), gameToAdd)
		require.Equal(t, game, savedGame)
		require.NoError(t, err)
	})
}

func TestGetGame(t *testing.T) {
	gameMockRepo := new(mockGameReposiroy)
	tagMockRepo := new(mockTagRepository)
	genreMockRepo := new(mockGenreRepository)
	s3Mock := new(mockS3Storager)
	mailerMock := new(mockEmailAlerter)
	gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
	t.Run("Игра не найдена", func(t *testing.T) {
		gameID := uint64(1)
		expectedError := outerror.ErrGameNotFound
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, expectedError).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, game)
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		gameID := uint64(1)
		expectedError := errors.New("some error")
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, expectedError).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, game)
	})
	t.Run("Успешное получение игры", func(t *testing.T) {
		gameID := uint64(1)
		expectedGame := &model.Game{
			Title:       "Dark Souls 3",
			Genres:      []model.Genre{{GenreID: 1, GenreName: "Hard"}},
			Description: "qwe",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			Tags:        []model.Tag{{TagID: 1, TagName: "Aboba"}},
			GameID:      2,
		}
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(expectedGame, nil).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.NoError(t, err)
		require.Equal(t, expectedGame, game)
	})
}

func TestDeleteGame(t *testing.T) {
	gameMockRepo := new(mockGameReposiroy)
	tagMockRepo := new(mockTagRepository)
	genreMockRepo := new(mockGenreRepository)
	s3Mock := new(mockS3Storager)
	mailerMock := new(mockEmailAlerter)
	gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
	t.Run("Успешное удаление игры", func(t *testing.T) {
		gameID := uint64(4)
		deletedGame := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			ImageURL:    "qwe",
		}
		gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseDate.Year()))
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(deletedGame, nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(nil).Once()

		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.NoError(t, err)
		require.Equal(t, deletedGame, game)
	})
	t.Run("Нет игры для удаления", func(t *testing.T) {
		gameID := uint64(4)
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(nil, outerror.ErrGameNotFound).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.ErrorIs(t, err, outerror.ErrGameNotFound)
		require.Nil(t, game)
	})
	t.Run("Неожиданная ошибка при удалении игры", func(t *testing.T) {
		gameID := uint64(4)
		someErr := errors.New("some err")
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(nil, someErr).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.ErrorIs(t, err, someErr)
		require.Nil(t, game)
	})
	t.Run("У игры нет обложки в S3", func(t *testing.T) {
		gameID := uint64(4)
		deletedGame := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			ImageURL:    "qwe",
		}
		gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseDate.Year()))
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(deletedGame, nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(outerror.ErrImageNotFoundS3).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.Equal(t, deletedGame, game)
		require.ErrorIs(t, err, outerror.ErrImageNotFoundS3)
	})
	t.Run("Не удалось удалить обложку из S3", func(t *testing.T) {
		gameID := uint64(4)
		deletedGame := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			ImageURL:    "qwe",
		}
		someErr := errors.New("some error")
		gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseDate.Year()))
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(deletedGame, nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(someErr).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.Equal(t, deletedGame, game)
		require.ErrorIs(t, err, someErr)
	})
}
