package clientgrpc

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GameServiceTestClient struct {
	cl gamev2.GameServiceClient
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
	grpcClient := gamev2.NewGameServiceClient(conn)
	if grpcClient == nil {
		panic(err)
	}
	return &GameServiceTestClient{grpcClient}
}

func (g *GameServiceTestClient) GetClient() gamev2.GameServiceClient {
	return g.cl
}

func (g *GameServiceTestClient) SetUp(t *testing.T) {
	t.Helper()
}

func (g *GameServiceTestClient) TearDown(t *testing.T) {
	t.Helper()
}
