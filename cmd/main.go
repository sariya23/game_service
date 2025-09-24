package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/grpchandlers"
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
	log.Info("env is", slog.String("env", cfg.Env.EnvType))
	grpcServer := grpc.NewServer()
	db := postgresql.MustNewConnection(ctx, log, cfg.Postgres.PostgresURL)
	s3Client := minioclient.MustPrepareMinio(ctx, log, cfg.Minio, false)

	gameService := gameservice.NewGameService(log, db, db, db, s3Client)
	grpchandlers.RegisterGrpcHandlers(grpcServer, gameService)
	log.Info("server ready to get connections")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GrpcServerPort))
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}
