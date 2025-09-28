//go:build unit
// +build unit

package gameservice

import (
	"context"
	"errors"
	"testing"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteGame(t *testing.T) {
	t.Run("Успешное удаление игры", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameID := uint64(4)
		deletedGame := random.NewRandomGame()
		gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseDate.Year()))
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(dto.DeletedGameFromGame(deletedGame), nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(nil).Once()

		deletedGameID, err := gameService.DeleteGame(context.Background(), gameID)
		require.NoError(t, err)
		require.Equal(t, deletedGame.GameID, deletedGameID)
	})
	t.Run("Нет игры для удаления", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameID := uint64(4)
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(nil, outerror.ErrGameNotFound).Once()
		deletedGameID, err := gameService.DeleteGame(context.Background(), gameID)
		require.ErrorIs(t, err, outerror.ErrGameNotFound)
		require.Zero(t, deletedGameID)
	})
	t.Run("Неожиданная ошибка при удалении игры", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameID := uint64(4)
		someErr := errors.New("some err")
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(nil, someErr).Once()
		deletedGameID, err := gameService.DeleteGame(context.Background(), gameID)
		require.ErrorIs(t, err, someErr)
		require.Zero(t, deletedGameID)
	})
	t.Run("У игры нет обложки в S3", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameID := uint64(4)
		deletedGame := random.NewRandomGame()
		deletedGame.ImageURL = ""
		gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseDate.Year()))
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(dto.DeletedGameFromGame(deletedGame), nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(nil).Once()
		deletedGameID, err := gameService.DeleteGame(context.Background(), gameID)
		require.Equal(t, deletedGame.GameID, deletedGameID)
		require.NoError(t, err)
	})
	t.Run("Не удалось удалить обложку из S3", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameID := uint64(4)
		deletedGame := random.NewRandomGame()
		someErr := errors.New("some error")
		gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseDate.Year()))
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(dto.DeletedGameFromGame(deletedGame), nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(someErr).Once()
		deletedGameID, err := gameService.DeleteGame(context.Background(), gameID)
		require.Equal(t, deletedGame.GameID, deletedGameID)
		require.ErrorIs(t, err, someErr)
	})
}
