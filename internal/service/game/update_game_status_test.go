package gameservice

import (
	"context"
	"errors"
	"testing"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateGameStatus(t *testing.T) {
	t.Run("Игра не найдена, возвращаем ошибку Not Found", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)

		gameID := uint64(228)
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, outerror.ErrGameNotFound).Once()
		err := gameService.UpdateGameStatus(context.Background(), gameID, gamev4.GameStatusType_PENDING)
		require.ErrorIs(t, err, outerror.ErrGameNotFound)
	})
	t.Run("Internal ошибка при получении игры", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameID := uint64(228)
		someErr := errors.New("err")
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, someErr).Once()
		err := gameService.UpdateGameStatus(context.Background(), gameID, gamev4.GameStatusType_DRAFT)
		require.ErrorIs(t, err, someErr)
	})
}
