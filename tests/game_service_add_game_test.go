package tests

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	checkers "github.com/sariya23/game_service/tests/checkers/handlers"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
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
	t.Run("Тест ручки AddGame; Успешное сохранение игры; все поля", func(t *testing.T) {
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
		resp, err := grpcClient.AddGame(ctx, &request)
		checkers.AssertAddGame(t, err, resp)

		gameDB, err := db.GetGameByID(ctx, resp.GetGameId())
		require.NoError(t, err)
		checkers.AssertAddGameRequestAndDB(ctx, t, &request, *gameDB, s3)

	})
	t.Run("Тест ручки AddGame; Игра не создается если передан хотя бы один несущетвующий тэг", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		tags, err := db.GetTags(ctx)
		require.NoError(t, err)
		gameToAdd.Tags = append(model.GetTagNames(tags), gameToAdd.Tags...)
		gameToAdd.Genres = nil
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage

		request := gamev4.AddGameRequest{Game: gameToAdd}
		resp, err := grpcClient.AddGame(ctx, &request)

		checkers.AssertAddGameTagNotFound(t, err, resp)
	})
	t.Run("Тест ручки AddGame; Игра не создается если передан хотя бы один несущетвующий жанр", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		genres, err := db.GetGenres(ctx)
		require.NoError(t, err)
		gameToAdd.Genres = append(model.GetGenreNames(genres), gameToAdd.Genres...)
		gameToAdd.Tags = nil
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage

		request := gamev4.AddGameRequest{Game: gameToAdd}
		resp, err := grpcClient.AddGame(ctx, &request)

		checkers.AssertAddGameGenreNotFound(t, err, resp)

	})
	t.Run("Тест ручки AddGame; Нельзя сохранить точно такую же игру", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Genres = nil
		gameToAdd.Tags = nil
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage
		request := gamev4.AddGameRequest{Game: gameToAdd}
		resp, err := grpcClient.AddGame(ctx, &request)
		require.NoError(t, err)
		require.NotZero(t, resp.GetGameId())

		duplicateGame := random.RandomAddGameRequest()
		duplicateGame.Title = gameToAdd.GetTitle()
		duplicateGame.ReleaseDate = gameToAdd.GetReleaseDate()

		duplicateRequest := gamev4.AddGameRequest{Game: duplicateGame}
		resp, err = grpcClient.AddGame(ctx, &duplicateRequest)

		checkers.AssertAddGameDuplicateGame(t, err, resp)
	})
}
