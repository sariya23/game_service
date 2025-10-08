//go:build integrations

package game_test

import (
	"context"
	"io"
	"sort"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/sariya23/game_service/internal/lib/converters"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/tests/clientgrpc"
	"github.com/sariya23/game_service/tests/utils/random"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddGame(t *testing.T) {
	t.Run("Успешное добавление игры, все поля", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		request := gamev2.AddGameRequest{Game: &gameToAdd}

		response, err := client.GetClient().AddGame(ctx, &request)

		require.NoError(t, err)
		assert.NotZero(t, response.GameId)
		gameDB := dbT.GetGameById(ctx, response.GameId)

		assert.Equal(t, gameToAdd.Title, gameDB.Title)
		assert.Equal(t, gameToAdd.Description, gameDB.Description)
		assert.Equal(t, converters.FromProtoDate(gameToAdd.ReleaseDate), gameDB.ReleaseDate)
		assert.Equal(t, gamev2.GameStatusType_DRAFT, gameDB.GameStatus)
		sort.Slice(gameDB.Genres, func(i, j int) bool {
			return gameDB.Genres[i].GenreName < gameDB.Genres[j].GenreName
		})
		sort.Slice(gameToAdd.Genres, func(i, j int) bool {
			return gameToAdd.Genres[i] < gameToAdd.Genres[j]
		})
		assert.Equal(t, gameToAdd.Genres, model.GenreNames(gameDB.Genres))
		sort.Slice(gameDB.Tags, func(i, j int) bool {
			return gameDB.Tags[i].TagName < gameDB.Tags[j].TagName
		})
		sort.Slice(gameToAdd.Tags, func(i, j int) bool {
			return gameToAdd.Tags[i] < gameToAdd.Tags[j]
		})
		assert.Equal(t, gameToAdd.Tags, model.TagNames(gameDB.Tags))
		reader, err := minioT.GetClient().GetObject(ctx, minioT.BucketName, gameDB.ImageURL, minio.GetObjectOptions{})
		require.NoError(t, err)
		defer reader.Close()
		imageData, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, gameToAdd.CoverImage, imageData)
	})
}
