package grcpgatewayapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/sariya23/api_game_service/gen/game"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcGatewayApp struct {
	server *http.Server
	log    *slog.Logger
	host   string
	port   int
}

func NewGrpcGatewayApp(ctx context.Context, log *slog.Logger, grpcPort, httpPort int, grpcHost, httpHost string, allowedOrigins string) *GrpcGatewayApp {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := game.RegisterGameServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", grpcHost, grpcPort), opts)
	if err != nil {
		panic(fmt.Sprintf("cannot register game service endpoints: %v", err))
	}
	corsHandler := setupCORS(mux, allowedOrigins, log)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", httpHost, httpPort),
		Handler: corsHandler,
	}

	return &GrpcGatewayApp{
		server: httpServer,
		log:    log,
		host:   httpHost,
		port:   httpPort,
	}
}

func setupCORS(next http.Handler, allowedOrigins string, log *slog.Logger) http.Handler {
	origins := parseOrigins(allowedOrigins)
	if len(origins) == 0 {
		log.Debug("CORS disabled: no allowed origins configured")
		return next
	}
	c := cors.New(cors.Options{
		AllowedOrigins: origins,
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		MaxAge:         3600,
	})

	log.Debug("CORS enabled", slog.Any("allowed_origins", origins))
	return c.Handler(next)
}

func parseOrigins(originsStr string) []string {
	if originsStr == "" {
		return []string{}
	}
	origins := strings.Split(originsStr, ",")
	result := make([]string, 0, len(origins))
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			result = append(result, origin)
		}
	}
	return result
}

func (gw *GrpcGatewayApp) MustRun() {
	gw.log.Info("grpc gateway listening", slog.String("host", gw.host), slog.Int("port", gw.port))
	if err := gw.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("error while executing grpc gateway: %v", err))
	}
}

func (gw *GrpcGatewayApp) Stop(ctx context.Context) {
	const operationPlace = "grpcgateway.Stop"
	log := gw.log.With("operationPlace", operationPlace)
	if err := gw.server.Shutdown(ctx); err != nil {
		log.Error(fmt.Sprintf("HTTP server shutdown failed: %v", err))
	}
}
