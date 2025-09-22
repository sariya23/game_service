package tests

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	"github.com/sariya23/game_service/tests/checkers"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDeteteGame(t *testing.T) {
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
	t.Run("Тест ручки DeleteGame; Успешное удаление игры", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		tags, err := db.GetTags(ctx)
		require.NoError(t, err)
		genres, err := db.GetGenres(ctx)
		require.NoError(t, err)
		gameToAdd.Tags = model.GetTagNames(tags)
		gameToAdd.Genres = model.GetGenreNames(genres)
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage
		request := gamev4.AddGameRequest{Game: gameToAdd}
		addResp, err := grpcClient.AddGame(ctx, &request)
		checkers.AssertAddGame(t, err, addResp)

		respDelete, err := grpcClient.DeleteGame(ctx, &gamev4.DeleteGameRequest{GameId: addResp.GameId})
		checkers.AssertDeleteGame(t, err, addResp.GameId, respDelete)

		respGet, err := grpcClient.GetGame(ctx, &gamev4.GetGameRequest{GameId: addResp.GameId})
		checkers.AssertGetGameNotFound(t, err, respGet)

		obj, err := s3.GetObject(ctx, minioclient.GameKey(gameToAdd.Title, int(gameToAdd.ReleaseDate.Year)))
		require.Error(t, err)
		require.Empty(t, obj)
	})
	t.Run("Тест ручки DeleteGame; Удаление игры без обложки", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		tags, err := db.GetTags(ctx)
		require.NoError(t, err)
		genres, err := db.GetGenres(ctx)
		require.NoError(t, err)
		gameToAdd.Tags = model.GetTagNames(tags)
		gameToAdd.Genres = model.GetGenreNames(genres)
		gameToAdd.CoverImage = nil

		request := gamev4.AddGameRequest{Game: gameToAdd}
		addResp, err := grpcClient.AddGame(ctx, &request)
		checkers.AssertAddGame(t, err, addResp)

		respDelete, err := grpcClient.DeleteGame(ctx, &gamev4.DeleteGameRequest{GameId: addResp.GameId})
		checkers.AssertDeleteGame(t, err, addResp.GameId, respDelete)
	})
	t.Run("Тест ручки DeleteGame; Игра не найдена", func(t *testing.T) {
		resp, err := grpcClient.DeleteGame(ctx, &gamev4.DeleteGameRequest{GameId: uint64(gofakeit.Uint64())})
		checkers.AssertDeleteGameNotFound(t, err, resp)
	})
}
