package checkers

import (
	"context"
	"testing"

	"github.com/sariya23/game_service/internal/model"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/assert"
)

// AssertAddGameRequestAndDB сравнивает данные для сохранения игры с сохраненной игрой
func AssertAddGameRequestAndDB(ctx context.Context,
	t *testing.T,
	request *gamev4.AddGameRequest,
	gameDB model.Game,
	expctedImageS3 []byte,
) {
	t.Helper()
	assert.Equal(t, request.Game.GetTitle(), gameDB.Title)
	assert.Equal(t, request.Game.GetDescription(), gameDB.Description)
	assert.Equal(t, request.Game.GetReleaseDate().GetYear(), int32(gameDB.ReleaseDate.Year()))
	assert.Equal(t, request.Game.GetReleaseDate().GetMonth(), int32(gameDB.ReleaseDate.Month()))
	assert.Equal(t, request.Game.GetReleaseDate().GetDay(), int32(gameDB.ReleaseDate.Day()))
	assert.Equal(t, request.Game.GetTags(), model.GetTagNames(gameDB.Tags))
	assert.Equal(t, request.Game.GetGenres(), model.GetGenreNames(gameDB.Genres))
	assert.Equal(t, request.Game.CoverImage, expctedImageS3)
}
