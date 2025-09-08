package tests

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/sariya23/game_service/internal/config"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestAddGame(t *testing.T) {
	t.Run("Успешное сохранение игры", func(t *testing.T) {
		ctx := context.Background()
		cfg := config.MustLoadByPath("../config/local.env")
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
		request := gamev4.AddGameRequest{
			Game: &gamev4.GameRequest{
				Title:       "test",
				Genres:      []string{"Экшен"},
				Tags:        []string{"Пиксельная графика"},
				ReleaseDate: &date.Date{Year: 2024, Month: 3, Day: 2},
				Description: "test",
				CoverImage:  nil,
			},
		}
		_, err = grpcClient.AddGame(ctx, &request)
		require.NoError(t, err)

	})
}
