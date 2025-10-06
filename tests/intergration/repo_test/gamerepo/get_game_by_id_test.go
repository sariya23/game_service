//go:build integrations

package gamerepo

import (
	"context"
	"sort"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/game_service/internal/storage/postgresql/gamerepo"
	"github.com/sariya23/game_service/tests/utils/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGameByID(t *testing.T) {
	t.Run("Успешное получение игры", func(t *testing.T) {
		ctx := context.Background()
		sl := mockslog.NewDiscardLogger()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		gameRepo := gamerepo.NewGameRepository(dbT.DB, sl)
		gameToAdd := random.GameToAddService(genreIDs, tagIDs)
		gameID, err := gameRepo.SaveGame(ctx, gameToAdd)
		require.NoError(t, err)
		assert.NotZero(t, gameID)

		game, err := gameRepo.GetGameByID(ctx, gameID)
		require.NoError(t, err)
		assert.Equal(t, gameToAdd.Title, game.Title)
		assert.Equal(t, gameToAdd.Description, game.Description)
		assert.Equal(t, gameToAdd.ImageURL, game.ImageURL)
		assert.Equal(t, gameToAdd.ReleaseDate, game.ReleaseDate)
		sort.Slice(gameToAdd.TagIDs, func(i, j int) bool {
			return gameToAdd.TagIDs[i] < gameToAdd.TagIDs[j]
		})
		sort.Slice(game.Tags, func(i, j int) bool {
			return game.Tags[i].TagID < game.Tags[j].TagID
		})
		sort.Slice(gameToAdd.GenreIDs, func(i, j int) bool {
			return gameToAdd.GenreIDs[i] < gameToAdd.GenreIDs[j]
		})
		sort.Slice(game.Genres, func(i, j int) bool {
			return game.Genres[i].GenreID < game.Genres[j].GenreID
		})
		assert.Equal(t, gameToAdd.TagIDs, model.TagIDs(game.Tags))
		assert.Equal(t, gameToAdd.GenreIDs, model.GenreIDs(game.Genres))
	})
	t.Run("Игра не найдена", func(t *testing.T) {
		ctx := context.Background()
		sl := mockslog.NewDiscardLogger()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		gameRepo := gamerepo.NewGameRepository(dbT.DB, sl)
		game, err := gameRepo.GetGameByID(ctx, gofakeit.Int64())
		require.ErrorIs(t, err, outerror.ErrGameNotFound)
		assert.Nil(t, game)
	})
}
