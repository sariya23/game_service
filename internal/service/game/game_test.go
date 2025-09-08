package gameservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/date"
)

type mockEmailAlerter struct {
	mock.Mock
}

func (m *mockEmailAlerter) SendMessage(subject string, body string) error {
	args := m.Called(subject, body)

	return args.Error(0)
}

type mockTagRepository struct {
	mock.Mock
}

func (m *mockTagRepository) GetTags(ctx context.Context, tags []string) ([]model.Tag, error) {
	args := m.Called(ctx, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Tag), args.Error(1)
}

type mockGenreRepository struct {
	mock.Mock
}

func (m *mockGenreRepository) GetGenres(ctx context.Context, genres []string) ([]model.Genre, error) {
	args := m.Called(ctx, genres)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Genre), args.Error(1)
}

func TestAddGame(t *testing.T) {
	t.Run("Нельзя добавить игру, так как она уже есть в БД", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		expectedError := outerror.ErrGameAlreadyExist
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		game := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
		}
		gameMockRepo.On("GetGameByTitleAndReleaseYear", mock.Anything, gameToAdd.Title, gameToAdd.GetReleaseDate().Year).Return(game, nil).Once()

		err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
	})
	t.Run("Игра не создается с несуществующими тегами", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		expectedError := outerror.ErrTagNotFound
		tags := []string{"test"}
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			Tags:        tags,
		}
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		tagMockRepo.On("GetTags", mock.Anything, tags).Return(nil, outerror.ErrTagNotFound).Once()
		err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
	})
	t.Run("Игра не создается с несуществующими жанрами", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		expectedError := outerror.ErrGenreNotFound
		genres := []string{"test"}
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			Genres:      genres,
		}
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		genreMockRepo.On("GetGenres", mock.Anything, genres).Return(nil, outerror.ErrGenreNotFound).Once()
		err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedError)
	})
	t.Run("Не удалось сохранить игру", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
		}
		game := model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
		}
		expectedErr := errors.New("some error")
		gameMockRepo.On(
			"GetGameByTitleAndReleaseYear",
			mock.Anything,
			gameToAdd.GetTitle(),
			gameToAdd.GetReleaseDate().Year,
		).Return(nil, outerror.ErrGameNotFound).Once()
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(expectedErr).Once()
		err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("Игра сохраняется даже в случае не сохранения обложки в S3", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
		}
		game := model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
		}
		expectedErr := outerror.ErrCannotSaveGameImage
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
		).Return("", expectedErr).Once()
		mailerMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil).Once()
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(nil).Once()
		err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("Сохранение игры без ошибок", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
		}
		game := model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			ImageURL:    "qwe",
		}
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
		mailerMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil).Once()
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(nil).Once()
		err := gameService.AddGame(context.Background(), gameToAdd)
		require.NoError(t, err)
	})
	t.Run("Сохранение игры с тэгами и жанрами без ошибок", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"TAG"},
			Genres:      []string{"GENRE"},
		}
		game := model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			ImageURL:    "qwe",
			Tags:        []model.Tag{{1, "TAG"}},
			Genres:      []model.Genre{{1, "GENRE"}},
		}
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
		tagMockRepo.On("GetTags", mock.Anything, gameToAdd.GetTags()).Return([]model.Tag{{1, "TAG"}}, nil)
		genreMockRepo.On("GetGenres", mock.Anything, gameToAdd.GetGenres()).Return([]model.Genre{{1, "GENRE"}}, nil)
		mailerMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil).Once()
		gameMockRepo.On("SaveGame", mock.Anything, game).Return(nil).Once()
		err := gameService.AddGame(context.Background(), gameToAdd)
		require.NoError(t, err)
	})
	t.Run("Игра не сохраняется при ошибке получении тэгов", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		// TODO: сделать через gofakeit
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 34",
			Description: "test",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Tags:        []string{"ABOBA"},
		}
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
		expectedErr := errors.New("some err")
		tagMockRepo.On("GetTags", mock.Anything, gameToAdd.GetTags()).Return(nil, expectedErr).Once()
		err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("Игра не сохранилась при при ошибке получения жанров", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameToAdd := &gamev4.GameRequest{
			Title:       "Dark Souls 34",
			Description: "fgudgdifg",
			ReleaseDate: &date.Date{Year: 2016, Month: 3, Day: 16},
			CoverImage:  []byte("qwe"),
			Genres:      []string{"ABOBA"},
		}
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
		expectedErr := errors.New("some err")
		tagMockRepo.AssertNotCalled(t, "GetTags")
		genreMockRepo.On("GetGenres", mock.Anything, gameToAdd.GetGenres()).Return(nil, expectedErr).Once()
		err := gameService.AddGame(context.Background(), gameToAdd)
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestGetGame(t *testing.T) {
	t.Run("Игра не найдена", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameID := uint64(1)
		expectedError := outerror.ErrGameNotFound
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, expectedError).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, game)
	})
	t.Run("Internal ошибка", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameID := uint64(1)
		expectedError := errors.New("some error")
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(nil, expectedError).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.ErrorIs(t, err, expectedError)
		require.Nil(t, game)
	})
	t.Run("Успешное получение игры", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameID := uint64(1)
		expectedGame := &model.Game{
			Title:       "Dark Souls 3",
			Genres:      []model.Genre{{GenreID: 1, GenreName: "Hard"}},
			Description: "qwe",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			Tags:        []model.Tag{{TagID: 1, TagName: "Aboba"}},
			GameID:      2,
		}
		gameMockRepo.On("GetGameByID", mock.Anything, gameID).Return(expectedGame, nil).Once()
		game, err := gameService.GetGame(context.Background(), gameID)
		require.NoError(t, err)
		require.Equal(t, expectedGame, game)
	})
}

