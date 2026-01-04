package grpchandlers

import (
	"context"
	"log/slog"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/lib/converters"
	"github.com/sariya23/game_service/internal/lib/logger"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srvApi *serverAPI) GameList(
	ctx context.Context,
	request *game.GameListRequest,
) (*game.GameListResponse, error) {
	log := logger.EnrichRequestID(ctx, srvApi.log)
	log.Info("request to handler",
		slog.String("handler", "GameList"),
		slog.Any("request", request),
	)
	if request.Year < 0 {
		log.Warn("negative year")
		return &game.GameListResponse{}, status.Error(codes.InvalidArgument, outerror.NegativeYearMessage)
	}
	games, err := srvApi.gameServicer.GameList(
		ctx,
		dto.GameFilters{
			ReleaseYear: request.GetYear(),
			Genres:      request.GetGenres(),
			Tags:        request.GetTags(),
		},
		request.GetLimit(),
	)
	if err != nil {
		return &game.GameListResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	result := make([]*game.GameListResponse_ShortGame, 0, len(games))
	for _, g := range games {
		result = append(result, converters.ToShortGameResponse(g))
	}
	log.Info("success get games")
	return &game.GameListResponse{Games: result}, nil
}
