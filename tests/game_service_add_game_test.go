//go:build integrations
// +build integrations

package tests

import (
	"io"
	"testing"

	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/model"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	checkers "github.com/sariya23/game_service/tests/checkers"
	hadlerchecker "github.com/sariya23/game_service/tests/checkers/handlers"
	"github.com/sariya23/game_service/tests/suite"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/require"
)

func TestAddGame(t *testing.T) {
	ctx, suite := suite.NewSuite(t)
	t.Run("Тест ручки AddGame; Успешное сохранение игры; все поля", func(t *testing.T) {
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
		resp, err := suite.GrpcClient.AddGame(ctx, &request)
		hadlerchecker.AssertAddGame(t, err, resp)

		gameDB, err := suite.Db.GetGameByID(ctx, resp.GetGameId())
		require.NoError(t, err)

		image, err := suite.S3.GetObject(ctx, minioclient.GameKey(request.Game.GetTitle(), int(request.Game.ReleaseDate.GetYear())))
		require.NoError(t, err)
		imageBytes, err := io.ReadAll(image)
		require.NoError(t, err)
		checkers.AssertAddGameRequestAndDB(ctx, t, &request, *gameDB, imageBytes, gamev4.GameStatusType_DRAFT)

	})
	t.Run("Тест ручки AddGame; Игра не создается если передан хотя бы один несущетвующий тэг", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		tags, err := suite.Db.GetTags(ctx)
		require.NoError(t, err)
		gameToAdd.Tags = append(model.GetTagNames(tags), gameToAdd.Tags...)
		gameToAdd.Genres = nil
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage

		request := gamev4.AddGameRequest{Game: gameToAdd}
		resp, err := suite.GrpcClient.AddGame(ctx, &request)

		hadlerchecker.AssertAddGameTagNotFound(t, err, resp)
	})
	t.Run("Тест ручки AddGame; Игра не создается если передан хотя бы один несущетвующий жанр", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		genres, err := suite.Db.GetGenres(ctx)
		require.NoError(t, err)
		gameToAdd.Genres = append(model.GetGenreNames(genres), gameToAdd.Genres...)
		gameToAdd.Tags = nil
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage

		request := gamev4.AddGameRequest{Game: gameToAdd}
		resp, err := suite.GrpcClient.AddGame(ctx, &request)

		hadlerchecker.AssertAddGameGenreNotFound(t, err, resp)

	})
	t.Run("Тест ручки AddGame; Нельзя сохранить точно такую же игру", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Genres = nil
		gameToAdd.Tags = nil
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage
		request := gamev4.AddGameRequest{Game: gameToAdd}
		resp, err := suite.GrpcClient.AddGame(ctx, &request)
		require.NoError(t, err)
		require.NotZero(t, resp.GetGameId())

		duplicateGame := random.RandomAddGameRequest()
		duplicateGame.Title = gameToAdd.GetTitle()
		duplicateGame.ReleaseDate = gameToAdd.GetReleaseDate()

		duplicateRequest := gamev4.AddGameRequest{Game: duplicateGame}
		resp, err = suite.GrpcClient.AddGame(ctx, &duplicateRequest)

		hadlerchecker.AssertAddGameDuplicateGame(t, err, resp)
	})
}
