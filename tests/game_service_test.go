package tests

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/sariya23/game_service/internal/config"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSaveGame(t *testing.T) {
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

	gameToSave := gamev4.GameRequest{
		Title:       "Test Title",
		Description: "description",
		Genres:      []string{"TEST"},
		Tags:        []string{"TEST"},
		ReleaseDate: &date.Date{Year: 2025, Month: 2, Day: 24},
	}
	_, err = grpcClient.AddGame(ctx, &gamev4.AddGameRequest{Game: &gameToSave})
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}
}
