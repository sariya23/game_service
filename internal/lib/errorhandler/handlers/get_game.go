package handlers

import (
	"errors"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/outerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetGame(err error) (*game.GetGameResponse, error) {
	if errors.Is(err, outerror.ErrGameNotFound) {
		return &game.GetGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
	}
	return &game.GetGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
}
