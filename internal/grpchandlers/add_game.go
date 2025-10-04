package grpchandlers

import (
	"context"
	"errors"

	"github.com/sariya23/game_service/internal/lib/validators"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srvApi *serverAPI) AddGame(
	ctx context.Context,
	request *gamev2.AddGameRequest,
) (*gamev2.AddGameResponse, error) {
	if valid, msg := validators.AddGame(request); !valid {
		return &gamev2.AddGameResponse{}, status.Error(codes.InvalidArgument, msg)
	}
	gameID, err := srvApi.gameServicer.AddGame(ctx, request.GetGame())
	if err != nil {
		if errors.Is(err, outerror.ErrGameAlreadyExist) {
			return &gamev2.AddGameResponse{}, status.Error(codes.AlreadyExists, outerror.GameAlreadyExistMessage)
		} else if errors.Is(err, outerror.ErrCannotSaveGameImage) {
			return &gamev2.AddGameResponse{GameId: gameID}, nil
		} else if errors.Is(err, outerror.ErrGenreNotFound) {
			return &gamev2.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.GenreNotFoundMessage)
		} else if errors.Is(err, outerror.ErrTagNotFound) {
			return &gamev2.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.TagNotFoundMessage)
		}
		return &gamev2.AddGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}

	return &gamev2.AddGameResponse{GameId: gameID}, nil
}
