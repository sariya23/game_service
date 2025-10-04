package handlers

import (
	"errors"

	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UpdateGameStatus(err error) (*gamev2.UpdateGameStatusResponse, error) {
	switch {
	case errors.Is(err, outerror.ErrUnknownGameStatus):
		return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.InvalidArgument, outerror.UnknownGameStatusMessage)
	case errors.Is(err, outerror.ErrInvalidNewGameStatus):
		return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.InvalidArgument, outerror.InvalidNewGameStatusMessage)
	case errors.Is(err, outerror.ErrGameNotFound):
		return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
	default:
		return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
}
