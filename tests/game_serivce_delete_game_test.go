//go:build integrations
// +build integrations

package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/lib/random"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	checkers "github.com/sariya23/game_service/tests/checkers/handlers"
	"github.com/sariya23/game_service/tests/suite"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/require"
)

func TestDeteteGame(t *testing.T) {
	ctx, suite := suite.NewSuite(t)
	t.Run("Тест ручки DeleteGame; Успешное удаление игры", func(t *testing.T) {
		gameToAdd := random.WithOnlyRequireFields()
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage
		request := gamev4.AddGameRequest{Game: gameToAdd}
		addResp, err := suite.GrpcClient.AddGame(ctx, &request)
		checkers.AssertAddGame(t, err, addResp)

		respDelete, err := suite.GrpcClient.DeleteGame(ctx, &gamev4.DeleteGameRequest{GameId: addResp.GameId})
		checkers.AssertDeleteGame(t, err, addResp.GameId, respDelete)

		respGet, err := suite.GrpcClient.GetGame(ctx, &gamev4.GetGameRequest{GameId: addResp.GameId})
		checkers.AssertGetGameNotFound(t, err, respGet)

		obj, err := suite.S3.GetObject(ctx, minioclient.GameKey(gameToAdd.Title, int(gameToAdd.ReleaseDate.Year)))
		require.Error(t, err)
		require.Empty(t, obj)
	})
	t.Run("Тест ручки DeleteGame; Удаление игры без обложки", func(t *testing.T) {
		gameToAdd := random.WithOnlyRequireFields()
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage

		request := gamev4.AddGameRequest{Game: gameToAdd}
		addResp, err := suite.GrpcClient.AddGame(ctx, &request)
		checkers.AssertAddGame(t, err, addResp)

		respDelete, err := suite.GrpcClient.DeleteGame(ctx, &gamev4.DeleteGameRequest{GameId: addResp.GameId})
		checkers.AssertDeleteGame(t, err, addResp.GameId, respDelete)
	})
	t.Run("Тест ручки DeleteGame; Игра не найдена", func(t *testing.T) {
		resp, err := suite.GrpcClient.DeleteGame(ctx, &gamev4.DeleteGameRequest{GameId: uint64(gofakeit.Int64())})
		checkers.AssertDeleteGameNotFound(t, err, resp)
	})
}
