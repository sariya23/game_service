//go:build integrations

package game_test

import (
	"context"
	"sort"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/tests/clientgrpc"
	"github.com/sariya23/game_service/tests/utils/ds"
	"github.com/sariya23/game_service/tests/utils/random"
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
			require.NoError(t, err)
			gameIDs = append(gameIDs, responseAdd.GameId)
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
		targetYear := ds.PickMostFrequentValue(func() []int32 {
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
				res = append(res, model.GenreNames(game.Genres)...)
			}
			return res
		}(), 3)
		expectedGames := func() []model.Game {
			var res []model.Game
			for _, game := range games {
				if len(ds.Intersection(model.GenreNames(game.Genres), targetGenres)) > 0 {
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
			GameList(ctx, &gamev2.GameListRequest{Genres: targetGenres})
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
	t.Run("Список игр, фильтрация по тэгам, дефолтный лимит", func(t *testing.T) {
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
			require.NoError(t, err)
			assert.NotZero(t, responseAdd.GameId)
			games = append(games, *dbT.GetGameById(ctx, responseAdd.GameId))
			dbT.UpdateGameStatus(ctx, responseAdd.GameId, gamev2.GameStatusType_PUBLISH)
		}

		targetTags := random.Sample(func() []string {
			var res []string
			for _, game := range games {
				res = append(res, model.TagNames(game.Tags)...)
			}
			return res
		}(), 3)
		expectedGames := func() []model.Game {
			var res []model.Game
			for _, game := range games {
				if len(ds.Intersection(model.TagNames(game.Tags), targetTags)) > 0 {
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
			GameList(ctx, &gamev2.GameListRequest{Tags: targetTags})
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
	t.Run("Список игр, все поля", func(t *testing.T) {
		ctx := context.Background()
		client := clientgrpc.NewGameServiceTestClient()
		dbT.SetUp(ctx, t, tables...)
		defer dbT.TearDown(t)
		genres, tags := dbT.GetGenres(ctx), dbT.GetTags(ctx)
		n := gofakeit.Number(20, 40)
		games := make([]model.Game, 0, n)
		for range n {
			gameToAdd := random.GameToAddRequest(model.GenreNames(genres), model.TagNames(tags))
			request := gamev2.AddGameRequest{Game: gameToAdd}
			responseAdd, err := client.GetClient().AddGame(ctx, &request)
			require.NoError(t, err)
			assert.NotZero(t, responseAdd.GameId)
			games = append(games, *dbT.GetGameById(ctx, responseAdd.GameId))
			dbT.UpdateGameStatus(ctx, responseAdd.GameId, gamev2.GameStatusType_PUBLISH)
		}

		targetTags := random.Sample(func() []string {
			var res []string
			for _, game := range games {
				res = append(res, model.TagNames(game.Tags)...)
			}
			return res
		}(), 3)
		targetGenres := random.Sample(func() []string {
			var res []string
			for _, game := range games {
				res = append(res, model.GenreNames(game.Genres)...)
			}
			return res
		}(), 3)
		targetYear := ds.PickMostFrequentValue(func() []int32 {
			years := make([]int32, 0, n)
			for _, v := range games {
				years = append(years, int32(v.ReleaseDate.Year()))
			}
			return years
		}())
		expectedGames := func() []model.Game {
			var res []model.Game
			for _, game := range games {
				if len(ds.Intersection(model.TagNames(game.Tags), targetTags)) > 0 && len(ds.Intersection(model.GenreNames(game.Genres), targetGenres)) > 0 && game.ReleaseDate.Year() == int(targetYear) {
					res = append(res, game)
				}
			}
			return res
		}()
		limit := gofakeit.Number(3, 7)
		if len(expectedGames) < limit {
			limit = len(expectedGames)
		}
		expectedGames = expectedGames[:limit]
		response, err := client.
			GetClient().
			GameList(ctx, &gamev2.GameListRequest{Tags: targetTags, Genres: targetGenres, Year: targetYear, Limit: uint32(limit)})
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
