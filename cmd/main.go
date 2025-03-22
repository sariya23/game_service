package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/grpchandlers"
	"github.com/sariya23/game_service/internal/lib/kafka"
	gameservice "github.com/sariya23/game_service/internal/service/game"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	"github.com/sariya23/game_service/internal/storage/s3"
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
	kafkaProducer := kafka.MustNewKafkaProducer([]string{""}, "qwe")
	db := postgresql.MustNewConnection(log)
	s3Client := s3.NewS3Storage(log)
	gameService := gameservice.NewGameService(log, kafkaProducer, db, s3Client, db)
	grpchandlers.RegisterGrpcHandlers(grpcServer, gameService)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcServerPort))
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}
