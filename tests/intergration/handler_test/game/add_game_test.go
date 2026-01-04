//go:build integrations

package game_test

import (
	"context"
	"io"
	"sort"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/minio/minio-go/v7"
	game_api "github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/lib/converters"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/tests/clientgrpc"
	"github.com/sariya23/game_service/tests/utils/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddGame(t *testing.T) {
	t.Run("Успешное добавление игры, все поля", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		request := game_api.AddGameRequest{Game: gameToAdd}

		response, err := client.GetClient().AddGame(ctx, &request)

		require.NoError(t, err)
		assert.NotZero(t, response.GameId)
		gameNoImageURL := dbT.GetGameById(ctx, response.GameId)
		imageURL, err := minioT.GetClient().PresignedGetObject(ctx, minioT.BucketName, gameNoImageURL.ImageKey, time.Minute, nil)
		require.NoError(t, err)

		gameDB := gameNoImageURL.ToDomain(imageURL.String())

		assert.Equal(t, gameToAdd.Title, gameDB.Title)
		assert.Equal(t, gameToAdd.Description, gameDB.Description)
		assert.Equal(t, converters.FromProtoDate(gameToAdd.ReleaseDate), gameDB.ReleaseDate)
		assert.Equal(t, game_api.GameStatusType_DRAFT, gameDB.GameStatus)
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
	t.Run("Успешное добавление игры, только обязательные поля", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		gameToAdd.Genres = []string{}
		gameToAdd.Tags = []string{}
		gameToAdd.CoverImage = []byte{}
		request := game_api.AddGameRequest{Game: gameToAdd}

		response, err := client.GetClient().AddGame(ctx, &request)

		require.NoError(t, err)
		assert.NotZero(t, response.GameId)
		gameNoImageURL := dbT.GetGameById(ctx, response.GameId)
		imageURL, err := minioT.GetClient().PresignedGetObject(ctx, minioT.BucketName, gameNoImageURL.ImageKey, time.Minute, nil)
		require.NoError(t, err)

		gameDB := gameNoImageURL.ToDomain(imageURL.String())

		assert.Equal(t, gameToAdd.Title, gameDB.Title)
		assert.Equal(t, gameToAdd.Description, gameDB.Description)
		assert.Equal(t, converters.FromProtoDate(gameToAdd.ReleaseDate), gameDB.ReleaseDate)
		assert.Equal(t, game_api.GameStatusType_DRAFT, gameDB.GameStatus)
		assert.Nil(t, gameDB.Genres)
		assert.Nil(t, gameDB.Tags)
		assert.Empty(t, gameDB.ImageURL)
	})
	t.Run("Если не переданы обязательные поля, игра не создается", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		cases := []struct {
			name      string
			gameToAdd *game_api.GameRequest
		}{
			{
				name: "no title",
				gameToAdd: func() *game_api.GameRequest {
					gameToAdd.Title = ""
					return gameToAdd
				}(),
			},
			{
				name: "no description",
				gameToAdd: func() *game_api.GameRequest {
					gameToAdd.Description = ""
					return gameToAdd
				}(),
			},
			{
				name: "no release date",
				gameToAdd: func() *game_api.GameRequest {
					gameToAdd.ReleaseDate = nil
					return gameToAdd
				}(),
			},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				request := game_api.AddGameRequest{Game: tc.gameToAdd}
				response, err := client.GetClient().AddGame(ctx, &request)
				st, _ := status.FromError(err)
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, st.Code())
				assert.Nil(t, response)
			})
		}
	})
	t.Run("Нельзя создать игру с таким же названием и датой выпуска (дубликат)", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		request := game_api.AddGameRequest{Game: gameToAdd}
		response, err := client.GetClient().AddGame(ctx, &request)
		require.NoError(t, err)
		assert.NotZero(t, response.GameId)

		request = game_api.AddGameRequest{Game: gameToAdd}
		response, err = client.GetClient().AddGame(ctx, &request)

		st, _ := status.FromError(err)
		assert.Equal(t, codes.AlreadyExists, st.Code())
		assert.Nil(t, response)
	})
	t.Run("Нельзя создать игру с несуществующими жанрами", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		gameToAdd.Genres = append(gameToAdd.Genres, gofakeit.LetterN(30))
		request := game_api.AddGameRequest{Game: gameToAdd}

		response, err := client.GetClient().AddGame(ctx, &request)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Nil(t, response)
	})
	t.Run("Нельзя создать игру с несуществующими тэгами", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
		gameToAdd.Tags = append(gameToAdd.Tags, gofakeit.LetterN(30))
		request := game_api.AddGameRequest{Game: gameToAdd}

		response, err := client.GetClient().AddGame(ctx, &request)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Nil(t, response)
	})
}
