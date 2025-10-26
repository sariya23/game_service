package grpchandlers

import (
	"context"
	"errors"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/outerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srvAPI *serverAPI) DeleteGame(
	ctx context.Context,
	request *game.DeleteGameRequest,
) (*game.DeleteGameResponse, error) {
	gameID, err := srvAPI.gameServicer.DeleteGame(ctx, request.GetGameId())
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			return &game.DeleteGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &game.DeleteGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &game.DeleteGameResponse{GameId: gameID}, nil
}
