package errorhandler

import (
	"errors"

	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AddGame(err error, gameID int64) (*gamev2.AddGameResponse, error) {
	switch {
	case errors.Is(err, outerror.ErrGameAlreadyExist):
		return &gamev2.AddGameResponse{}, status.Error(codes.AlreadyExists, outerror.GameAlreadyExistMessage)
	case errors.Is(err, outerror.ErrCannotSaveGameImage):
		return &gamev2.AddGameResponse{GameId: gameID}, nil
	case errors.Is(err, outerror.ErrGenreNotFound):
		return &gamev2.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.GenreNotFoundMessage)
	case errors.Is(err, outerror.ErrTagNotFound):
		return &gamev2.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.TagNotFoundMessage)
	default:
		return &gamev2.AddGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
}
