package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/sariya23/game_service/internal/app/grcpgatewayapp"
	"github.com/sariya23/game_service/internal/app/grpcserviceapp"
	"github.com/sariya23/game_service/internal/config"
	gameservice "github.com/sariya23/game_service/internal/service/game"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
)

type App struct {
	log            *slog.Logger
	Config         *config.Config
	Db             *postgresql.PostgreSQL
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
		log:            log,
	}
}

func (a *App) MustRun() {
	runActions := []struct {
		action func()
		errMsg string
	}{
		{action: a.GrpcApp.MustRun, errMsg: "Error while starting gRPC server"},
		{action: a.GrpcGateWayApp.MustRun, errMsg: "Error while starting grpc-gateway"},
	}
	wg := sync.WaitGroup{}
	wg.Add(len(runActions))

	for _, action := range runActions {
		go func() {
			defer wg.Done()
			action.action()
		}()
	}

	wg.Wait()
	a.log.Info("App is running")
}

func (a *App) Stop(ctx context.Context, cancel context.CancelFunc, sigchan <-chan os.Signal) {
	const operationPlace = "app.Stop"
	log := a.log.With("operationPlace", operationPlace)
	sig := <-sigchan
	log.Info(fmt.Sprintf("Received signal: %v, shutting down...\n", sig))
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()
	defer cancel()
	a.GrpcGateWayApp.Stop(shutdownCtx)
	log.Info("Grpc gateway server stopped")
	a.GrpcApp.Stop()
	log.Info("GRPC server stopped")
	a.Db.Close()
	log.Info("DB closed")
}
