package checkers

import (
	"testing"

	"github.com/sariya23/game_service/internal/outerror"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AssertDeleteGame проверяет ответ ручки DeleteGame
func AssertDeleteGame(t *testing.T, err error, expectedID uint64, resp *gamev4.DeleteGameResponse) {
	t.Helper()
	require.NoError(t, err)
	require.Equal(t, expectedID, resp.GameId)
}

func AssertDeleteGameNotFound(t *testing.T, err error, response *gamev4.DeleteGameResponse) {
	t.Helper()
	s, _ := status.FromError(err)
	require.Equal(t, int(codes.NotFound), int(s.Code()))
	require.Equal(t, outerror.GameNotFoundMessage, s.Message())
	require.Nil(t, response)
}
