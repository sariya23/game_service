package grpchandlers

import (
	"context"
	"log/slog"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"google.golang.org/grpc"
)

type GameServicer interface {
	AddGame(ctx context.Context, game dto.AddGameHandler) (gameID int64, err error)
	GetGame(ctx context.Context, gameID int64) (game *model.Game, err error)
	GameList(ctx context.Context, gameFilters dto.GameFilters, limit uint32) (games []model.ShortGame, err error)
	DeleteGame(ctx context.Context, gameID int64) (deletedGameID int64, err error)
	UpdateGameStatus(ctx context.Context, gameID int64, newStatus game.GameStatusType) error
}

type serverAPI struct {
	game.UnimplementedGameServiceServer
	gameServicer GameServicer
	log          *slog.Logger
}

func RegisterGrpcHandlers(grpcServer *grpc.Server, gameServicer GameServicer, log *slog.Logger) {
	game.RegisterGameServiceServer(grpcServer, &serverAPI{gameServicer: gameServicer, log: log})
}
