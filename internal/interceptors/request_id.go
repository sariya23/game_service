package interceptors

import (
	"context"

	"github.com/sariya23/game_service/internal/lib/generate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const (
	RequestIDKey = "request_id"
)

func RequestIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var requestID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ids := md.Get(RequestIDKey); len(ids) > 0 {
			requestID = ids[0]
		}
	}
	if requestID == "" {
		requestID = generate.GenerateRequestID()
	}

	ctx = context.WithValue(ctx, RequestIDKey, requestID)
	resp, err := handler(ctx, req)
	return resp, err
}
