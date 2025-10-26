package handlers

import (
	"errors"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/outerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UpdateGameStatus(err error) (*game.UpdateGameStatusResponse, error) {
	switch {
	case errors.Is(err, outerror.ErrUnknownGameStatus):
		return &game.UpdateGameStatusResponse{}, status.Error(codes.InvalidArgument, outerror.UnknownGameStatusMessage)
	case errors.Is(err, outerror.ErrInvalidNewGameStatus):
		return &game.UpdateGameStatusResponse{}, status.Error(codes.InvalidArgument, outerror.InvalidNewGameStatusMessage)
	case errors.Is(err, outerror.ErrGameNotFound):
		return &game.UpdateGameStatusResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
	default:
		return &game.UpdateGameStatusResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
}
