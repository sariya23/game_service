package grpchandlers

import (
	"context"
	"log/slog"

	game_api "github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/lib/converters"
	errorhandler "github.com/sariya23/game_service/internal/lib/errorhandler/handlers"
	"github.com/sariya23/game_service/internal/lib/logger"
	"github.com/sariya23/game_service/internal/outerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srvApi *serverAPI) GetGame(
	ctx context.Context,
	request *game_api.GetGameRequest,
) (*game_api.GetGameResponse, error) {
	log := logger.EnrichRequestID(ctx, srvApi.log)
	log.Info("request to handler",
		slog.String("handler", "GetGame"),
		slog.Any("request", request),
	)
	if request.GameId < 0 {
		log.Warn("invalid request, game_id is negative")
		return &game_api.GetGameResponse{}, status.Error(codes.InvalidArgument, outerror.NegativeGameIDMessage)
	}
	game, err := srvApi.gameServicer.GetGame(ctx, request.GetGameId())
	if err != nil {
		log.Error("internal error")
		return errorhandler.GetGame(err)
	}
	log.Info("success request")
	return &game_api.GetGameResponse{Game: converters.ToProtoGame(game)}, nil
}
