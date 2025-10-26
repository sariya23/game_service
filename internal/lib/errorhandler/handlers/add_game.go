package handlers

import (
	"errors"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/outerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AddGame(err error, gameID int64) (*game.AddGameResponse, error) {
	switch {
	case errors.Is(err, outerror.ErrGameAlreadyExist):
		return &game.AddGameResponse{}, status.Error(codes.AlreadyExists, outerror.GameAlreadyExistMessage)
	case errors.Is(err, outerror.ErrCannotSaveGameImage):
		return &game.AddGameResponse{GameId: gameID}, nil
	case errors.Is(err, outerror.ErrGenreNotFound):
		return &game.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.GenreNotFoundMessage)
	case errors.Is(err, outerror.ErrTagNotFound):
		return &game.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.TagNotFoundMessage)
	default:
		return &game.AddGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
}
