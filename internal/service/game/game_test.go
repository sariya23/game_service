package gameservice

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/model/domain"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/date"
)

type mockGameProvider struct {
	mock.Mock
}

func (m *mockGameProvider) GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (domain.Game, error) {
	args := m.Called(ctx, title, releaseYear)
	return args.Get(0).(domain.Game), args.Error(1)
}

type mockKafkaProducer struct {
	mock.Mock
}

func (m *mockKafkaProducer) SendMessage(message string) error {
	args := m.Called(message)
	return args.Error(0)
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

func (m *mockGameSaver) SaveGame(ctx context.Context, game domain.Game) (uint64, error) {
	args := m.Called(ctx, game)

	return args.Get(0).(uint64), args.Error(1)
}

func TestAddGame(t *testing.T) {
	gameProviderMock := new(mockGameProvider)
	gameSaverMock := new(mockGameSaver)
	kafkaMock := new(mockKafkaProducer)
	s3Mock := new(mockS3Storager)
	gameService := NewGameService(mockslog.NewDiscardLogger(), kafkaMock, gameProviderMock, s3Mock, gameSaverMock)
	t.Run("Нельзя добавить игру, так как она уже есть в БД", func(t *testing.T) {
		expectedError := outerror.ErrGameAlreadyExist
		gameToAdd := gamev4.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		domainGame := domain.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		gameProviderMock.On("GetGameByTitleAndReleaseYear", mock.Anything, gameToAdd.Title, gameToAdd.ReleaseYear.Year).Return(domainGame, nil).Once()

		gameID, err := gameService.AddGame(context.Background(), &gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, gameID, uint64(0))
	})
	t.Run("Не удалось сохранить игру", func(t *testing.T) {
		gameToAdd := gamev4.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		domainGame := domain.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		expectedErr := errors.New("some error")
		gameProviderMock.On("GetGameByTitleAndReleaseYear", mock.Anything, gameToAdd.Title, gameToAdd.ReleaseYear.Year).Return(domain.Game{}, outerror.ErrGameNotFound).Once()
		gameSaverMock.On("SaveGame", mock.Anything, domainGame).Return(uint64(0), expectedErr)
		gameID, err := gameService.AddGame(context.Background(), &gameToAdd)
		require.ErrorIs(t, err, expectedErr)
		require.Equal(t, uint64(0), gameID)
	})
}
