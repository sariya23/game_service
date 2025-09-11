package tests

import (
	"context"
	"io"
	"net"
	"strconv"
	"testing"

	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestAddGame(t *testing.T) {
	ctx := context.Background()
	cfg := config.MustLoadByPath("../config/local.env")
	db := postgresql.MustNewConnection(ctx, mockslog.NewDiscardLogger(), cfg.Postgres.PostgresURL)
	s3 := minioclient.MustPrepareMinio(ctx, mockslog.NewDiscardLogger(), cfg.Minio, false)
	conn, err := grpc.NewClient(
		net.JoinHostPort("127.0.0.1", strconv.Itoa(cfg.Server.GrpcServerPort)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil || conn == nil {
		t.Fatalf("cannot start client; err = %v", err)
	}
	grpcClient := gamev4.NewGameServiceClient(conn)
	if grpcClient == nil {
		t.Fatal("cannot create grpcClient")
	}
	t.Run("Успешное сохранение игры; все поля", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		tags, err := db.GetTags(ctx)
		require.NoError(t, err)
		genres, err := db.GetGenres(ctx)
		require.NoError(t, err)
		gameToAdd.Tags = model.TagNames(tags)
		gameToAdd.Genres = model.GenreNames(genres)
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage
		request := gamev4.AddGameRequest{Game: gameToAdd}
		resp, err := grpcClient.AddGame(ctx, &request)
		require.NoError(t, err)
		require.NotZero(t, resp.GetGameId())

		gameDB, err := db.GetGameByID(ctx, resp.GetGameId())
		require.NoError(t, err)

		assert.Equal(t, gameToAdd.GetTitle(), gameDB.Title)
		assert.Equal(t, gameToAdd.GetDescription(), gameDB.Description)
		assert.Equal(t, gameToAdd.GetReleaseDate().GetYear(), int32(gameDB.ReleaseDate.Year()))
		assert.Equal(t, gameToAdd.GetReleaseDate().GetMonth(), int32(gameDB.ReleaseDate.Month()))
		assert.Equal(t, gameToAdd.GetReleaseDate().GetDay(), int32(gameDB.ReleaseDate.Day()))
		assert.Equal(t, gameToAdd.GetTags(), model.TagNames(gameDB.Tags))
		assert.Equal(t, gameToAdd.GetGenres(), model.GenreNames(gameDB.Genres))
		image, err := s3.GetObject(ctx, minioclient.GameKey(gameToAdd.GetTitle(), int(gameToAdd.ReleaseDate.GetYear())))
		require.NoError(t, err)
		imageBytes, err := io.ReadAll(image)
		require.NoError(t, err)
		assert.Equal(t, gameToAdd.CoverImage, imageBytes)

	})
}
