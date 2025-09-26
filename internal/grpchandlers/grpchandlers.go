package grpchandlers

import (
	"context"
	"errors"

	"github.com/sariya23/game_service/internal/converters"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameServicer interface {
	AddGame(ctx context.Context, game *gamev4.GameRequest) (gameID uint64, err error)
	GetGame(ctx context.Context, gameID uint64) (game *model.Game, err error)
	GetTopGames(ctx context.Context, gameFilters dto.GameFilters, limit uint32) (games []model.ShortGame, err error)
	DeleteGame(ctx context.Context, gameID uint64) (deletedGameID uint64, err error)
	UpdateGameStatus(ctx context.Context, gameID uint64, newStatus gamev4.GameStatusType) error
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
	if request.Game.ReleaseDate == nil {
		return &gamev4.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.ReleaseYearRequiredMessage)
	}
	gameID, err := srvApi.gameServicer.AddGame(ctx, request.GetGame())
	if err != nil {
		if errors.Is(err, outerror.ErrGameAlreadyExist) {
			return &gamev4.AddGameResponse{}, status.Error(codes.AlreadyExists, outerror.GameAlreadyExistMessage)
		} else if errors.Is(err, outerror.ErrCannotSaveGameImage) {
			return &gamev4.AddGameResponse{GameId: gameID}, nil
		} else if errors.Is(err, outerror.ErrGenreNotFound) {
			return &gamev4.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.GenreNotFoundMessage)
		} else if errors.Is(err, outerror.ErrTagNotFound) {
			return &gamev4.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.TagNotFoundMessage)
		}
		return &gamev4.AddGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}

	return &gamev4.AddGameResponse{GameId: gameID}, nil
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
	return &gamev4.GetGameResponse{Game: converters.ToProtoGame(*game)}, nil
}

func (srvApi *serverAPI) GetTopGames(
	ctx context.Context,
	request *gamev4.GetTopGamesRequest,
) (*gamev4.GetTopGamesResponse, error) {
	games, err := srvApi.gameServicer.GetTopGames(
		ctx,
		dto.GameFilters{
			ReleaseYear: request.GetYear(),
			Genres:      request.GetGenres(),
			Tags:        request.GetTags(),
		},
		request.GetLimit(),
	)
	if err != nil {
		return &gamev4.GetTopGamesResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	result := make([]*gamev4.GetTopGamesResponse_ShortGame, 0, len(games))
	for _, g := range games {
		result = append(result, converters.ToShortGameResponse(g))
	}
	return &gamev4.GetTopGamesResponse{Games: result}, nil
}

func (srvAPI *serverAPI) DeleteGame(
	ctx context.Context,
	request *gamev4.DeleteGameRequest,
) (*gamev4.DeleteGameResponse, error) {
	gameID, err := srvAPI.gameServicer.DeleteGame(ctx, request.GetGameId())
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			return &gamev4.DeleteGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &gamev4.DeleteGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &gamev4.DeleteGameResponse{GameId: gameID}, nil
}

func (srvAPI *serverAPI) UpdateGameStatus(
	ctx context.Context,
	request *gamev4.UpdateGameStatusRequest,
) (*gamev4.UpdateGameStatusReponse, error) {
	_, err := srvAPI.gameServicer.GetGame(ctx, request.GetGameId())
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			return &gamev4.UpdateGameStatusReponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &gamev4.UpdateGameStatusReponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	err = srvAPI.gameServicer.UpdateGameStatus(ctx, request.GetGameId(), request.GetNewStautus())
	if err != nil {
		if errors.Is(err, outerror.ErrUnknownGameStatus) {
			return &gamev4.UpdateGameStatusReponse{}, status.Error(codes.InvalidArgument, outerror.UnknownGameStatusMessage)
		}
		return &gamev4.UpdateGameStatusReponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &gamev4.UpdateGameStatusReponse{}, nil
}
