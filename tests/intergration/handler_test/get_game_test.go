//go:build integrations

package handler_test

import (
	"context"
	"testing"

	"github.com/sariya23/game_service/internal/lib/converters"
	"github.com/sariya23/game_service/tests/clientgrpc"
	"github.com/sariya23/game_service/tests/utils/filldb"
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
		gameID := filldb.InsertGame(ctx, dbT, gameToAdd)
		filldb.InsertGameGenre(ctx, dbT, gameID, gameToAdd.GenreIDs)
		filldb.InsertGameTag(ctx, dbT, gameID, gameToAdd.TagIDs)
		request := gamev2.GetGameRequest{GameId: gameID}

		response, err := client.GetClient().GetGame(ctx, &request)
		require.NoError(t, err)
		assert.Equal(t, gameID, response.GetGame().ID)
		assert.Equal(t, gameToAdd.Title, response.GetGame().Title)
		assert.Equal(t, gameToAdd.Description, response.GetGame().Description)
		assert.Equal(t, gameToAdd.ReleaseDate, converters.FromProtoDate(response.GetGame().ReleaseDate))
		assert.Equal(t, gameToAdd.ImageURL, response.GetGame().CoverImageUrl)

	})
}
