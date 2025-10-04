package grpchandlers

import (
	"context"

	errorhandler "github.com/sariya23/game_service/internal/lib/errorhandler/handlers"
	"github.com/sariya23/game_service/internal/lib/validators"
	"github.com/sariya23/game_service/internal/model/dto"
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
	newGame := &dto.AddGame{
		Title:       request.Game.Title,
		Genres:      request.Game.Genres,
		Description: request.Game.Description,
		ReleaseDate: request.Game.ReleaseDate,
		CoverImage:  request.Game.CoverImage,
		Tags:        request.Game.Tags,
	}
	gameID, err := srvApi.gameServicer.AddGame(ctx, newGame)
	if err != nil {
		return errorhandler.AddGame(err, gameID)
	}

	return &gamev2.AddGameResponse{GameId: gameID}, nil
}
