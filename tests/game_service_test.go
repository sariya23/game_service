package tests

import (
	"context"
	"io"
	"net"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
	t.Run("Игра не создается если передан хотя бы один несущетвующий тэг", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		tags, err := db.GetTags(ctx)
		require.NoError(t, err)
		gameToAdd.Tags = append(model.TagNames(tags), gameToAdd.Tags...)
		gameToAdd.Genres = nil
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage
		request := gamev4.AddGameRequest{Game: gameToAdd}
		resp, err := grpcClient.AddGame(ctx, &request)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.TagNotFoundMessage, s.Message())
		require.Zero(t, resp.GetGameId())
	})
	t.Run("Игра не создается если передан хотя бы один несущетвующий жанр", func(t *testing.T) {
		gameToAdd := random.RandomAddGameRequest()
		genres, err := db.GetGenres(ctx)
		require.NoError(t, err)
		gameToAdd.Genres = append(model.GenreNames(genres), gameToAdd.Genres...)
		gameToAdd.Tags = nil
		expectedImage, err := random.Image()
		require.NoError(t, err)
		gameToAdd.CoverImage = expectedImage
		request := gamev4.AddGameRequest{Game: gameToAdd}
		resp, err := grpcClient.AddGame(ctx, &request)
		s, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, s.Code())
		require.Equal(t, outerror.GenreNotFoundMessage, s.Message())
		require.Zero(t, resp.GetGameId())
	})
	t.Run("Нельзя сохранить точно такую же игру", func(t *testing.T) {
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
		s, _ := status.FromError(err)
		require.Equal(t, codes.AlreadyExists, s.Code())
		require.Equal(t, outerror.GameAlreadyExistMessage, s.Message())
		require.Zero(t, resp.GetGameId())
	})
}

func TestGetGame(t *testing.T) {
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
	t.Run("Успешное получение игры", func(t *testing.T) {
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
		addResp, err := grpcClient.AddGame(ctx, &request)
		require.NoError(t, err)
		require.NotZero(t, addResp.GetGameId())

		getResp, err := grpcClient.GetGame(ctx, &gamev4.GetGameRequest{GameId: addResp.GetGameId()})
		require.NoError(t, err)
		assert.Equal(t, gameToAdd.GetTitle(), getResp.Game.GetTitle())
		assert.Equal(t, gameToAdd.GetDescription(), getResp.Game.GetDescription())
		assert.Equal(t, gameToAdd.GetReleaseDate().GetYear(), getResp.Game.GetReleaseDate().GetYear())
		assert.Equal(t, gameToAdd.GetReleaseDate().GetMonth(), getResp.Game.GetReleaseDate().GetMonth())
		assert.Equal(t, gameToAdd.GetReleaseDate().GetDay(), getResp.Game.GetReleaseDate().GetDay())
		assert.Equal(t, gameToAdd.GetGenres(), getResp.Game.GetGenres())
		assert.Equal(t, gameToAdd.GetTags(), getResp.Game.GetTags())
		obj, err := s3.GetObject(ctx, getResp.Game.GetCoverImageUrl())
		require.NoError(t, err)
		imageBytes, err := io.ReadAll(obj)
		require.NoError(t, err)
		assert.Equal(t, gameToAdd.GetCoverImage(), imageBytes)
	})
	t.Run("Ошибка при получени несуществующей игры", func(t *testing.T) {
		resp, err := grpcClient.GetGame(ctx, &gamev4.GetGameRequest{GameId: uint64(gofakeit.IntRange(10000, 40000))})
		s, _ := status.FromError(err)
		require.Equal(t, codes.NotFound, s.Code())
		require.Equal(t, outerror.GameNotFoundMessage, s.Message())
		require.Nil(t, resp.GetGame())
	})
}
