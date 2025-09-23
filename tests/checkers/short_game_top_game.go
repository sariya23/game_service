package checkers

import (
	"testing"

	"github.com/sariya23/game_service/internal/model"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/assert"
)

// AssertShortGameTopGame сранивает игру из БД с игрой, которая пришла из ответа GetTopGames
func AssertShortGameTopGame(t *testing.T, gameDB model.ShortGame, gameSRV *gamev4.GetTopGamesResponse_ShortGame) {
	t.Helper()
	assert.Equal(t, int(gameDB.GameID), int(gameSRV.ID))
	assert.Equal(t, gameDB.Title, gameSRV.Title)
	assert.Equal(t, gameDB.Description, gameSRV.Description)
	assert.Equal(t, int32(gameDB.ReleaseDate.Year()), gameSRV.GetReleaseDate().GetYear())
	assert.Equal(t, int32(gameDB.ReleaseDate.Month()), gameSRV.GetReleaseDate().GetMonth())
	assert.Equal(t, int32(gameDB.ReleaseDate.Day()), gameSRV.GetReleaseDate().GetDay())
	assert.Equal(t, gameDB.ImageURL, gameSRV.CoverImageUrl)
}
