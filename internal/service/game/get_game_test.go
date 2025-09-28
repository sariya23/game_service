//go:build unit
// +build unit

package gameservice

import (
	"context"
	"errors"
	"testing"

	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetGame(t *testing.T) {
	t.Run("Игра не найдена", func(t *testing.T) {
		suite := NewSuite()
		gameID := uint64(1)
		expectedError := outerror.ErrGameNotFound
		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, expectedError).Once()
		game, err := suite.gameService.GetGame(context.Background(), gameID)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, game)
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		suite := NewSuite()
		gameID := uint64(1)
		expectedError := errors.New("some error")
		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, expectedError).Once()
		game, err := suite.gameService.GetGame(context.Background(), gameID)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, game)
	})
	t.Run("Успешное получение игры", func(t *testing.T) {
		suite := NewSuite()
		gameID := uint64(1)
		expectedGame := random.NewRandomGame()
		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(expectedGame, nil).Once()
		game, err := suite.gameService.GetGame(context.Background(), gameID)
		require.NoError(t, err)
		require.Equal(t, expectedGame, game)
	})
}
