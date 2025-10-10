package grpchandlers

import (
	"context"

	errorhandler "github.com/sariya23/game_service/internal/lib/errorhandler/handlers"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srvAPI *serverAPI) UpdateGameStatus(
	ctx context.Context,
	request *gamev2.UpdateGameStatusRequest,
) (*gamev2.UpdateGameStatusResponse, error) {
	if request.GameId < 0 {
		return nil, status.Error(codes.InvalidArgument, outerror.NegativeGameIDMessage)
	}
	err := srvAPI.gameServicer.UpdateGameStatus(ctx, request.GetGameId(), request.GetNewStatus())
	if err != nil {
		return errorhandler.UpdateGameStatus(err)
	}
	return &gamev2.UpdateGameStatusResponse{}, nil
}
