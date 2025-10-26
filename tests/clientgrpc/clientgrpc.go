package clientgrpc

import (
	"fmt"
	"path/filepath"
	"testing"

	game_api "github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GameServiceTestClient struct {
	cl game_api.GameServiceClient
}

func NewGameServiceTestClient() *GameServiceTestClient {
	cfg := config.MustLoadByPath(filepath.Join("..", "..", "..", "..", "config", "test.env"))
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Server.GRPCServerHost, cfg.Server.GrpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	grpcClient := game_api.NewGameServiceClient(conn)
	if grpcClient == nil {
		panic(err)
	}
	return &GameServiceTestClient{grpcClient}
}

func (g *GameServiceTestClient) GetClient() game_api.GameServiceClient {
	return g.cl
}

func (g *GameServiceTestClient) SetUp(t *testing.T) {
	t.Helper()
}

func (g *GameServiceTestClient) TearDown(t *testing.T) {
	t.Helper()
}
