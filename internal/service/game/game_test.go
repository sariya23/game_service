package gameservice

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/model/domain"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/game_service/internal/storage/postgresql"
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
	return args.String(0), args.Error(1)
}

type mockGameSaver struct {
	mock.Mock
}

func (m *mockGameSaver) SaveGame(ctx context.Context, game domain.Game) (*postgresql.GameTransaction, error) {
	args := m.Called(ctx, game)

	if args.Get(0) != nil {
		return args.Get(0).(*postgresql.GameTransaction), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *mockS3Storager) Get(ctx context.Context, bucket, key string) io.Reader {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(io.Reader)
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
	t.Run("Не удалось начать транзакцию", func(t *testing.T) {
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
		gameProviderMock.On("GetGameByTitleAndReleaseYear", mock.Anything, gameToAdd.Title, gameToAdd.ReleaseYear.Year).Return(domainGame, outerror.ErrGameAlreadyExist).Once()
		gameSaverMock.On("SaveGame", mock.Anything, domainGame).Return(nil, errors.New("some error"))
		gameID, err := gameService.AddGame(context.Background(), &gameToAdd)
		require.ErrorIs(t, err, outerror.ErrCannotStartGameTransaction)
		require.Equal(t, uint64(0), gameID)
	})
}
