package checkers

import (
	"context"
	"testing"

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
	expected *gamev4.GameRequest,
	response *gamev4.GetGameResponse,
	expctedImageS3 []byte,
	err error,
) {
	t.Helper()
	require.NoError(t, err)
	assert.Equal(t, expected.GetTitle(), response.Game.GetTitle())
	assert.Equal(t, expected.GetDescription(), response.Game.GetDescription())
	assert.Equal(t, expected.GetReleaseDate().GetYear(), response.Game.GetReleaseDate().GetYear())
	assert.Equal(t, expected.GetReleaseDate().GetMonth(), response.Game.GetReleaseDate().GetMonth())
	assert.Equal(t, expected.GetReleaseDate().GetDay(), response.Game.GetReleaseDate().GetDay())
	assert.Equal(t, expected.GetGenres(), response.Game.GetGenres())
	assert.Equal(t, expected.GetTags(), response.Game.GetTags())
	assert.Equal(t, expected.GetCoverImage(), expctedImageS3)
}

// AssertGetGameNotFound проверяет ответ ручки GetGame при попытке получить несуществующую игру
func AssertGetGameNotFound(t *testing.T, err error, response *gamev4.GetGameResponse) {
	t.Helper()
	s, _ := status.FromError(err)
	require.Equal(t, codes.NotFound, s.Code())
	require.Equal(t, outerror.GameNotFoundMessage, s.Message())
	require.Nil(t, response.GetGame())
}
