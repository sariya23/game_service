//go:build integrations
// +build integrations

package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/lib/random"
	checkers "github.com/sariya23/game_service/tests/checkers"
	hadlerchecker "github.com/sariya23/game_service/tests/checkers/handlers"
	"github.com/sariya23/game_service/tests/suite"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/require"
)

func TestUpdateGameStatus(t *testing.T) {
	ctx, suite := suite.NewSuite(t)
	t.Run("Успешное обновление статуса игры", func(t *testing.T) {
		gameToAdd := random.WithOnlyRequireFields()
		addGameReq := gamev4.AddGameRequest{Game: gameToAdd}
		addResp, err := suite.GrpcClient.AddGame(ctx, &addGameReq)
		hadlerchecker.AssertAddGame(t, err, addResp)

		reqUpdateGameStatus := gamev4.UpdateGameStatusRequest{GameId: addResp.GameId, NewStautus: gamev4.GameStatusType_PENDING}
		_, err = suite.GrpcClient.UpdateGameStatus(ctx, &reqUpdateGameStatus)
		require.NoError(t, err)

		gameDB, err := suite.Db.GetGameByID(ctx, addResp.GameId)
		require.NoError(t, err)

		checkers.AssertAddGameRequestAndDB(ctx, t, &addGameReq, *gameDB, nil, gamev4.GameStatusType_PENDING)
	})
	t.Run("Нет игры по переданному айди", func(t *testing.T) {
		_, err := suite.GrpcClient.UpdateGameStatus(ctx, &gamev4.UpdateGameStatusRequest{GameId: uint64(gofakeit.Int64()), NewStautus: gamev4.GameStatusType_DRAFT})
		hadlerchecker.AssertUpdateGameStatusGameNotFound(t, err)
	})
	t.Run("Нельзя сменить жанр с DRAFT на PUBLISH", func(t *testing.T) {
		gameToAdd := random.WithOnlyRequireFields()
		addGameReq := gamev4.AddGameRequest{Game: gameToAdd}
		addResp, err := suite.GrpcClient.AddGame(ctx, &addGameReq)
		hadlerchecker.AssertAddGame(t, err, addResp)

		reqUpdateGameStatus := gamev4.UpdateGameStatusRequest{GameId: addResp.GameId, NewStautus: gamev4.GameStatusType_PUBLISH}
		_, err = suite.GrpcClient.UpdateGameStatus(ctx, &reqUpdateGameStatus)
		hadlerchecker.AssertUpdateGameStatusInvalidStatus(t, err)
	})
	t.Run("Нельзя сменить жанр с PUBLISH на PEDNING", func(t *testing.T) {
		gameToAdd := random.WithOnlyRequireFields()
		addGameReq := gamev4.AddGameRequest{Game: gameToAdd}
		addResp, err := suite.GrpcClient.AddGame(ctx, &addGameReq)
		hadlerchecker.AssertAddGame(t, err, addResp)

		reqUpdateGameStatus := gamev4.UpdateGameStatusRequest{GameId: addResp.GameId, NewStautus: gamev4.GameStatusType_PENDING}
		_, err = suite.GrpcClient.UpdateGameStatus(ctx, &reqUpdateGameStatus)
		require.NoError(t, err)

		reqUpdateGameStatus = gamev4.UpdateGameStatusRequest{GameId: addResp.GameId, NewStautus: gamev4.GameStatusType_PUBLISH}
		_, err = suite.GrpcClient.UpdateGameStatus(ctx, &reqUpdateGameStatus)
		require.NoError(t, err)

		reqUpdateGameStatus = gamev4.UpdateGameStatusRequest{GameId: addResp.GameId, NewStautus: gamev4.GameStatusType_PENDING}
		_, err = suite.GrpcClient.UpdateGameStatus(ctx, &reqUpdateGameStatus)
		hadlerchecker.AssertUpdateGameStatusInvalidStatus(t, err)
	})
	t.Run("Нельзя сменить жанр с PUBLISH на DRAFT", func(t *testing.T) {
		gameToAdd := random.WithOnlyRequireFields()
		addGameReq := gamev4.AddGameRequest{Game: gameToAdd}
		addResp, err := suite.GrpcClient.AddGame(ctx, &addGameReq)
		hadlerchecker.AssertAddGame(t, err, addResp)

		reqUpdateGameStatus := gamev4.UpdateGameStatusRequest{GameId: addResp.GameId, NewStautus: gamev4.GameStatusType_PENDING}
		_, err = suite.GrpcClient.UpdateGameStatus(ctx, &reqUpdateGameStatus)
		require.NoError(t, err)

		reqUpdateGameStatus = gamev4.UpdateGameStatusRequest{GameId: addResp.GameId, NewStautus: gamev4.GameStatusType_PUBLISH}
		_, err = suite.GrpcClient.UpdateGameStatus(ctx, &reqUpdateGameStatus)
		require.NoError(t, err)

		reqUpdateGameStatus = gamev4.UpdateGameStatusRequest{GameId: addResp.GameId, NewStautus: gamev4.GameStatusType_DRAFT}
		_, err = suite.GrpcClient.UpdateGameStatus(ctx, &reqUpdateGameStatus)
		hadlerchecker.AssertUpdateGameStatusInvalidStatus(t, err)
	})
}
