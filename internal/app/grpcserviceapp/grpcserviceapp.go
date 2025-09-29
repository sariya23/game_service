package grpcserviceapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/sariya23/game_service/internal/grpchandlers"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	port   int
	host   string
	server *grpc.Server
	log    *slog.Logger
}

func NewGrpcServer(log *slog.Logger, port int, host string, implementation grpchandlers.GameServicer) *GrpcServer {
	grpcServer := grpc.NewServer()
	grpchandlers.RegisterGrpcHandlers(grpcServer, implementation)
	return &GrpcServer{
		port:   port,
		host:   host,
		server: grpcServer,
		log:    log,
	}
}

func (g *GrpcServer) MustRun() {
	const operationPlace = "grpcserviceapp.MustRun"
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", g.host, g.port))
	if err != nil {
		panic(fmt.Sprintf("%s: %v", operationPlace, err))
	}
	g.log.Info("grpc server listening", slog.String("host", g.host), slog.Int("port", g.port))
	if err := g.server.Serve(listener); err != nil {
		panic(fmt.Sprintf("%s: %v", operationPlace, err))
	}
}

func (g *GrpcServer) Stop() {
	g.server.GracefulStop()
}
