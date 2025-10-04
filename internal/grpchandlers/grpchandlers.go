package grpchandlers

import (
	"context"
	"errors"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	gamev2 "github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameServicer interface {
	AddGame(ctx context.Context, game *gamev2.GameRequest) (gameID int64, err error)
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

func (srvAPI *serverAPI) UpdateGameStatus(
	ctx context.Context,
	request *gamev2.UpdateGameStatusRequest,
) (*gamev2.UpdateGameStatusResponse, error) {
	err := srvAPI.gameServicer.UpdateGameStatus(ctx, request.GetGameId(), request.GetNewStatus())
	if err != nil {
		if errors.Is(err, outerror.ErrUnknownGameStatus) {
			return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.NotFound, outerror.UnknownGameStatusMessage)
		} else if errors.Is(err, outerror.ErrInvalidNewGameStatus) {
			return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.InvalidArgument, outerror.InvalidNewGameStatusMessage)
		} else if errors.Is(err, outerror.ErrGameNotFound) {
			return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &gamev2.UpdateGameStatusResponse{}, nil
}
