package grpchandlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/outerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srvApi *serverAPI) DeleteGame(
	ctx context.Context,
	request *game.DeleteGameRequest,
) (*game.DeleteGameResponse, error) {
	requestID := ctx.Value("request_id").(string)
	log := srvApi.log.With("request_id", requestID)
	log.Info("request to handler",
		slog.String("handler", "AddGame"),
		slog.Any("request", request),
	)
	log.Info("request to handler", slog.String("handler", "DeleteGame"), slog.Any("request", request))
	gameID, err := srvApi.gameServicer.DeleteGame(ctx, request.GetGameId())
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			return &game.DeleteGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &game.DeleteGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	log.Info("game deleted successfully")
	return &game.DeleteGameResponse{GameId: gameID}, nil
}
