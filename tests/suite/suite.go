package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suit struct {
	*testing.T
	Cfg        config.Config
	Db         postgresql.PostgreSQL
	S3         *minioclient.Minio
	GrpcClient gamev4.GameServiceClient
}

func NewSuite(t *testing.T) (context.Context, Suit) {
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
	return ctx, Suit{T: t, Cfg: cfg, Db: db, GrpcClient: grpcClient, S3: s3}
}
