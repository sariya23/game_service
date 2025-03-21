package gameservice

import (
	"context"
	"testing"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/date"
)

type mockGameProvider struct {
	mock.Mock
}

func (m *mockGameProvider) GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (gamev4.Game, error) {
	args := m.Called(ctx, title, releaseYear)
	return args.Get(0).(gamev4.Game), args.Error(1)
}

type mockKafkaProducer struct {
	mock.Mock
}

func (m *mockKafkaProducer) SendMessage(message string) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestAddGame(t *testing.T) {
	gameProviderMock := new(mockGameProvider)
	kafkaMock := new(mockKafkaProducer)
	gameService := NewGameService(mockslog.NewDiscardLogger(), kafkaMock, gameProviderMock)
	t.Run("Нельзя добавить игру, так как она уже есть в БД", func(t *testing.T) {
		expectedError := outerror.ErrGameAlreadyExist
		game := gamev4.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseYear: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		gameProviderMock.On("GetGameByTitleAndReleaseYear", mock.Anything, game.Title, game.ReleaseYear.Year).Return(gamev4.Game{}, expectedError)

		gameID, err := gameService.AddGame(context.Background(), &game)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, gameID, uint64(0))
	})
}
