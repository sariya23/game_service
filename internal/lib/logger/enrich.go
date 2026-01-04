package logger

import (
	"context"
	"log/slog"

	"github.com/sariya23/game_service/internal/interceptors"
)

func EnrichRequestID(ctx context.Context, log *slog.Logger) *slog.Logger {
	requestID := ctx.Value(interceptors.RequestIDKey).(string)
	log = log.With("request_id", requestID)
	return log
}
