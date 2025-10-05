//go:build integrations

package gamerepo

import (
	"context"
	"testing"

	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/storage/postgresql/gamerepo"
	"github.com/sariya23/game_service/tests/utils/random"
	"github.com/stretchr/testify/require"
)

func TestSaveGame_success(t *testing.T) {
	ctx := context.Background()
	sl := mockslog.NewDiscardLogger()
	dbT.SetUp(ctx, t, tables...)
	defer dbT.TearDown(t)
	gameRepo := gamerepo.NewGameRepository(dbT.DB, sl)
	gameToAdd := random.GameToAddService(genreIDs, tagIDs)

	gameID, err := gameRepo.SaveGame(ctx, gameToAdd)
	require.NoError(t, err)
	require.NotZero(t, gameID)
}
