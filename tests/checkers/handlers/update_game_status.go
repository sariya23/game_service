package checkers

import (
	"testing"

	"github.com/sariya23/game_service/internal/outerror"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AssertUpdateGameStatusGameNotFound(t *testing.T, err error) {
	t.Helper()
	s, _ := status.FromError(err)
	require.Equal(t, codes.NotFound, s.Code())
	require.Equal(t, outerror.GameNotFoundMessage, s.Message())
}

func AssertUpdateGameStatusInvalidStatus(t *testing.T, err error) {
	t.Helper()
	s, _ := status.FromError(err)
	require.Equal(t, codes.InvalidArgument, s.Code())
	require.Equal(t, outerror.InvalidNewGameStatusMessage, s.Message())
}
