//go:build unit
// +build unit

package gameservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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
		suite := NewSuite()
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
		suite.gameMockRepo.On("GetGameByTitleAndReleaseYear", mock.Anything, gameToAdd.Title, gameToAdd.GetReleaseDate().Year).Return(game, nil).Once()

		gameID, err := suite.gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
	t.Run("Игра не создается с несуществующими тегами", func(t *testing.T) {
		suite := NewSuite()
		expectedError := outerror.ErrTagNotFound
		gameToAdd := random.RandomAddGameRequest()
		tags := gameToAdd.Tags
		gameToAdd.CoverImage = nil
		suite.gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		suite.tagMockRepo.On("GetTagByNames", mock.Anything, tags).Return(nil, outerror.ErrTagNotFound).Once()
		gameID, err := suite.gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
	t.Run("Игра не создается с несуществующими жанрами", func(t *testing.T) {
		suite := NewSuite()
		expectedError := outerror.ErrGenreNotFound
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.CoverImage = nil
		gameToAdd.Tags = nil
		genres := gameToAdd.Genres
		suite.gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		suite.genreMockRepo.On("GetGenreByNames", mock.Anything, genres).Return(nil, outerror.ErrGenreNotFound).Once()
		gameID, err := suite.gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
	t.Run("Не удалось сохранить игру", func(t *testing.T) {
		suite := NewSuite()
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
		suite.gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		suite.gameMockRepo.On("SaveGame", mock.Anything, game).Return(GameNotSaveID, expectedError).Once()
		gameID, err := suite.gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
	t.Run("Игра сохраняется даже в случае не сохранения обложки в S3", func(t *testing.T) {
		suite := NewSuite()
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
		suite.gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		suite.s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return("", expectedError).Once()
		suite.gameMockRepo.On("SaveGame", mock.Anything, game).Return(savedGameID, nil).Once()
		suite.gameMockRepo.On("UpdateGameStatus", mock.Anything, savedGameID, gamev4.GameStatusType_PENDING).Return(nil).Once()
		gameID, err := suite.gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Equal(t, savedGameID, gameID)
	})
	t.Run("Сохранение игры без ошибок", func(t *testing.T) {
		suite := NewSuite()
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
		suite.gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		suite.s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return(string(gameToAdd.CoverImage), nil).Once()
		suite.gameMockRepo.On("SaveGame", mock.Anything, game).Return(expectedGameID, nil).Once()
		suite.gameMockRepo.On("UpdateGameStatus", mock.Anything, expectedGameID, gamev4.GameStatusType_PENDING).Return(nil).Once()
		gameID, err := suite.gameService.AddGame(context.Background(), gameToAdd)
		require.NoError(t, err)
		require.Equal(t, expectedGameID, gameID)
	})
	t.Run("Сохранение игры с тэгами и жанрами без ошибок", func(t *testing.T) {
		suite := NewSuite()
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
		suite.gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		suite.s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return(game.ImageURL, nil).Once()
		suite.tagMockRepo.On("GetTagByNames", mock.Anything, gameToAdd.GetTags()).Return(modelTags, nil)
		suite.genreMockRepo.On("GetGenreByNames", mock.Anything, gameToAdd.GetGenres()).Return(modelGenres, nil)
		suite.gameMockRepo.On("SaveGame", mock.Anything, game).Return(expectedGameID, nil).Once()
		suite.gameMockRepo.On("UpdateGameStatus", mock.Anything, expectedGameID, gamev4.GameStatusType_PENDING).Return(nil).Once()
		gameID, err := suite.gameService.AddGame(context.Background(), gameToAdd)
		require.NoError(t, err)
		require.Equal(t, expectedGameID, gameID)
	})
	t.Run("Игра не сохраняется при ошибке получении тэгов", func(t *testing.T) {
		suite := NewSuite()
		gameToAdd := random.RandomAddGameRequest()
		suite.gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		suite.s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return("qwe", nil).Once()
		expectedError := errors.New("some err")
		suite.tagMockRepo.On("GetTagByNames", mock.Anything, gameToAdd.GetTags()).Return(nil, expectedError).Once()
		gameID, err := suite.gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
	t.Run("Игра не сохранилась при при ошибке получения жанров", func(t *testing.T) {
		suite := NewSuite()
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Tags = nil
		suite.gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		suite.s3Mock.On(
			"SaveObject",
			mock.Anything,
			fmt.Sprintf("%s_%d", gameToAdd.GetTitle(), int(gameToAdd.GetReleaseDate().Year)),
			bytes.NewReader(gameToAdd.GetCoverImage()),
		).Return("qwe", nil).Once()
		expectedError := errors.New("some err")
		suite.genreMockRepo.On("GetGenreByNames", mock.Anything, gameToAdd.GetGenres()).Return(nil, expectedError).Once()
		gameID, err := suite.gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
		require.Zero(t, gameID)
	})
}
