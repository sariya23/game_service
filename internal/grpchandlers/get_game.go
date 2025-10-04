package grpchandlers

import (
	"context"

	"github.com/sariya23/game_service/internal/lib/converters"
	errorhandler "github.com/sariya23/game_service/internal/lib/errorhandler/handlers"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

func (srvApi *serverAPI) GetGame(
	ctx context.Context,
	request *gamev2.GetGameRequest,
) (*gamev2.GetGameResponse, error) {
	game, err := srvApi.gameServicer.GetGame(ctx, request.GetGameId())
	if err != nil {
		return errorhandler.GetGame(err)
	}
	return &gamev2.GetGameResponse{Game: converters.ToProtoGame(*game)}, nil
}
