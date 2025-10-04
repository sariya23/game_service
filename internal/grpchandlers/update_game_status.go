package grpchandlers

import (
	"context"
	"errors"

	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
