package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sariya23/game_service/internal/app/app"
	"github.com/sariya23/game_service/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	application := app.NewApp(ctx, log, cfg)
	go func() {
		application.MustRun()
	}()
	go func() {
		application.Stop(ctx, cancel, sigchan)
	}()
	log.Info("service is ready to recieve requests")
	<-ctx.Done()
	log.Info("service stopped gracefully")
}
