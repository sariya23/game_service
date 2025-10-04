package grpchandlers

import (
	"context"
	"errors"

	"github.com/sariya23/game_service/internal/converters"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/outerror"
	gamev2 "github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameServicer interface {
	AddGame(ctx context.Context, game *gamev2.GameRequest) (gameID int64, err error)
	GetGame(ctx context.Context, gameID int64) (game *model.Game, err error)
	GameList(ctx context.Context, gameFilters dto.GameFilters, limit uint32) (games []model.ShortGame, err error)
	DeleteGame(ctx context.Context, gameID int64) (deletedGameID int64, err error)
	UpdateGameStatus(ctx context.Context, gameID int64, newStatus gamev2.GameStatusType) error
}

type serverAPI struct {
	gamev2.UnimplementedGameServiceServer
	gameServicer GameServicer
}

func RegisterGrpcHandlers(grpcServer *grpc.Server, gameServicer GameServicer) {
	gamev2.RegisterGameServiceServer(grpcServer, &serverAPI{gameServicer: gameServicer})
}

func (srvApi *serverAPI) AddGame(
	ctx context.Context,
	request *gamev2.AddGameRequest,
) (*gamev2.AddGameResponse, error) {
	if request.Game.Title == "" {
		return &gamev2.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.TitleRequiredMessage)
	}
	if request.Game.Description == "" {
		return &gamev2.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.DescriptionRequiredMessage)
	}
	if request.Game.ReleaseDate == nil {
		return &gamev2.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.ReleaseYearRequiredMessage)
	}
	gameID, err := srvApi.gameServicer.AddGame(ctx, request.GetGame())
	if err != nil {
		if errors.Is(err, outerror.ErrGameAlreadyExist) {
			return &gamev2.AddGameResponse{}, status.Error(codes.AlreadyExists, outerror.GameAlreadyExistMessage)
		} else if errors.Is(err, outerror.ErrCannotSaveGameImage) {
			return &gamev2.AddGameResponse{GameId: gameID}, nil
		} else if errors.Is(err, outerror.ErrGenreNotFound) {
			return &gamev2.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.GenreNotFoundMessage)
		} else if errors.Is(err, outerror.ErrTagNotFound) {
			return &gamev2.AddGameResponse{}, status.Error(codes.InvalidArgument, outerror.TagNotFoundMessage)
		}
		return &gamev2.AddGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}

	return &gamev2.AddGameResponse{GameId: gameID}, nil
}

func (srvApi *serverAPI) GetGame(
	ctx context.Context,
	request *gamev2.GetGameRequest,
) (*gamev2.GetGameResponse, error) {
	game, err := srvApi.gameServicer.GetGame(ctx, request.GetGameId())
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			return &gamev2.GetGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &gamev2.GetGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &gamev2.GetGameResponse{Game: converters.ToProtoGame(*game)}, nil
}

func (srvApi *serverAPI) GameList(
	ctx context.Context,
	request *gamev2.GameListRequest,
) (*gamev2.GameListResponse, error) {
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
		return &gamev2.GameListResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	result := make([]*gamev2.GameListResponse_ShortGame, 0, len(games))
	for _, g := range games {
		result = append(result, converters.ToShortGameResponse(g))
	}
	return &gamev2.GameListResponse{Games: result}, nil
}

func (srvAPI *serverAPI) DeleteGame(
	ctx context.Context,
	request *gamev2.DeleteGameRequest,
) (*gamev2.DeleteGameResponse, error) {
	gameID, err := srvAPI.gameServicer.DeleteGame(ctx, request.GetGameId())
	if err != nil {
		if errors.Is(err, outerror.ErrGameNotFound) {
			return &gamev2.DeleteGameResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &gamev2.DeleteGameResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &gamev2.DeleteGameResponse{GameId: gameID}, nil
}

func (srvAPI *serverAPI) UpdateGameStatus(
	ctx context.Context,
	request *gamev2.UpdateGameStatusRequest,
) (*gamev2.UpdateGameStatusResponse, error) {
	err := srvAPI.gameServicer.UpdateGameStatus(ctx, request.GetGameId(), request.GetNewStatus())
	if err != nil {
		if errors.Is(err, outerror.ErrUnknownGameStatus) {
			return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.NotFound, outerror.UnknownGameStatusMessage)
		} else if errors.Is(err, outerror.ErrInvalidNewGameStatus) {
			return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.InvalidArgument, outerror.InvalidNewGameStatusMessage)
		} else if errors.Is(err, outerror.ErrGameNotFound) {
			return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.NotFound, outerror.GameNotFoundMessage)
		}
		return &gamev2.UpdateGameStatusResponse{}, status.Error(codes.Internal, outerror.InternalMessage)
	}
	return &gamev2.UpdateGameStatusResponse{}, nil
}
