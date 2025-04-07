package grpchandlers

import (
	"context"
	"errors"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameServicer interface {
	AddGame(ctx context.Context, game *gamev4.GameRequest) (savedGame *gamev4.DomainGame, err error)
	GetGame(ctx context.Context, gameID uint64) (game *gamev4.DomainGame, err error)
	GetTopGames(ctx context.Context, gameFilters model.GameFilters, limit uint32) (games []*gamev4.DomainGame, err error)
	DeleteGame(ctx context.Context, gameID uint64) (deletedGame *gamev4.DomainGame, err error)
}

type serverAPI struct {
	gamev4.UnimplementedGameServiceServer
	gameServicer GameServicer
}

func RegisterGrpcHandlers(grpcServer *grpc.Server, gameServicer GameServicer) {
	gamev4.RegisterGameServiceServer(grpcServer, &serverAPI{gameServicer: gameServicer})
}

func (srvApi *serverAPI) AddGame(
	ctx context.Context,
	request *gamev4.AddGameRequest,
) (*gamev4.AddGameResponse, error) {
	if request.Game.Title == "" {
		return &gamev4.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.TitleRequiredMessage)
	}
	if request.Game.Description == "" {
		return &gamev4.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.DescriptionRequiredMessage)
	}
	if request.Game.ReleaseYear == nil {
		return &gamev4.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.ReleaseYearRequiredMessage)
	}
	savedGame, err := srvApi.gameServicer.AddGame(ctx, request.Game)
	if err != nil {
		if errors.Is(err, outerror.ErrGameAlreadyExist) {
			return &gamev4.AddGameResponse{}, status.Error(codes.AlreadyExists, outerror.GameAlreadyExistMessage)
		} else if errors.Is(err, outerror.ErrCannotSaveGameImage) {
			return &gamev4.AddGameResponse{Game: savedGame}, nil
		}
		return &gamev4.AddGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	response := gamev4.AddGameResponse{
		Game: savedGame,
	}

	return &response, nil
}

func (srvApi *serverAPI) GetGame(
	ctx context.Context,
	request *gamev4.GetGameRequest,
) (*gamev4.GetGameResponse, error) {
	game, err := srvApi.gameServicer.GetGame(ctx, request.GetGameId())
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			return &gamev4.GetGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &gamev4.GetGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &gamev4.GetGameResponse{Game: game}, nil
}

func (srvApi *serverAPI) GetTopGames(
	ctx context.Context,
	request *gamev4.GetTopGamesRequest,
) (*gamev4.GetTopGamesResponse, error) {
	games, err := srvApi.gameServicer.GetTopGames(
		ctx,
		model.GameFilters{
			ReleaseYear: request.GetYear(),
			Genres:      request.GetGenres(),
			Tags:        request.GetTags(),
		},
		request.GetLimit(),
	)
	if err != nil {
		return &gamev4.GetTopGamesResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &gamev4.GetTopGamesResponse{Games: games}, nil
}

func (srvAPI *serverAPI) DeleteGame(
	ctx context.Context,
	request *gamev4.DeleteGameRequest,
) (*gamev4.DeleteGameResponse, error) {
	game, err := srvAPI.gameServicer.DeleteGame(ctx, request.GetGameId())
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			return &gamev4.DeleteGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &gamev4.DeleteGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &gamev4.DeleteGameResponse{Game: game}, nil
}