func TestDeleteGame(t *testing.T) {
	t.Run("Успешное удаление игры", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameID := uint64(4)
		deletedGame := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			ImageURL:    "qwe",
		}
		gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseDate.Year()))
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(deletedGame, nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(nil).Once()

		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.NoError(t, err)
		require.Equal(t, deletedGame, game)
	})
	t.Run("Нет игры для удаления", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameID := uint64(4)
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(nil, outerror.ErrGameNotFound).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.ErrorIs(t, err, outerror.ErrGameNotFound)
		require.Nil(t, game)
	})
	t.Run("Неожиданная ошибка при удалении игры", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameID := uint64(4)
		someErr := errors.New("some err")
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(nil, someErr).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.ErrorIs(t, err, someErr)
		require.Nil(t, game)
	})
	t.Run("У игры нет обложки в S3", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameID := uint64(4)
		deletedGame := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			ImageURL:    "qwe",
		}
		gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseDate.Year()))
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(deletedGame, nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(outerror.ErrImageNotFoundS3).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.Equal(t, deletedGame, game)
		require.ErrorIs(t, err, outerror.ErrImageNotFoundS3)
	})
	t.Run("Не удалось удалить обложку из S3", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		gameID := uint64(4)
		deletedGame := &model.Game{
			Title:       "Dark Souls 3",
			Description: "test",
			ReleaseDate: time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC),
			ImageURL:    "qwe",
		}
		someErr := errors.New("some error")
		gameKey := minioclient.GameKey(deletedGame.Title, int(deletedGame.ReleaseDate.Year()))
		gameMockRepo.On("DaleteGame", mock.Anything, gameID).Return(deletedGame, nil).Once()
		s3Mock.On("DeleteObject", mock.Anything, gameKey).Return(someErr).Once()
		game, err := gameService.DeleteGame(context.Background(), gameID)
		require.Equal(t, deletedGame, game)
		require.ErrorIs(t, err, someErr)
	})
}

func TestGetTopGames(t *testing.T) {
	t.Run("Internal ошибка", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		filters := model.GameFilters{ReleaseYear: 2020}
		gameMockRepo.On("GetTopGames", mock.Anything, filters, uint32(10)).Return(([]model.Game)(nil), errors.New("err")).Once()
		games, err := gameService.gameRepository.GetTopGames(context.Background(), filters, uint32(10))
		require.Error(t, err)
		require.Nil(t, games)
	})
	t.Run("Если игр под фильтры не нашлось, ошибки нет", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		filters := model.GameFilters{ReleaseYear: 2020}
		gameMockRepo.On("GetTopGames", mock.Anything, filters, uint32(10)).Return(([]model.Game)(nil), nil).Once()
		games, err := gameService.gameRepository.GetTopGames(context.Background(), filters, uint32(10))
		require.NoError(t, err)
		require.Nil(t, games)
	})
	t.Run("Успешное получение топа игр", func(t *testing.T) {
		gameMockRepo := new(mockGameReposiroy)
		tagMockRepo := new(mockTagRepository)
		genreMockRepo := new(mockGenreRepository)
		s3Mock := new(mockS3Storager)
		mailerMock := new(mockEmailAlerter)
		gameService := NewGameService(mockslog.NewDiscardLogger(), gameMockRepo, tagMockRepo, genreMockRepo, s3Mock, mailerMock)
		filters := model.GameFilters{ReleaseYear: 2020}
		gameMockRepo.On("GetTopGames", mock.Anything, filters, uint32(10)).Return([]model.Game{{GameID: 1, Title: "qwe", Description: "qe"}}, nil).Once()
		games, err := gameService.gameRepository.GetTopGames(context.Background(), filters, uint32(10))
		require.NoError(t, err)
		require.Equal(t, games, []model.Game{{GameID: 1, Title: "qwe", Description: "qe"}})
	})
}
