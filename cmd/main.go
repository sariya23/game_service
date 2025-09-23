package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/grpchandlers"
	"github.com/sariya23/game_service/internal/lib/email"
	gameservice "github.com/sariya23/game_service/internal/service/game"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	log.Info(
		"starting game service",
		slog.Int("grpc port", cfg.Server.GrpcServerPort),
		slog.Int("htpp port", cfg.Server.HttpServerPort),
	)
	grpcServer := grpc.NewServer()
	db := postgresql.MustNewConnection(ctx, log, cfg.Postgres.PostgresURL)
	s3Client := minioclient.MustPrepareMinio(ctx, log, cfg.Minio, false)
	var mailer gameservice.EmailAlerter
	if cfg.Env.EnvType == config.TestEnvType {
		mailer = email.NewDialerMock(cfg.Email)
	} else {
		mailer = email.NewDialer(cfg.Email)
	}
	gameService := gameservice.NewGameService(log, db, db, db, s3Client, mailer)
	grpchandlers.RegisterGrpcHandlers(grpcServer, gameService)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GrpcServerPort))
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}
