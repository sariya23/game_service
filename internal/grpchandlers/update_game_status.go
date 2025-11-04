package grpchandlers

import (
	"context"
	"log/slog"

	"github.com/sariya23/api_game_service/gen/game"
	errorhandler "github.com/sariya23/game_service/internal/lib/errorhandler/handlers"
	"github.com/sariya23/game_service/internal/outerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srvApi *serverAPI) UpdateGameStatus(
	ctx context.Context,
	request *game.UpdateGameStatusRequest,
) (*game.UpdateGameStatusResponse, error) {
	srvApi.log.Info("request to handler", slog.String("handler", "UpdateGameStatus"), slog.Any("request", request))
	if request.GameId < 0 {
		return nil, status.Error(codes.InvalidArgument, outerror.NegativeGameIDMessage)
	}
	err := srvApi.gameServicer.UpdateGameStatus(ctx, request.GetGameId(), request.GetNewStatus())
	if err != nil {
		return errorhandler.UpdateGameStatus(err)
	}
	return &game.UpdateGameStatusResponse{}, nil
}
