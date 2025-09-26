package gameservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const GameNotSaveID uint64 = 0

func TestAddGame(t *testing.T) {
	t.Run("Нельзя добавить игру, так как она уже есть в БД", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		expectedError := outerror.ErrGameAlreadyExist
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Genres = nil
		gameToAdd.Tags = nil
		gameToAdd.CoverImage = nil
		game := &model.Game{
			Title:       gameToAdd.GetTitle(),
			Description: gameToAdd.GetDescription(),
			ReleaseDate: time.Date(
				int(gameToAdd.ReleaseDate.Year),
				time.Month(gameToAdd.ReleaseDate.Month),
				int(gameToAdd.ReleaseDate.Day),
				0,
				0,
				0,
				0,
				time.UTC),
		}
		gameMockRepo.On("GetGameByTitleAndReleaseYear", mock.Anything, gameToAdd.Title, gameToAdd.GetReleaseDate().Year).Return(game, nil).Once()

		gameID, err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
	t.Run("Игра не создается с несуществующими тегами", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		expectedError := outerror.ErrTagNotFound
		gameToAdd := random.RandomAddGameRequest()
		tags := gameToAdd.Tags
		gameToAdd.CoverImage = nil
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		tagMockRepo.On("GetTagByNames", mock.Anything, tags).Return(nil, outerror.ErrTagNotFound).Once()
		gameID, err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
	t.Run("Игра не создается с несуществующими жанрами", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		expectedError := outerror.ErrGenreNotFound
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.CoverImage = nil
		gameToAdd.Tags = nil
		genres := gameToAdd.Genres
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		genreMockRepo.On("GetGenreByNames", mock.Anything, genres).Return(nil, outerror.ErrGenreNotFound).Once()
		gameID, err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
	t.Run("Не удалось сохранить игру", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Genres = nil
		gameToAdd.Tags = nil
		gameToAdd.CoverImage = nil
		game := model.Game{
			Title:       gameToAdd.GetTitle(),
			Description: gameToAdd.GetDescription(),
			ReleaseDate: time.Date(
				int(gameToAdd.ReleaseDate.Year),
				time.Month(gameToAdd.ReleaseDate.Month),
				int(gameToAdd.ReleaseDate.Day),
				0,
				0,
				0,
				0,
				time.UTC),
		}
		expectedError := errors.New("some error")
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(GameNotSaveID, expectedError).Once()
		gameID, err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
	t.Run("Игра сохраняется даже в случае не сохранения обложки в S3", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Genres = nil
		gameToAdd.Tags = nil
		game := model.Game{
			Title:       gameToAdd.GetTitle(),
			Description: gameToAdd.GetDescription(),
			ReleaseDate: time.Date(
				int(gameToAdd.ReleaseDate.Year),
				time.Month(gameToAdd.ReleaseDate.Month),
				int(gameToAdd.ReleaseDate.Day),
				0,
				0,
				0,
				0,
				time.UTC),
		}
		savedGameID := uint64(23)
		expectedError := outerror.ErrCannotSaveGameImage
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return("", expectedError).Once()
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(savedGameID, nil).Once()
		gameMockRepo.On("UpdateGameStatus", mock.Anything, savedGameID, gamev4.GameStatusType_PENDING).Return(nil).Once()
		gameID, err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, savedGameID, gameID)
	})
	t.Run("Сохранение игры без ошибок", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Genres = nil
		gameToAdd.Tags = nil
		game := model.Game{
			Title:       gameToAdd.GetTitle(),
			Description: gameToAdd.GetDescription(),
			ReleaseDate: time.Date(
				int(gameToAdd.ReleaseDate.Year),
				time.Month(gameToAdd.ReleaseDate.Month),
				int(gameToAdd.ReleaseDate.Day),
				0,
				0,
				0,
				0,
				time.UTC),
			ImageURL: string(gameToAdd.CoverImage),
		}
		expectedGameID := uint64(228)
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return(string(gameToAdd.CoverImage), nil).Once()
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(expectedGameID, nil).Once()
		gameMockRepo.On("UpdateGameStatus", mock.Anything, expectedGameID, gamev4.GameStatusType_PENDING).Return(nil).Once()
		gameID, err := gameService.AddGame(context.Background(), gameToAdd)
		require.NoError(t, err)
		require.Equal(t, expectedGameID, gameID)
	})
	t.Run("Сохранение игры с тэгами и жанрами без ошибок", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Tags = []string{"TAG"}
		gameToAdd.Genres = []string{"GENRE"}
		modelTags := []model.Tag{{TagID: 1, TagName: "TAG"}}
		modelGenres := []model.Genre{{GenreID: 1, GenreName: "GENRE"}}
		game := model.Game{
			Title:       gameToAdd.Title,
			Description: gameToAdd.Description,
			ReleaseDate: time.Date(
				int(gameToAdd.ReleaseDate.Year),
				time.Month(gameToAdd.ReleaseDate.Month),
				int(gameToAdd.ReleaseDate.Day),
				0,
				0,
				0,
				0,
				time.UTC,
			),
			ImageURL: string(gameToAdd.CoverImage),
			Tags:     modelTags,
			Genres:   modelGenres,
		}
		expectedGameID := uint64(288)
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return(game.ImageURL, nil).Once()
		tagMockRepo.On("GetTagByNames", mock.Anything, gameToAdd.GetTags()).Return(modelTags, nil)
		genreMockRepo.On("GetGenreByNames", mock.Anything, gameToAdd.GetGenres()).Return(modelGenres, nil)
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(expectedGameID, nil).Once()
		gameMockRepo.On("UpdateGameStatus", mock.Anything, expectedGameID, gamev4.GameStatusType_PENDING).Return(nil).Once()
		gameID, err := gameService.AddGame(context.Background(), gameToAdd)
		require.NoError(t, err)
		require.Equal(t, expectedGameID, gameID)
	})
	t.Run("Игра не сохраняется при ошибке получении тэгов", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameToAdd := random.RandomAddGameRequest()
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return("qwe", nil).Once()
		expectedError := errors.New("some err")
		tagMockRepo.On("GetTagByNames", mock.Anything, gameToAdd.GetTags()).Return(nil, expectedError).Once()
		gameID, err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
	t.Run("Игра не сохранилась при при ошибке получения жанров", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock)
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Tags = nil
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return("qwe", nil).Once()
		expectedError := errors.New("some err")
		genreMockRepo.On("GetGenreByNames", mock.Anything, gameToAdd.GetGenres()).Return(nil, expectedError).Once()
		gameID, err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
}
