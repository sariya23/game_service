//go:build unit
// +build unit

package gameservice

import (
	"context"
	"errors"
	"testing"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateGameStatus(t *testing.T) {
	t.Run("Игра не найдена, возвращаем ошибку Not Found", func(t *testing.T) {
		suite := NewSuite()

		gameID := uint64(228)
		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, outerror.ErrGameNotFound).Once()
		err := suite.gameService.UpdateGameStatus(context.Background(), gameID, gamev4.GameStatusType_PENDING)
		require.ErrorIs(t, err, outerror.ErrGameNotFound)
	})
	t.Run("Internal ошибка при получении игры", func(t *testing.T) {
		suite := NewSuite()

		gameID := uint64(228)
		someErr := errors.New("err")
		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, someErr).Once()
		err := suite.gameService.UpdateGameStatus(context.Background(), gameID, gamev4.GameStatusType_DRAFT)
		require.ErrorIs(t, err, someErr)
	})
	t.Run("Нельзя сменить статус из DRAFT в PUBLISH", func(t *testing.T) {
		suite := NewSuite()

		gameID := uint64(228)
		newStatus := gamev4.GameStatusType_PUBLISH
		expectedGame := model.Game{GameStatus: int(gamev4.GameStatusType_DRAFT)}

		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(&expectedGame, nil).Once()
		err := suite.gameService.UpdateGameStatus(context.Background(), gameID, newStatus)
		require.ErrorIs(t, err, outerror.ErrInvalidNewGameStatus)
	})
	t.Run("Нельзя сменить статус из PUBLISH в PENDING", func(t *testing.T) {
		suite := NewSuite()

		gameID := uint64(228)
		newStatus := gamev4.GameStatusType_PENDING
		expectedGame := model.Game{GameStatus: int(gamev4.GameStatusType_PUBLISH)}

		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(&expectedGame, nil).Once()
		err := suite.gameService.UpdateGameStatus(context.Background(), gameID, newStatus)
		require.ErrorIs(t, err, outerror.ErrInvalidNewGameStatus)
	})
	t.Run("Нельзя сменить статус из PUBLISH в DRAFR", func(t *testing.T) {
		suite := NewSuite()

		gameID := uint64(228)
		newStatus := gamev4.GameStatusType_DRAFT
		expectedGame := model.Game{GameStatus: int(gamev4.GameStatusType_PUBLISH)}

		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(&expectedGame, nil).Once()
		err := suite.gameService.UpdateGameStatus(context.Background(), gameID, newStatus)
		require.ErrorIs(t, err, outerror.ErrInvalidNewGameStatus)
	})
	t.Run("Передан невалидный статус", func(t *testing.T) {
		suite := NewSuite()

		gameID := uint64(228)
		newStatus := gamev4.GameStatusType(228)
		expectedGame := model.Game{GameStatus: int(gamev4.GameStatusType_PUBLISH)}

		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(&expectedGame, nil).Once()
		err := suite.gameService.UpdateGameStatus(context.Background(), gameID, newStatus)
		require.ErrorIs(t, err, outerror.ErrUnknownGameStatus)
	})
	t.Run("Internal ошибка при обновлении статуса", func(t *testing.T) {
		suite := NewSuite()

		gameID := uint64(228)
		newStatus := gamev4.GameStatusType_PENDING
		expectedGame := model.Game{GameStatus: int(gamev4.GameStatusType_DRAFT)}
		expectedErr := errors.New("err")

		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(&expectedGame, nil).Once()
		suite.gameMockRepo.On("UpdateGameStatus", mock.Anything, gameID, newStatus).Return(expectedErr)
		err := suite.gameService.UpdateGameStatus(context.Background(), gameID, newStatus)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("Успешное обновление статуса игры", func(t *testing.T) {
		suite := NewSuite()

		gameID := uint64(228)
		newStatus := gamev4.GameStatusType_PENDING
		expectedGame := model.Game{GameStatus: int(gamev4.GameStatusType_DRAFT)}

		suite.gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(&expectedGame, nil).Once()
		suite.gameMockRepo.On("UpdateGameStatus", mock.Anything, gameID, newStatus).Return(nil)
		err := suite.gameService.UpdateGameStatus(context.Background(), gameID, newStatus)
		require.NoError(t, err)
	})
}
