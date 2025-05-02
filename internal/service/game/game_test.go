package gameservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/outerror"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/date"
)

type mockGameProvider struct {
	mock.Mock
}

func (m *mockGameProvider) GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (*gamev4.DomainGame, error) {
	args := m.Called(ctx, title, releaseYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gamev4.DomainGame), args.Error(1)
}

func (m *mockGameProvider) GetGameByID(ctx context.Context, gameID uint64) (*gamev4.DomainGame, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gamev4.DomainGame), args.Error(1)
}

func (m *mockGameProvider) GetTopGames(ctx context.Context, releaseYear int32, tags []string, genres []string, limit uint32) (games []*gamev4.DomainGame, err error) {
	args := m.Called(ctx, releaseYear, tags, genres, limit)
	return args.Get(0).([]*gamev4.DomainGame), args.Error(1)
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

type mockGameSaver struct {
	mock.Mock
}

func (m *mockGameSaver) SaveGame(ctx context.Context, game *gamev4.DomainGame) (*gamev4.DomainGame, error) {
	args := m.Called(ctx, game)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gamev4.DomainGame), args.Error(1)
}

type mockEmailAlerter struct {
	mock.Mock
}

func (m *mockEmailAlerter) SendMessage(subject string, body string) error {
	args := m.Called(subject, body)

	return args.Error(0)
}

type mockGameDeleter struct {
	mock.Mock
}

func (m *mockGameDeleter) DaleteGame(ctx context.Context, gameID uint64) (*gamev4.DomainGame, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gamev4.DomainGame), args.Error(1)
}

func TestAddGame(t *testing.T) {
	gameProviderMock := new(mockGameProvider)
	gameSaverMock := new(mockGameSaver)
	gameDeleterMock := new(mockGameDeleter)
	s3Mock := new(mockS3Storager)
	mailerMock := new(mockEmailAlerter)
	gameService := NewGameService(mockslog.NewDiscardLogger(), gameProviderMock, s3Mock, gameSaverMock, mailerMock, gameDeleterMock)
	t.Run("Нельзя добавить игру, так как она уже есть в БД", func(t *testing.T) {
		expectedError := outerror.ErrGameAlreadyExist
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		domainGame := &gamev4.DomainGame{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		gameProviderMock.On("GetGameByTitleAndReleaseYear", mock.Anything, gameToAdd.Title, gameToAdd.GetReleaseDate().Year).Return(domainGame, nil).Once()

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
		domainGame := &gamev4.DomainGame{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		expectedErr := errors.New("some error")
		gameProviderMock.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.Title,
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		gameSaverMock.On("SaveGame", mock.Anything, domainGame).Return(nil, expectedErr).Once()
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
		domainGame := &gamev4.DomainGame{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		expectedErr := outerror.ErrCannotSaveGameImage
		gameProviderMock.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.Title,
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.Title, int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return("", expectedErr).Once()
		mailerMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil).Once()
		gameSaverMock.On("SaveGame", mock.Anything, domainGame).Return(domainGame, nil).Once()
		savedGame, err := gameService.AddGame(context.Background(), gameToAdd)
		require.Equal(t, domainGame, savedGame)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("Сохранение игры без ошибок", func(t *testing.T) {
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
		}
		domainGame := &gamev4.DomainGame{
			Title:         "Dark Souls 3",
			Description:   "test",
			ReleaseDate:   &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImageUrl: "qwe",
		}
		gameProviderMock.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.Title, int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return("qwe", nil).Once()
		mailerMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil).Once()
		gameSaverMock.On("SaveGame", mock.Anything, domainGame).Return(domainGame, nil).Once()
		savedGame, err := gameService.AddGame(context.Background(), gameToAdd)
		require.Equal(t, domainGame, savedGame)
		require.NoError(t, err)
	})
}

func TestGetGame(t *testing.T) {
	gameProviderMock := new(mockGameProvider)
	gameSaverMock := new(mockGameSaver)
	s3Mock := new(mockS3Storager)
	mailerMock := new(mockEmailAlerter)
	gameDeleterMock := new(mockGameDeleter)
	gameService := NewGameService(mockslog.NewDiscardLogger(), gameProviderMock, s3Mock, gameSaverMock, mailerMock, gameDeleterMock)
	t.Run("Игра не найдена", func(t *testing.T) {
		gameID := uint64(1)
		expectedError := outerror.ErrGameNotFound
		gameProviderMock.On("GetGameByID", mock.Anything, gameID).Return(nil, expectedError).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, game)
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		gameID := uint64(1)
		expectedError := errors.New("some error")
		gameProviderMock.On("GetGameByID", mock.Anything, gameID).Return(nil, expectedError).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, game)
	})
	t.Run("Успешное получение игры", func(t *testing.T) {
		gameID := uint64(1)
		expectedGame := &gamev4.DomainGame{
			Title:       "Dark Souls 3",
			Genres:      []string{"Hard"},
			Description: "qwe",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			Tags:        []string{"asd"},
			ID:          2,
		}
		gameProviderMock.On("GetGameByID", mock.Anything, gameID).Return(expectedGame, nil).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.NoError(t, err)
		require.Equal(t, expectedGame, game)
	})
}

func TestDeleteGame(t *testing.T) {
	gameProviderMock := new(mockGameProvider)
	gameSaverMock := new(mockGameSaver)
	s3Mock := new(mockS3Storager)
	mailerMock := new(mockEmailAlerter)
	gameDeleterMock := new(mockGameDeleter)
	gameService := NewGameService(mockslog.NewDiscardLogger(), gameProviderMock, s3Mock, gameSaverMock, mailerMock, gameDeleterMock)
	t.Run("Успешное удаление игры", func(t *testing.T) {
		gameID := uint64(4)
		deletedGame := &gamev4.DomainGame{
			Title:         "Dark Souls 3",
			Description:   "test",
			ReleaseDate:   &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImageUrl: "qwe",
		}
		gameKey := minioclient.GameKey(deletedGame.GetTitle(), int(deletedGame.GetReleaseDate().Year))
		gameDeleterMock.On("DaleteGame", mock.Anything, gameID).Return(deletedGame, nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(nil).Once()

		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.NoError(t, err)
		require.Equal(t, deletedGame, game)
	})
	t.Run("Нет игры для удаления", func(t *testing.T) {
		gameID := uint64(4)
		gameDeleterMock.On("DaleteGame", mock.Anything, gameID).Return(nil, outerror.ErrGameNotFound).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.ErrorIs(t, err, outerror.ErrGameNotFound)
		require.Nil(t, game)
	})
	t.Run("Неожиданная ошибка при удалении игры", func(t *testing.T) {
		gameID := uint64(4)
		someErr := errors.New("some err")
		gameDeleterMock.On("DaleteGame", mock.Anything, gameID).Return(nil, someErr).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.ErrorIs(t, err, someErr)
		require.Nil(t, game)
	})
	t.Run("У игры нет обложки в S3", func(t *testing.T) {
		gameID := uint64(4)
		deletedGame := &gamev4.DomainGame{
			Title:         "Dark Souls 3",
			Description:   "test",
			ReleaseDate:   &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImageUrl: "qwe",
		}
		gameKey := minioclient.GameKey(deletedGame.GetTitle(), int(deletedGame.GetReleaseDate().Year))
		gameDeleterMock.On("DaleteGame", mock.Anything, gameID).Return(deletedGame, nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(outerror.ErrImageNotFoundS3).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.Equal(t, deletedGame, game)
		require.ErrorIs(t, err, outerror.ErrImageNotFoundS3)
	})
}
