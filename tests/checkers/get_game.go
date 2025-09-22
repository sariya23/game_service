package checkers

import (
	"context"
	"io"
	"testing"

	"github.com/sariya23/game_service/internal/outerror"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AssertGetGame(ctx context.Context,
	t *testing.T,
	expected *gamev4.GameRequest,
	response *gamev4.GetGameResponse,
	s3 *minioclient.Minio,
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
	obj, err := s3.GetObject(ctx, response.Game.GetCoverImageUrl())
	require.NoError(t, err)
	imageBytes, err := io.ReadAll(obj)
	require.NoError(t, err)
	assert.Equal(t, expected.GetCoverImage(), imageBytes)
}

func AssertGetGameNotFound(t *testing.T, err error, response *gamev4.GetGameResponse) {
	t.Helper()
	s, _ := status.FromError(err)
	require.Equal(t, codes.NotFound, s.Code())
	require.Equal(t, outerror.GameNotFoundMessage, s.Message())
	require.Nil(t, response.GetGame())
}
