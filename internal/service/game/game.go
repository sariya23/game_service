package gameservice

import (
	"context"
	"io"
	"log/slog"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
)

type GameReposetory interface {
	GetGameByTitleAndReleaseYear(ctx context.Context, title string, releaseYear int32) (*model.Game, error)
	GetGameByID(ctx context.Context, gameID int64) (*model.Game, error)
	GameList(ctx context.Context, filters dto.GameFilters, limit uint32) ([]model.ShortGame, error)
	SaveGame(ctx context.Context, game dto.AddGameService) (int64, error)
	DaleteGame(ctx context.Context, gameID int64) (*dto.DeletedGame, error)
	UpdateGameStatus(ctx context.Context, gameID int64, newStatus game.GameStatusType) error
}

type TagRepository interface {
	GetTagByNames(ctx context.Context, tags []string) ([]model.Tag, error)
}

type GenreRepository interface {
	GetGenreByNames(ctx context.Context, genres []string) ([]model.Genre, error)
}

type S3Storager interface {
	SaveObject(ctx context.Context, name string, data io.Reader) (string, error)
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
