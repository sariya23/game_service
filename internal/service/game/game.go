package gameservice

import (
	"context"
	"io"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	gamev2 "github.com/sariya23/proto_api_games/v5/gen/gamev2"
)

type GameReposetory interface {
	GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (*model.Game, error)
	GetGameByID(ctx context.Context, gameID int64) (*model.Game, error)
	GameList(ctx context.Context, filters dto.GameFilters, limit uint32) ([]model.ShortGame, error)
	SaveGame(ctx context.Context, game model.Game) (int64, error)
	DaleteGame(ctx context.Context, gameID int64) (*dto.DeletedGame, error)
	UpdateGameStatus(ctx context.Context, gameID int64, newStatus gamev2.GameStatusType) error
}

type TagRepository interface {
	GetTagByNames(ctx context.Context, tags []string) ([]model.Tag, error)
}

type GenreRepository interface {
	GetGenreByNames(ctx context.Context, genres []string) ([]model.Genre, error)
}

type S3Storager interface {
	SaveObject(ctx context.Context, name string, data io.Reader) (string, error)
	GetObject(ctx context.Context, name string) (*minio.Object, error)
	DeleteObject(ctx context.Context, name string) error
}

type GameService struct {
	log             *slog.Logger
	gameRepository  GameReposetory
	tagReposetory   TagRepository
	genreReposetory GenreRepository
	s3Storager      S3Storager
}

func NewGameService(
	log *slog.Logger,
	gameReposiroy GameReposetory,
	tagReposetory TagRepository,
	genreReposetory GenreRepository,
	s3Storager S3Storager,

) *GameService {
	return &GameService{
		log:             log,
		s3Storager:      s3Storager,
		tagReposetory:   tagReposetory,
		genreReposetory: genreReposetory,
		gameRepository:  gameReposiroy,
	}
}
