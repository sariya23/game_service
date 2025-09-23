package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/model"
	checkers "github.com/sariya23/game_service/tests/checkers/handlers"
	"github.com/sariya23/game_service/tests/suite"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/require"
)

func TestGetGame(t *testing.T) {
	ctx, suite := suite.NewSuite(t)
	t.Run("Тест ручки GetGame; Успешное получение игры", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		tags, err := suite.Db.GetTags(ctx)
		require.NoError(t, err)
		genres, err := suite.Db.GetGenres(ctx)
		require.NoError(t, err)
		gameToAdd.Tags = model.GetTagNames(tags)
		gameToAdd.Genres = model.GetGenreNames(genres)
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage
		request := gamev4.AddGameRequest{Game: gameToAdd}
		addResp, err := suite.GrpcClient.AddGame(ctx, &request)
		require.NoError(t, err)
		require.NotZero(t, addResp.GetGameId())

		getResp, err := suite.GrpcClient.GetGame(ctx, &gamev4.GetGameRequest{GameId: addResp.GetGameId()})
		// не передавать s3
		checkers.AssertGetGame(ctx, t, gameToAdd, getResp, suite.S3, err)
	})
	t.Run("Тест ручки GetGame; Ошибка при получени несуществующей игры", func(t *testing.T) {
		resp, err := suite.GrpcClient.GetGame(ctx, &gamev4.GetGameRequest{GameId: uint64(gofakeit.IntRange(10000, 40000))})
		checkers.AssertGetGameNotFound(t, err, resp)
	})
}
