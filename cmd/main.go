package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/grpchandlers"
	gameservice "github.com/sariya23/game_service/internal/service/game"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	go func() {
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err := gamev4.RegisterGameServiceHandlerFromEndpoint(
			ctx, mux,
			fmt.Sprintf("%s:%d", cfg.Server.GRPCServerHost, cfg.Server.GrpcServerPort),
			opts)
		if err != nil {
			panic(fmt.Sprintf("failed to register game service handler: %v", err.Error()))
		}
		log.Info(
			"starting grpc gateway",
			slog.String("http host", cfg.Server.HttpServerHost),
			slog.Int("http port", cfg.Server.HttpServerPort),
		)
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Server.HttpServerHost, cfg.Server.HttpServerPort), mux); err != nil {
			panic(fmt.Sprintf("failed to start gRPC gateway: %v", err.Error()))
		}
	}()
	log.Info("grcp gateway ready to get connections")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GrpcServerPort))
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}
