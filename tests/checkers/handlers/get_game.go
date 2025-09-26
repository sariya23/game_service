package checkers

import (
	"context"
	"testing"

	"github.com/sariya23/game_service/internal/model"
	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AssertGetGame проверяет ответ ручки GetGame
func AssertGetGame(ctx context.Context,
	t *testing.T,
	expected *model.Game,
	response *gamev4.GetGameResponse,
	imageDB []byte,
	expectedImage []byte,
) {
	t.Helper()
	assert.Equal(t, expected.Title, response.Game.GetTitle())
	assert.Equal(t, expected.Description, response.Game.GetDescription())
	assert.Equal(t, int32(expected.ReleaseDate.Year()), response.Game.GetReleaseDate().GetYear())
	assert.Equal(t, int32(expected.ReleaseDate.Month()), response.Game.GetReleaseDate().GetMonth())
	assert.Equal(t, int32(expected.ReleaseDate.Day()), response.Game.GetReleaseDate().GetDay())
	assert.Equal(t, model.GetGenreNames(expected.Genres), response.Game.GetGenres())
	assert.Equal(t, model.GetTagNames(expected.Tags), response.Game.GetTags())
}

// AssertGetGameNotFound проверяет ответ ручки GetGame при попытке получить несуществующую игру
func AssertGetGameNotFound(t *testing.T, err error, response *gamev4.GetGameResponse) {
	t.Helper()
	s, _ := status.FromError(err)
	require.Equal(t, codes.NotFound, s.Code())
	require.Equal(t, outerror.GameNotFoundMessage, s.Message())
	require.Nil(t, response.GetGame())
}
