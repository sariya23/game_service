package grpchandlers

import (
	"context"

	"github.com/sariya23/game_service/internal/model"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/grpc"
)

type GameFilters struct {
	ReleaseYear string
	Genre       []string
	Tags        []string
}

type GameServicer interface {
	AddGame(ctx context.Context, game model.Game) (createdGame model.Game, err error)
	GetGame(ctx context.Context, gameTitle string) (game model.Game, err error)
	GetTopGames(ctx context.Context, gameFilters GameFilters, limit int32) (games []model.Game, err error)
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
	panic("impl me")
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
