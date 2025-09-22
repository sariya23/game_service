package tests

import (
	"context"
	"math/rand/v2"
	"net"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/lib/mockslog"
	"github.com/sariya23/game_service/internal/lib/random"
	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/model/dto"
	"github.com/sariya23/game_service/internal/storage/postgresql"
	"github.com/sariya23/game_service/tests/helpers"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGetTopGames(t *testing.T) {
	ctx := context.Background()
	cfg := config.MustLoadByPath("../config/local.env")
	db := postgresql.MustNewConnection(ctx, mockslog.NewDiscardLogger(), cfg.Postgres.PostgresURL)
	conn, err := grpc.NewClient(
		net.JoinHostPort("127.0.0.1", strconv.Itoa(cfg.Server.GrpcServerPort)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil || conn == nil {
		t.Fatalf("cannot start client; err = %v", err)
	}
	grpcClient := gamev4.NewGameServiceClient(conn)
	if grpcClient == nil {
		t.Fatal("cannot create grpcClient")
	}
	t.Run("Тест ручки GetTopGame; При пустом запросе вовзращается 10 игр", func(t *testing.T) {
		expctedGames, err := db.GetTopGames(ctx, dto.GameFilters{}, 10)
		require.NoError(t, err)
		req := gamev4.GetTopGamesRequest{}
		response, err := grpcClient.GetTopGames(ctx, &req)
		require.NoError(t, err)

		require.Equal(t, len(expctedGames), len(response.Games))
		for i := 0; i < len(expctedGames); i++ {
			gameDB := expctedGames[i]
			gameSRV := response.Games[i]

			assert.Equal(t, int(gameDB.GameID), int(gameSRV.ID))
			assert.Equal(t, gameDB.Title, gameSRV.Title)
			assert.Equal(t, gameDB.Description, gameSRV.Description)
			assert.Equal(t, int32(gameDB.ReleaseDate.Year()), gameSRV.GetReleaseDate().GetYear())
			assert.Equal(t, int32(gameDB.ReleaseDate.Month()), gameSRV.GetReleaseDate().GetMonth())
			assert.Equal(t, int32(gameDB.ReleaseDate.Day()), gameSRV.GetReleaseDate().GetDay())
			assert.Equal(t, gameDB.ImageURL, gameSRV.CoverImageUrl)
		}
	})
	t.Run("Тест ручки GetTopGame; Фильтрация по жанрам, без лимита", func(t *testing.T) {
		genres, err := db.GetGenres(ctx)
		require.NoError(t, err)
		genres = random.Sample(genres, 3)
		expectedGenreNames := model.GetGenreNames(genres)

		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Genres = expectedGenreNames
		gameToAdd.Tags = nil
		_, err = grpcClient.AddGame(ctx, &gamev4.AddGameRequest{Game: gameToAdd})
		require.NoError(t, err)

		expectedGenreNames = append(expectedGenreNames, gofakeit.MovieGenre())
		expctedGames, err := db.GetTopGames(ctx, dto.GameFilters{Genres: expectedGenreNames}, 10)
		require.NoError(t, err)
		req := gamev4.GetTopGamesRequest{Genres: expectedGenreNames}
		response, err := grpcClient.GetTopGames(ctx, &req)
		require.NoError(t, err)

		require.Equal(t, len(expctedGames), len(response.Games))
		for i := 0; i < len(expctedGames); i++ {
			gameDB := expctedGames[i]
			gameSRV := response.Games[i]

			fullGame, err := db.GetGameByID(ctx, gameDB.GameID)
			require.NoError(t, err)
			fullGameGenreNames := model.GetGenreNames(fullGame.Genres)
			assert.True(t, helpers.HasIntersection(expectedGenreNames, fullGameGenreNames))
			assert.Equal(t, int(gameDB.GameID), int(gameSRV.ID))
			assert.Equal(t, gameDB.Title, gameSRV.Title)
			assert.Equal(t, gameDB.Description, gameSRV.Description)
			assert.Equal(t, int32(gameDB.ReleaseDate.Year()), gameSRV.GetReleaseDate().GetYear())
			assert.Equal(t, int32(gameDB.ReleaseDate.Month()), gameSRV.GetReleaseDate().GetMonth())
			assert.Equal(t, int32(gameDB.ReleaseDate.Day()), gameSRV.GetReleaseDate().GetDay())
			assert.Equal(t, gameDB.ImageURL, gameSRV.CoverImageUrl)
		}
	})
	t.Run("Тест ручки GetTopGame; Фильтрация по тэгам, без лимита", func(t *testing.T) {
		tags, err := db.GetTags(ctx)
		require.NoError(t, err)
		tags = random.Sample(tags, 3)
		expectedTagNames := model.GetTagNames(tags)

		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Tags = expectedTagNames
		gameToAdd.Genres = nil
		_, err = grpcClient.AddGame(ctx, &gamev4.AddGameRequest{Game: gameToAdd})
		require.NoError(t, err)

		expectedTagNames = append(expectedTagNames, gofakeit.BookGenre())
		expctedGames, err := db.GetTopGames(ctx, dto.GameFilters{Tags: expectedTagNames}, 10)
		require.NoError(t, err)
		expectedTagNames = append(expectedTagNames, gofakeit.BookGenre())
		req := gamev4.GetTopGamesRequest{Tags: expectedTagNames}
		response, err := grpcClient.GetTopGames(ctx, &req)
		require.NoError(t, err)

		require.Equal(t, len(expctedGames), len(response.Games))
		for i := 0; i < len(expctedGames); i++ {
			gameDB := expctedGames[i]
			gameSRV := response.Games[i]

			fullGame, err := db.GetGameByID(ctx, gameDB.GameID)
			require.NoError(t, err)
			fullGameTagsNames := model.GetTagNames(fullGame.Tags)

			assert.True(t, helpers.HasIntersection(expectedTagNames, fullGameTagsNames))
			assert.Equal(t, int(gameDB.GameID), int(gameSRV.ID))
			assert.Equal(t, gameDB.Title, gameSRV.Title)
			assert.Equal(t, gameDB.Description, gameSRV.Description)
			assert.Equal(t, int32(gameDB.ReleaseDate.Year()), gameSRV.GetReleaseDate().GetYear())
			assert.Equal(t, int32(gameDB.ReleaseDate.Month()), gameSRV.GetReleaseDate().GetMonth())
			assert.Equal(t, int32(gameDB.ReleaseDate.Day()), gameSRV.GetReleaseDate().GetDay())
			assert.Equal(t, gameDB.ImageURL, gameSRV.CoverImageUrl)
		}
	})
	t.Run("Тест ручки GetTopGame; Фильтрация по тэгам и жанрам", func(t *testing.T) {
		tags, err := db.GetTags(ctx)
		require.NoError(t, err)
		genres, err := db.GetGenres(ctx)
		require.NoError(t, err)
		gameTagNames := model.GetTagNames(random.Sample(tags, 2))
		gameGenreNames := model.GetGenreNames(random.Sample(genres, 2))
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Genres = gameGenreNames
		gameToAdd.Tags = gameTagNames

		_, err = grpcClient.AddGame(ctx, &gamev4.AddGameRequest{Game: gameToAdd})
		require.NoError(t, err)

		expectedGames, err := db.GetTopGames(ctx, dto.GameFilters{
			Genres: gameGenreNames,
			Tags:   gameTagNames,
		}, 10)
		require.NoError(t, err)

		gameGenreNames = append(gameGenreNames, gofakeit.BookGenre())
		gameTagNames = append(gameTagNames, gofakeit.BookGenre())
		resp, err := grpcClient.GetTopGames(ctx,
			&gamev4.GetTopGamesRequest{
				Genres: gameGenreNames,
				Tags:   gameTagNames,
			})
		require.NoError(t, err)
		require.Equal(t, len(expectedGames), len(resp.Games))

		for i := 0; i < len(expectedGames); i++ {
			gameDB := expectedGames[i]
			gameSRV := resp.Games[i]

			assert.Equal(t, gameDB.GameID, gameSRV.ID)
			fullGame, err := grpcClient.GetGame(ctx, &gamev4.GetGameRequest{GameId: gameSRV.ID})
			require.NoError(t, err)
			assert.True(t, helpers.HasIntersection(fullGame.Game.Genres, gameGenreNames))
			assert.True(t, helpers.HasIntersection(fullGame.Game.Tags, gameTagNames))
			assert.Equal(t, gameDB.Title, gameSRV.Title)
			assert.Equal(t, gameDB.Description, gameSRV.Description)
			assert.Equal(t, int32(gameDB.ReleaseDate.Year()), gameSRV.GetReleaseDate().GetYear())
			assert.Equal(t, int32(gameDB.ReleaseDate.Month()), gameSRV.GetReleaseDate().GetMonth())
			assert.Equal(t, int32(gameDB.ReleaseDate.Day()), gameSRV.GetReleaseDate().GetDay())
			assert.Equal(t, gameDB.ImageURL, gameSRV.CoverImageUrl)
		}
	})
	t.Run("Тест ручки GetTopGame; Фильтрация по тэгам, жанрам и году", func(t *testing.T) {
		tags, err := db.GetTags(ctx)
		require.NoError(t, err)
		genres, err := db.GetGenres(ctx)
		require.NoError(t, err)
		gameTagNames := model.GetTagNames(random.Sample(tags, 2))
		gameGenreNames := model.GetGenreNames(random.Sample(genres, 2))
		gameToAdd := random.RandomAddGameRequest()
		gameToAdd.Genres = gameGenreNames
		gameToAdd.Tags = gameTagNames

		_, err = grpcClient.AddGame(ctx, &gamev4.AddGameRequest{Game: gameToAdd})
		require.NoError(t, err)

		gameGenreNames = append(gameGenreNames, gofakeit.BookGenre())
		gameTagNames = append(gameTagNames, gofakeit.BookGenre())
		expectedGames, err := db.GetTopGames(ctx, dto.GameFilters{
			Genres:      gameGenreNames,
			Tags:        gameTagNames,
			ReleaseYear: gameToAdd.ReleaseDate.Year,
		}, 10)
		require.NoError(t, err)

		resp, err := grpcClient.GetTopGames(ctx,
			&gamev4.GetTopGamesRequest{
				Genres: gameGenreNames,
				Tags:   gameTagNames,
				Year:   gameToAdd.ReleaseDate.Year,
			})
		require.NoError(t, err)
		require.Equal(t, len(expectedGames), len(resp.Games))

		for i := 0; i < len(expectedGames); i++ {
			gameDB := expectedGames[i]
			gameSRV := resp.Games[i]

			assert.Equal(t, gameDB.GameID, gameSRV.ID)
			fullGame, err := grpcClient.GetGame(ctx, &gamev4.GetGameRequest{GameId: gameSRV.ID})
			require.NoError(t, err)
			assert.True(t, helpers.HasIntersection(fullGame.Game.Genres, gameGenreNames))
			assert.True(t, helpers.HasIntersection(fullGame.Game.Tags, gameTagNames))
			assert.Equal(t, gameDB.Title, gameSRV.Title)
			assert.Equal(t, gameDB.Description, gameSRV.Description)
			assert.Equal(t, int32(gameDB.ReleaseDate.Year()), gameSRV.GetReleaseDate().GetYear())
			assert.Equal(t, int32(gameDB.ReleaseDate.Month()), gameSRV.GetReleaseDate().GetMonth())
			assert.Equal(t, int32(gameDB.ReleaseDate.Day()), gameSRV.GetReleaseDate().GetDay())
			assert.Equal(t, gameDB.ImageURL, gameSRV.CoverImageUrl)
		}
	})
	t.Run("Тест ручки GetTopGame; Указан только лимит", func(t *testing.T) {
		limit := rand.IntN(15) + 1
		expctedGames, err := db.GetTopGames(ctx, dto.GameFilters{}, uint32(limit))
		require.NoError(t, err)
		req := gamev4.GetTopGamesRequest{Limit: uint32(limit)}
		response, err := grpcClient.GetTopGames(ctx, &req)
		require.NoError(t, err)

		require.Equal(t, len(expctedGames), len(response.Games))
		for i := 0; i < len(expctedGames); i++ {
			gameDB := expctedGames[i]
			gameSRV := response.Games[i]

			assert.Equal(t, int(gameDB.GameID), int(gameSRV.ID))
			assert.Equal(t, gameDB.Title, gameSRV.Title)
			assert.Equal(t, gameDB.Description, gameSRV.Description)
			assert.Equal(t, int32(gameDB.ReleaseDate.Year()), gameSRV.GetReleaseDate().GetYear())
			assert.Equal(t, int32(gameDB.ReleaseDate.Month()), gameSRV.GetReleaseDate().GetMonth())
			assert.Equal(t, int32(gameDB.ReleaseDate.Day()), gameSRV.GetReleaseDate().GetDay())
			assert.Equal(t, gameDB.ImageURL, gameSRV.CoverImageUrl)
		}
	})
}
