//go:build integrations

package game_test

import (
	"context"
	"slices"
	"sort"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/tests/clientgrpc"
	"github.com/sariya23/game_service/tests/utils/random"
	"github.com/sariya23/game_service/tests/utils/struct"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGameList(t *testing.T) {
	t.Run("Список игр без фильтров, дефолтный лимит", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		n := gofakeit.Number(12, 20)
		gameIDs := make([]int64, 0, n)
		for range n {
			gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
			request := gamev2.AddGameRequest{Game: gameToAdd}
			responseAdd, err := client.GetClient().AddGame(ctx, &request)
			gameIDs = append(gameIDs, responseAdd.GameId)
			require.NoError(t, err)
			assert.NotZero(t, responseAdd.GameId)
			dbT.UpdateGameStatus(ctx, responseAdd.GameId, gamev2.GameStatusType_PUBLISH)
		}

		response, err := client.GetClient().GameList(ctx, &gamev2.GameListRequest{})
		require.NoError(t, err)
		assert.Len(t, response.Games, 10)
		expectedGames := make([]model.Game, 0, n)
		for _, gameID := range gameIDs {
			expectedGames = append(expectedGames, *dbT.GetGameById(ctx, gameID))
		}
		sort.Slice(expectedGames, func(i, j int) bool {
			if expectedGames[i].Title != expectedGames[j].Title {
				return expectedGames[i].Title < expectedGames[j].Title
			}
			return expectedGames[i].ReleaseDate.Before(expectedGames[j].ReleaseDate)
		})
		expectedGames = expectedGames[:10]
		for i, expectedGame := range expectedGames {
			assert.Equal(t, expectedGame.Title, response.Games[i].Title)
			assert.Equal(t, expectedGame.Description, response.Games[i].Description)
			assert.Equal(t, expectedGame.ImageURL, response.Games[i].CoverImageUrl)
		}
	})
	t.Run("Список игр, фильтрация по годам, дефолтный лимит", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		n := gofakeit.Number(12, 20)
		games := make([]model.Game, 0, n)
		for range n {
			gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
			request := gamev2.AddGameRequest{Game: gameToAdd}
			responseAdd, err := client.GetClient().AddGame(ctx, &request)
			games = append(games, *dbT.GetGameById(ctx, responseAdd.GameId))
			require.NoError(t, err)
			assert.NotZero(t, responseAdd.GameId)
			dbT.UpdateGameStatus(ctx, responseAdd.GameId, gamev2.GameStatusType_PUBLISH)
		}
		targetYear := _struct.PickMostFrequentValue(func() []int32 {
			years := make([]int32, 0, n)
			for _, v := range games {
				years = append(years, int32(v.ReleaseDate.Year()))
			}
			return years
		}())
		expectedGames := func() []model.Game {
			var res []model.Game
			for _, game := range games {
				if game.ReleaseDate.Year() == int(targetYear) {
					res = append(res, game)
				}
			}
			return res
		}()
		limit := 10
		if len(expectedGames) < 10 {
			limit = len(expectedGames)
		}
		expectedGames = expectedGames[:limit]
		response, err := client.
			GetClient().
			GameList(ctx, &gamev2.GameListRequest{Year: targetYear})
		require.NoError(t, err)
		assert.Len(t, response.Games, limit)
		sort.Slice(expectedGames, func(i, j int) bool {
			if expectedGames[i].Title != expectedGames[j].Title {
				return expectedGames[i].Title < expectedGames[j].Title
			}
			return expectedGames[i].ReleaseDate.Before(expectedGames[j].ReleaseDate)
		})
		for i, expectedGame := range expectedGames {
			assert.Equal(t, expectedGame.Title, response.Games[i].Title)
			assert.Equal(t, expectedGame.Description, response.Games[i].Description)
			assert.Equal(t, expectedGame.ImageURL, response.Games[i].CoverImageUrl)
		}
	})
	t.Run("Список игр, фильтрация по жанрам, дефолтный лимит", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		n := gofakeit.Number(12, 20)
		games := make([]model.Game, 0, n)
		for range n {
			gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
			request := gamev2.AddGameRequest{Game: gameToAdd}
			responseAdd, err := client.GetClient().AddGame(ctx, &request)
			games = append(games, *dbT.GetGameById(ctx, responseAdd.GameId))
			require.NoError(t, err)
			assert.NotZero(t, responseAdd.GameId)
			dbT.UpdateGameStatus(ctx, responseAdd.GameId, gamev2.GameStatusType_PUBLISH)
		}

		targetGenres := random.Sample(func() []string {
			var res []string
			for _, game := range games {
				for _, genre := range game.Genres {
					res = append(res, genre.GenreName)
				}
			}
			return res
		}(), 3)
		expectedGames := func() []model.Game {
			var res []model.Game
			for _, game := range games {
				if slices.Contains(game.Genres) {
					res = append(res, game)
				}
			}
			return res
		}()
		limit := 10
		if len(expectedGames) < 10 {
			limit = len(expectedGames)
		}
		expectedGames = expectedGames[:limit]
		response, err := client.
			GetClient().
			GameList(ctx, &gamev2.GameListRequest{Year: targetYear})
		require.NoError(t, err)
		assert.Len(t, response.Games, limit)
		sort.Slice(expectedGames, func(i, j int) bool {
			if expectedGames[i].Title != expectedGames[j].Title {
				return expectedGames[i].Title < expectedGames[j].Title
			}
			return expectedGames[i].ReleaseDate.Before(expectedGames[j].ReleaseDate)
		})
		for i, expectedGame := range expectedGames {
			assert.Equal(t, expectedGame.Title, response.Games[i].Title)
			assert.Equal(t, expectedGame.Description, response.Games[i].Description)
			assert.Equal(t, expectedGame.ImageURL, response.Games[i].CoverImageUrl)
		}
	})
}
