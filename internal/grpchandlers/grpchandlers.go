package grpchandlers

import (
	"context"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	gamev2 "github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc"
)

type GameServicer interface {
	AddGame(ctx context.Context, game dto.AddGameHandler) (gameID int64, err error)
	GetGame(ctx context.Context, gameID int64) (game *model.Game, err error)
	GameList(ctx context.Context, gameFilters dto.GameFilters, limit uint32) (games []model.ShortGame, err error)
	DeleteGame(ctx context.Context, gameID int64) (deletedGameID int64, err error)
	UpdateGameStatus(ctx context.Context, gameID int64, newStatus gamev2.GameStatusType) error
}

type serverAPI struct {
	gamev2.UnimplementedGameServiceServer
	gameServicer GameServicer
}

func RegisterGrpcHandlers(grpcServer *grpc.Server, gameServicer GameServicer) {
	gamev2.RegisterGameServiceServer(grpcServer, &serverAPI{gameServicer: gameServicer})
}
