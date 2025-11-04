package grpchandlers

import (
	"context"
	"log/slog"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/lib/converters"
	errorhandler "github.com/sariya23/game_service/internal/lib/errorhandler/handlers"
	"github.com/sariya23/game_service/internal/lib/validators"
	"github.com/sariya23/game_service/internal/model/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srvApi *serverAPI) AddGame(
	ctx context.Context,
	request *game.AddGameRequest,
) (*game.AddGameResponse, error) {
	srvApi.log.Info("request to handler", slog.String("handler", "AddGame"), slog.Any("request", request))
	if valid, msg := validators.AddGame(request); !valid {
		return &game.AddGameResponse{}, status.Error(codes.InvalidArgument, msg)
	}
	newGame := dto.AddGameHandler{
		Title:       request.Game.Title,
		Genres:      request.Game.Genres,
		Description: request.Game.Description,
		ReleaseDate: converters.FromProtoDate(request.Game.ReleaseDate),
		CoverImage:  request.Game.CoverImage,
		Tags:        request.Game.Tags,
	}
	gameID, err := srvApi.gameServicer.AddGame(ctx, newGame)
	if err != nil {
		return errorhandler.AddGame(err, gameID)
	}

	return &game.AddGameResponse{GameId: gameID}, nil
}
