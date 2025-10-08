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
			return gameDB.Genres[i].GenreID < gameDB.Genres[j].GenreID
		})
		genresExpected := dbT.GetGenresByNames(ctx, gameToAdd.Genres)
		sort.Slice(genresExpected, func(i, j int) bool {
			return genresExpected[i].GenreID < genresExpected[j].GenreID
		})
		assert.Equal(t, genresExpected, gameDB.Genres)
		tagsExpected := dbT.GetTagsByNames(ctx, gameToAdd.Tags)
		sort.Slice(tagsExpected, func(i, j int) bool {
			return tagsExpected[i].TagID < tagsExpected[j].TagID
		})
		sort.Slice(gameDB.Tags, func(i, j int) bool {
			return gameDB.Tags[i].TagID < gameDB.Tags[j].TagID
		})
		assert.Equal(t, tagsExpected, gameDB.Tags)
		reader, err := minioT.GetClient().GetObject(ctx, minioT.BucketName, gameDB.ImageURL, minio.GetObjectOptions{})
		require.NoError(t, err)
		defer reader.Close()
		imageData, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, gameToAdd.CoverImage, imageData)
	})
}
