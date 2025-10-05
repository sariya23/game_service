package handlers

import (
	"errors"

	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func DeleteGame(err error) (*gamev2.DeleteGameResponse, error) {
	if errors.Is(err, outerror.ErrGameNotFound) {
		return &gamev2.DeleteGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
	}
	return &gamev2.DeleteGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
}
