//go:build integrations

package game_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	game_api "github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/tests/clientgrpc"
	"github.com/sariya23/game_service/tests/utils/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateGameStatus(t *testing.T) {
	t.Run("Успешное обновление статуса из DRAFT в PENDING", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		responseSave, err := client.GetClient().AddGame(ctx, &game_api.AddGameRequest{Game: gameToAdd})
		require.NoError(t, err)
		require.NotZero(t, responseSave.GameId)

		_, err = client.GetClient().UpdateGameStatus(ctx, &game_api.UpdateGameStatusRequest{GameId: responseSave.GameId, NewStatus: game_api.GameStatusType_PENDING})
		require.NoError(t, err)
		game := dbT.GetGameById(ctx, responseSave.GameId)
		assert.Equal(t, game_api.GameStatusType_PENDING, game.GameStatus)
	})
	t.Run("Успешное обновление статуса из PENDING в PUBLISH", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		responseSave, err := client.GetClient().AddGame(ctx, &game_api.AddGameRequest{Game: gameToAdd})
		require.NoError(t, err)
		require.NotZero(t, responseSave.GameId)
		dbT.UpdateGameStatus(ctx, responseSave.GameId, game_api.GameStatusType_PENDING)

		_, err = client.GetClient().UpdateGameStatus(ctx, &game_api.UpdateGameStatusRequest{GameId: responseSave.GameId, NewStatus: game_api.GameStatusType_PUBLISH})
		require.NoError(t, err)
		game := dbT.GetGameById(ctx, responseSave.GameId)
		assert.Equal(t, game_api.GameStatusType_PUBLISH, game.GameStatus)
	})
	t.Run("Нельзя перевести из статуса DRAFT в PUBLISH", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		responseSave, err := client.GetClient().AddGame(ctx, &game_api.AddGameRequest{Game: gameToAdd})
		require.NoError(t, err)
		require.NotZero(t, responseSave.GameId)

		_, err = client.GetClient().UpdateGameStatus(ctx, &game_api.UpdateGameStatusRequest{GameId: responseSave.GameId, NewStatus: game_api.GameStatusType_PUBLISH})
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
	t.Run("Отрицательный айди игры", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		_, err := client.GetClient().UpdateGameStatus(ctx, &game_api.UpdateGameStatusRequest{GameId: -gofakeit.Int64(), NewStatus: game_api.GameStatusType_PUBLISH})
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}
