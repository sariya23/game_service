package grpchandlers

import (
	"context"

	"github.com/sariya23/game_service/internal/model"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameServicer interface {
	AddGame(ctx context.Context, game *gamev4.Game) (gameId uint64, err error)
	GetGame(ctx context.Context, gameTitle string) (game *gamev4.Game, err error)
	GetTopGames(ctx context.Context, gameFilters model.GameFilters, limit uint32) (games []gamev4.Game, err error)
}

type serverAPI struct {
	gamev4.UnimplementedGameServiceServer
	gameServicer GameServicer
}

func RegisterGrpcHandlers(grpcServer *grpc.Server, gameServicer GameServicer) {
	gamev4.RegisterGameServiceServer(grpcServer, &serverAPI{gameServicer: gameServicer})
}

func (srvApi *serverAPI) AddGame(
	ctx context.Context,
	request *gamev4.AddGameRequest,
) (*gamev4.AddGameResponse, error) {
	if request.Game.Title == "" {
		return &gamev4.AddGameResponse{}, status.Error(codes.InvalidArgument, "Title is required field")
	}
	if request.Game.Description == "" {
		return &gamev4.AddGameResponse{}, status.Error(codes.InvalidArgument, "Description is required field")
	}
	if request.Game.ReleaseYear == nil {
		return &gamev4.AddGameResponse{}, status.Error(codes.InvalidArgument, "Release year is required field")
	}
	return &gamev4.AddGameResponse{}, nil
}

func (srvApi *serverAPI) GetGame(
	ctx context.Context,
	request *gamev4.GetGameRequest,
) (*gamev4.GetGameResponse, error) {
	panic("impl me")
}

func (srvApi *serverAPI) GetTopGames(
	ctx context.Context,
	request *gamev4.GetTopGamesRequest,
) (*gamev4.GetTopGamesResponse, error) {
	panic("impl me")
}
