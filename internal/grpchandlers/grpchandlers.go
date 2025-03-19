package grpchandlers

import (
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/grpc"
)

type GameServicer interface {
}

type serverAPI struct {
	gamev4.UnimplementedGameServiceServer
	gameServicer GameServicer
}

func RegisterGrpcHandlers(grpcServer *grpc.Server, gameServicer GameServicer) {
	gamev4.RegisterGameServiceServer(grpcServer, &serverAPI{gameServicer: gameServicer})
}
