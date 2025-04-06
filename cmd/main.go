package main

import (
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
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	log.Info(
		"starting game service",
		slog.Int("grpc port", cfg.GrpcServerPort),
		slog.Int("htpp port", cfg.HttpServerPort),
	)
	grpcServer := grpc.NewServer()
	db := postgresql.MustNewConnection(log)
	s3Client := minioclient.NewMinioClient(log, cfg.MinioBucket)
	mailer := email.NewDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.EmailUser, cfg.EmailPassword, cfg.AdminEmail)
	gameService := gameservice.NewGameService(log, db, s3Client, db, mailer)
	grpchandlers.RegisterGrpcHandlers(grpcServer, gameService)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcServerPort))
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}
