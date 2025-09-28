//go:build unit
// +build unit

package gameservice

import (
	"context"
	"errors"
	"testing"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetGame(t *testing.T) {
	t.Run("Игра не найдена", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameID := uint64(1)
		expectedError := outerror.ErrGameNotFound
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, expectedError).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, game)
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameID := uint64(1)
		expectedError := errors.New("some error")
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, expectedError).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, game)
	})
	t.Run("Успешное получение игры", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameID := uint64(1)
		expectedGame := random.NewRandomGame()
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(expectedGame, nil).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.NoError(t, err)
		require.Equal(t, expectedGame, game)
	})
}
