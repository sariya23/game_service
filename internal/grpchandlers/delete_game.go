package grpchandlers

import (
	"context"
	"errors"

	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srvAPI *serverAPI) DeleteGame(
	ctx context.Context,
	request *gamev2.DeleteGameRequest,
) (*gamev2.DeleteGameResponse, error) {
	gameID, err := srvAPI.gameServicer.DeleteGame(ctx, request.GetGameId())
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			return &gamev2.DeleteGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &gamev2.DeleteGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &gamev2.DeleteGameResponse{GameId: gameID}, nil
}
