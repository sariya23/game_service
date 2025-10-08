//go:build integrations

package game_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
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
		gameToAdd := random.GameToAddService(tagIDs, genreIDs)
		gameID := dbT.InsertGame(ctx, gameToAdd)
		dbT.InsertGameGenre(ctx, gameID, gameToAdd.GenreIDs)
		dbT.InsertGameTag(ctx, gameID, gameToAdd.TagIDs)
		request := gamev2.DeleteGameRequest{GameId: gameID}

		response, err := client.GetClient().DeleteGame(ctx, &request)

		require.NoError(t, err)
		assert.Equal(t, gameID, response.GameId)
		assert.Len(t, dbT.GetGameGenreByGameID(ctx, gameID), 0)
		assert.Len(t, dbT.GetGameTagByGameID(ctx, gameID), 0)
		assert.Nil(t, dbT.GetGameById(ctx, gameID))
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
