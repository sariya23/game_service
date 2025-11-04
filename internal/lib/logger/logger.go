package logger

import (
	"log/slog"
	"os"
)

func NewLogger(logLvl slog.Level) *slog.Logger {
	sl := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLvl}))
	return sl
}
