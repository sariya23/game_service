package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sariya23/game_service/internal/app/app"
	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/lib/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.MustLoad()
	log := setUpLogger(cfg)
	log.Info("starting app", slog.String("env", cfg.Env.EnvType))
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

func setUpLogger(cfg *config.Config) *slog.Logger {
	logLevel := slog.LevelDebug
	if cfg.Env.EnvType == "prod" {
		logLevel = slog.LevelInfo
	}
	return logger.NewLogger(logLevel)
}
