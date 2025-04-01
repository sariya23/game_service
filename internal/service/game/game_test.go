package gameservice

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/game_service/internal/storage/s3"
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

type mockS3Storager struct {
	mock.Mock
}

func (m *mockS3Storager) Save(ctx context.Context, data io.Reader, key string) (string, error) {
	args := m.Called(ctx, data, key)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockS3Storager) Get(ctx context.Context, bucket, key string) io.Reader {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(io.Reader)
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

func TestAddGame(t *testing.T) {
	gameProviderMock := new(mockGameProvider)
	gameSaverMock := new(mockGameSaver)
	s3Mock := new(mockS3Storager)
	mailerMock := new(mockEmailAlerter)
	gameService := NewGameService(mockslog.NewDiscardLogger(), gameProviderMock, s3Mock, gameSaverMock, mailerMock)
	t.Run("Нельзя добавить игру, так как она уже есть в БД", func(t *testing.T) {
		expectedError := outerror.ErrGameAlreadyExist
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		domainGame := &gamev4.DomainGame{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		gameProviderMock.On("GetGameByTitleAndReleaseYear", mock.Anything, gameToAdd.Title, gameToAdd.ReleaseYear.Year).Return(domainGame, nil).Once()

		savedGame, err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, savedGame)
	})
	t.Run("Не удалось сохранить игру", func(t *testing.T) {
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		domainGame := &gamev4.DomainGame{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		expectedErr := errors.New("some error")
		gameProviderMock.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.Title,
			gameToAdd.ReleaseYear.Year,
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
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
		}
		domainGame := &gamev4.DomainGame{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		expectedErr := outerror.ErrCannotSaveGameImage
		gameProviderMock.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.Title,
			gameToAdd.ReleaseYear.Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		s3Mock.On(
			"Save",
			mock.Anything,
			bytes.NewReader(gameToAdd.GetCoverImage()),
			s3.CreateGameKey(gameToAdd.Title, int(gameToAdd.GetReleaseYear().Year)),
		).Return("", expectedErr).Once()
		gameSaverMock.On("SaveGame", mock.Anything, domainGame).Return(domainGame, nil).Once()
		savedGame, err := gameService.AddGame(context.Background(), gameToAdd)
		require.Equal(t, domainGame, savedGame)
		require.ErrorIs(t, err, expectedErr)
	})
}
