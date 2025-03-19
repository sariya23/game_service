package main

import (
	"log/slog"
	"os"

	"github.com/sariya23/game_service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	log.Info(
		"starting game service",
		slog.Int("grpc port", cfg.GrpcServerPort),
		slog.Int("htpp port", cfg.HttpServerPort),
	)
}
