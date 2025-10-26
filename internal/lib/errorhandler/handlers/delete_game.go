package handlers

import (
	"errors"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/outerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func DeleteGame(err error) (*game.DeleteGameResponse, error) {
	if errors.Is(err, outerror.ErrGameNotFound) {
		return &game.DeleteGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
	}
	return &game.DeleteGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
}
