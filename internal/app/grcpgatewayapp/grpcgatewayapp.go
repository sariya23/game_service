package grcpgatewayapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcGatewayApp struct {
	server *http.Server
	log    *slog.Logger
	host   string
	port   int
}

func NewGrpcGatewayApp(ctx context.Context, log *slog.Logger, grpcPort, httpPort int, grpcHost, httpHost string) *GrpcGatewayApp {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewBundle().TransportCredentials())}
	err := gamev4.RegisterGameServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", grpcHost, grpcPort), opts)
	if err != nil {
		panic(fmt.Sprintf("cannot register game service endpoints: %v", err))
	}
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", httpHost, httpPort),
		Handler: mux,
	}
	return &GrpcGatewayApp{
		server: httpServer,
		log:    log,
		host:   httpHost,
		port:   httpPort,
	}
}

func (gw *GrpcGatewayApp) MustRun() {
	gw.log.Info("grpc gateway listening", slog.String("host", gw.host), slog.Int("port", gw.port))
	if err := gw.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("error while executing grpc gateway: %v", err))
	}
}
