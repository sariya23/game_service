package interceptors

import (
	"context"

	"github.com/sariya23/game_service/internal/lib/generate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const (
	RequestIDKey contextKey = "x-request-id"
)

func RequestIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var requestID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ids := md.Get("x-request-id"); len(ids) > 0 {
			requestID = ids[0]
		}
	}
	if requestID == "" {
		requestID = generate.GenerateRequestID()
	}

	ctx = context.WithValue(ctx, "x-request-id", requestID)
	resp, err := handler(ctx, req)
	return resp, err
}
