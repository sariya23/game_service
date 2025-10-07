//go:build integrations

package game_test

import (
	"context"
	"sort"
	"testing"

	"github.com/sariya23/game_service/internal/lib/converters"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/tests/clientgrpc"
	"github.com/sariya23/game_service/tests/utils/random"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGame(t *testing.T) {
	t.Run("Успешное получение игры", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		gameToAdd := random.GameToAddService(tagIDs, genreIDs)
		gameID := dbT.InsertGame(ctx, gameToAdd)
		dbT.InsertGameGenre(ctx, gameID, gameToAdd.GenreIDs)
		dbT.InsertGameTag(ctx, gameID, gameToAdd.TagIDs)
		request := gamev2.GetGameRequest{GameId: gameID}

		response, err := client.GetClient().GetGame(ctx, &request)
		require.NoError(t, err)
		assert.Equal(t, gameID, response.GetGame().ID)
		assert.Equal(t, gameToAdd.Title, response.GetGame().Title)
		assert.Equal(t, gameToAdd.Description, response.GetGame().Description)
		assert.Equal(t, converters.ToProtoDate(gameToAdd.ReleaseDate), response.GetGame().ReleaseDate)
		assert.Equal(t, gameToAdd.ImageURL, response.GetGame().CoverImageUrl)

		expectedGenres := dbT.GetGenresByIDs(ctx, genreIDs)
		sort.Slice(expectedGenres, func(i, j int) bool {
			return expectedGenres[i].GenreName < expectedGenres[j].GenreName
		})
		actualGenres := response.GetGame().GetGenres()
		sort.Slice(actualGenres, func(i, j int) bool {
			return actualGenres[i] < actualGenres[j]
		})
		assert.Equal(t, model.GenreNames(expectedGenres), actualGenres)
	})
}
