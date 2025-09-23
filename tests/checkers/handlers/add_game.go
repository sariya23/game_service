package checkers

import (
	"testing"

	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AssertAddGame проверяет ответ ручки AddGame при успешном сохранении
func AssertAddGame(t *testing.T, err error, response *gamev4.AddGameResponse) {
	t.Helper()
	require.NoError(t, err)
	require.NotZero(t, response.GetGameId())
}

// AssertAddGameGenreNotFound проверяет ответ ручки AddGame при попытке создании игры с
// несуществующим жанром
func AssertAddGameTagNotFound(t *testing.T, err error, response *gamev4.AddGameResponse) {
	t.Helper()
	s, _ := status.FromError(err)
	require.Equal(t, codes.InvalidArgument, s.Code())
	require.Equal(t, outerror.TagNotFoundMessage, s.Message())
	require.Zero(t, response.GetGameId())
}

// AssertAddGameGenreNotFound проверяет ответ ручки AddGame при попытке создании игры с
// несуществующим жанром
func AssertAddGameGenreNotFound(t *testing.T, err error, response *gamev4.AddGameResponse) {
	t.Helper()
	s, _ := status.FromError(err)
	require.Equal(t, codes.InvalidArgument, s.Code())
	require.Equal(t, outerror.GenreNotFoundMessage, s.Message())
	require.Zero(t, response.GetGameId())
}

// AssertAddGameDuplicateGame проверяет ответ ручки при попытке создать дублирующую игру
func AssertAddGameDuplicateGame(t *testing.T, err error, response *gamev4.AddGameResponse) {
	t.Helper()
	s, _ := status.FromError(err)
	require.Equal(t, codes.AlreadyExists, s.Code())
	require.Equal(t, outerror.GameAlreadyExistMessage, s.Message())
	require.Zero(t, response.GetGameId())
}
