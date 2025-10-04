package grpchandlers

import (
	"context"

	"github.com/sariya23/game_service/internal/lib/converters"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srvApi *serverAPI) GameList(
	ctx context.Context,
	request *gamev2.GameListRequest,
) (*gamev2.GameListResponse, error) {
	games, err := srvApi.gameServicer.GameList(
		ctx,
		&dto.GameFilters{
			ReleaseYear: request.GetYear(),
			Genres:      request.GetGenres(),
			Tags:        request.GetTags(),
		},
		request.GetLimit(),
	)
	if err != nil {
		return &gamev2.GameListResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	result := make([]*gamev2.GameListResponse_ShortGame, 0, len(games))
	for _, g := range games {
		result = append(result, converters.ToShortGameResponse(g))
	}
	return &gamev2.GameListResponse{Games: result}, nil
}
