//go:build integrations

package gamerepo

import (
	"context"
	"sort"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/storage/postgresql/gamerepo"
	"github.com/sariya23/game_service/tests/utils/random"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListGame(t *testing.T) {
	t.Run("Список игр без фильтрации, возвращаются все игры", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		sl := mockslog.NewDiscardLogger()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		gameRepo := gamerepo.NewGameRepository(dbT.DB, sl)
		totalGames := gofakeit.Number(5, 15)
		limit := uint32(totalGames - gofakeit.Number(3, 7))
		games := make([]dto.AddGameService, 0, totalGames)
		for range totalGames {
			gameToAdd := random.GameToAddService(genreIDs, tagIDs)
			gameID, err := gameRepo.SaveGame(ctx, gameToAdd)
			require.NoError(t, err)
			err = gameRepo.UpdateGameStatus(ctx, gameID, gamev2.GameStatusType_PUBLISH)
			require.NoError(t, err)
			games = append(games, gameToAdd)
		}

		// Act
		shortGames, err := gameRepo.GameList(ctx, dto.GameFilters{}, limit)

		// Assert
		require.NoError(t, err)
		sort.Slice(games, func(i, j int) bool {
			if games[i].Title != games[j].Title {
				return games[i].Title < games[j].Title
			} else {
				return games[i].ReleaseDate.Before(games[j].ReleaseDate)
			}
		})
		games = games[:limit]
		require.Equal(t, len(games), len(shortGames))
		for i := 0; i < len(games); i++ {
			assert.Equal(t, games[i].Title, shortGames[i].Title)
			assert.Equal(t, games[i].Description, shortGames[i].Description)
			assert.Equal(t, games[i].ReleaseDate, shortGames[i].ReleaseDate)
			assert.Equal(t, games[i].ImageURL, shortGames[i].ImageURL)
		}
	})
}
