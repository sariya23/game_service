package checkers

import (
	"context"
	"io"
	"testing"

	"github.com/sariya23/game_service/internal/model"
	minioclient "github.com/sariya23/game_service/internal/storage/s3/minio"
	gamev4 "github.com/sariya23/proto_api_games/v4/gen/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertAddGameRequestAndDB сравнивает данные для сохранения игры с сохраненной игрой
func AssertAddGameRequestAndDB(ctx context.Context,
	t *testing.T,
	request *gamev4.AddGameRequest,
	gameDB model.Game,
	s3 *minioclient.Minio,
) {
	t.Helper()
	assert.Equal(t, request.Game.GetTitle(), gameDB.Title)
	assert.Equal(t, request.Game.GetDescription(), gameDB.Description)
	assert.Equal(t, request.Game.GetReleaseDate().GetYear(), int32(gameDB.ReleaseDate.Year()))
	assert.Equal(t, request.Game.GetReleaseDate().GetMonth(), int32(gameDB.ReleaseDate.Month()))
	assert.Equal(t, request.Game.GetReleaseDate().GetDay(), int32(gameDB.ReleaseDate.Day()))
	assert.Equal(t, request.Game.GetTags(), model.GetTagNames(gameDB.Tags))
	assert.Equal(t, request.Game.GetGenres(), model.GetGenreNames(gameDB.Genres))
	image, err := s3.GetObject(ctx, minioclient.GameKey(request.Game.GetTitle(), int(request.Game.ReleaseDate.GetYear())))
	require.NoError(t, err)
	imageBytes, err := io.ReadAll(image)
	require.NoError(t, err)
	assert.Equal(t, request.Game.CoverImage, imageBytes)
}
