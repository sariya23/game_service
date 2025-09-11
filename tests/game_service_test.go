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
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestAddGame(t *testing.T) {
	t.Run("Успешное сохранение игры", func(t *testing.T) {
		ctx := context.Background()
		cfg := config.MustLoadByPath("../config/local.env")
		db := postgresql.MustNewConnection(ctx, mockslog.NewDiscardLogger(), cfg.Postgres.PostgresURL)
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
		gameToAdd := random.RandomAddGameRequest()
		availableTags, err := db.GetTags(ctx)
		require.NoError(t, err)
		gameToAdd.Tags = model.TagNames(availableTags)
		gameToAdd.CoverImage = nil
		gameToAdd.Genres = nil
		request := gamev4.AddGameRequest{Game: gameToAdd}
		resp, err := grpcClient.AddGame(ctx, &request)
		require.NoError(t, err)
		require.NotZero(t, resp.GetGameId())

		game, err := grpcClient.GetGame(ctx, &gamev4.GetGameRequest{GameId: resp.GetGameId()})
		require.NoError(t, err)
		require.Equal(t, gameToAdd.Title, game.Game.Title)
	})
}
