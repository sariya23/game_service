//go:build integrations

package game_test

import (
	"context"
	"io"
	"sort"
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

func TestGetGame(t *testing.T) {
	t.Run("Успешное получение игры", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		responseSave, err := client.GetClient().AddGame(ctx, &gamev2.AddGameRequest{Game: &gameToAdd})
		require.NoError(t, err)
		require.NotZero(t, responseSave.GameId)
		request := gamev2.GetGameRequest{GameId: responseSave.GameId}

		response, err := client.GetClient().GetGame(ctx, &request)

		require.NoError(t, err)
		assert.Equal(t, responseSave.GameId, response.GetGame().ID)
		assert.Equal(t, gameToAdd.Title, response.GetGame().Title)
		assert.Equal(t, gameToAdd.Description, response.GetGame().Description)
		assert.Equal(t, gameToAdd.ReleaseDate.String(), response.GetGame().ReleaseDate.String())
		reader, err := minioT.GetClient().GetObject(ctx, minioT.BucketName, response.Game.CoverImageUrl, minio.GetObjectOptions{})
		require.NoError(t, err)
		defer reader.Close()
		imageData, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, gameToAdd.CoverImage, imageData)

		expectedGenres := dbT.GetGenresByNames(ctx, gameToAdd.Genres)
		sort.Slice(expectedGenres, func(i, j int) bool {
			return expectedGenres[i].GenreName < expectedGenres[j].GenreName
		})
		actualGenres := response.GetGame().Genres
		sort.Slice(actualGenres, func(i, j int) bool {
			return actualGenres[i] < actualGenres[j]
		})
		assert.Equal(t, model.GenreNames(expectedGenres), actualGenres)

		expectedTags := dbT.GetTagsByNames(ctx, gameToAdd.Tags)
		sort.Slice(expectedTags, func(i, j int) bool {
			return expectedTags[i].TagName < expectedTags[j].TagName
		})
		actualTags := response.GetGame().Tags
		sort.Slice(actualTags, func(i, j int) bool {
			return actualTags[i] < actualTags[j]
		})
		assert.Equal(t, model.TagNames(expectedTags), actualTags)
	})
	t.Run("Нет игры с таким ID", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		nonExistGameID := gofakeit.Int64()
		request := gamev2.GetGameRequest{GameId: nonExistGameID}

		response, err := client.GetClient().GetGame(ctx, &request)
		require.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Nil(t, response)
	})
	t.Run("Отрицательный ID", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		nonExistGameID := gofakeit.Uint32()
		request := gamev2.GetGameRequest{GameId: -int64(nonExistGameID)}

		response, err := client.GetClient().GetGame(ctx, &request)
		require.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Nil(t, response)
	})
}
