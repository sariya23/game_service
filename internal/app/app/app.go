package app

import (
	"context"
	"log/slog"

	"github.com/sariya23/game_service/internal/app/grcpgatewayapp"
	"github.com/sariya23/game_service/internal/app/grpcserviceapp"
	"github.com/sariya23/game_service/internal/config"
	gameservice "github.com/sariya23/game_service/internal/service/game"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
)

type App struct {
	Config         *config.Config
	Db             postgresql.PostgreSQL
	Minio          *minioclient.Minio
	GrpcApp        *grpcserviceapp.GrpcServer
	GrpcGateWayApp *grcpgatewayapp.GrpcGatewayApp
}

func NewApp(ctx context.Context, log *slog.Logger, cfg *config.Config) *App {
	db := postgresql.MustNewConnection(ctx, log, cfg.Postgres.PostgresURL)
	s3Client := minioclient.MustPrepareMinio(ctx, log, cfg.Minio, false)
	gameService := gameservice.NewGameService(log, db, db, db, s3Client)
	grpcApp := grpcserviceapp.NewGrpcServer(log, cfg.Server.GrpcServerPort, cfg.Server.GRPCServerHost, gameService)
	gwApp := grcpgatewayapp.NewGrpcGatewayApp(ctx, log, cfg.Server.GrpcServerPort, cfg.Server.HttpServerPort, cfg.Server.GRPCServerHost, cfg.Server.HttpServerHost)
	return &App{
		Config:         cfg,
		Db:             db,
		Minio:          s3Client,
		GrpcApp:        grpcApp,
		GrpcGateWayApp: gwApp,
	}
}
