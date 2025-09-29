package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/sariya23/game_service/internal/app/app"
	"github.com/sariya23/game_service/internal/config"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	application := app.NewApp(ctx, log, cfg)
	
	application.GrpcApp.MustRun()
}
