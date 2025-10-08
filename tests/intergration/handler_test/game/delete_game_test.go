//go:build integrations

package game_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/minio/minio-go/v7"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/tests/clientgrpc"
	"github.com/sariya23/game_service/tests/utils/random"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteGame(t *testing.T) {
	t.Run("Успешное удаление игры", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		respAddGame, err := client.GetClient().AddGame(ctx, &gamev2.AddGameRequest{Game: &gameToAdd})
		require.NoError(t, err)
		assert.NotZero(t, respAddGame.GameId)
		game := dbT.GetGameById(ctx, respAddGame.GameId)
		request := gamev2.DeleteGameRequest{GameId: respAddGame.GameId}

		response, err := client.GetClient().DeleteGame(ctx, &request)

		require.NoError(t, err)
		assert.Equal(t, respAddGame.GameId, response.GameId)
		assert.Len(t, dbT.GetGameGenreByGameID(ctx, respAddGame.GameId), 0)
		assert.Len(t, dbT.GetGameTagByGameID(ctx, respAddGame.GameId), 0)
		assert.Nil(t, dbT.GetGameById(ctx, respAddGame.GameId))
		_, err = minioT.GetClient().StatObject(ctx, minioT.BucketName, game.ImageURL, minio.GetObjectOptions{})
		require.Equal(t, "The specified key does not exist.", err.Error())
	})
	t.Run("Удаление несуществующей игры, возвращается ошибка", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		gameID := gofakeit.Int64()
		request := gamev2.DeleteGameRequest{GameId: gameID}

		response, err := client.GetClient().DeleteGame(ctx, &request)
		st, _ := status.FromError(err)
		require.Error(t, err)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Nil(t, response)
	})
}
