//go:build integrations

package gamerepo

import (
	"context"
	"testing"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/storage/postgresql/gamerepo"
	"github.com/sariya23/game_service/tests/utils/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveGame(t *testing.T) {
	t.Run("Успешное сохранение игры со всеми полями", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		sl := mockslog.NewDiscardLogger()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		gameRepo := gamerepo.NewGameRepository(dbT.DB, sl)
		gameToAdd := random.GameToAddService(genreIDs, tagIDs)

		// Act
		gameID, err := gameRepo.SaveGame(ctx, gameToAdd)

		// Assert
		require.NoError(t, err)
		assert.NotZero(t, gameID)
	})
}
