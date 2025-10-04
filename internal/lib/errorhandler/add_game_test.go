package errorhandler

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddGame_errorhandler(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name             string
		err              error
		gameID           int64
		expectedError    error
		expectedResponse *gamev2.AddGameResponse
	}{
		{
			name:             "GameAlreadyExist",
			err:              fmt.Errorf("%s: %w", "qweo", outerror.ErrGameAlreadyExist),
			gameID:           0,
			expectedError:    status.Error(codes.AlreadyExists, outerror.GameAlreadyExistMessage),
			expectedResponse: &gamev2.AddGameResponse{},
		},
		{
			name:             "ErrCannotSaveGameImage",
			err:              fmt.Errorf("%s: %w", "qweo", outerror.ErrCannotSaveGameImage),
			gameID:           23,
			expectedError:    nil,
			expectedResponse: &gamev2.AddGameResponse{GameId: 23},
		},
		{
			name:             "GenreNotFound",
			err:              fmt.Errorf("%s: %w", "qweo", outerror.ErrGenreNotFound),
			gameID:           0,
			expectedError:    status.Error(codes.InvalidArgument, outerror.GenreNotFoundMessage),
			expectedResponse: &gamev2.AddGameResponse{},
		},
		{
			name:             "TagNotFound",
			err:              fmt.Errorf("%s: %w", "qweo", outerror.ErrTagNotFound),
			gameID:           0,
			expectedError:    status.Error(codes.InvalidArgument, outerror.TagNotFoundMessage),
			expectedResponse: &gamev2.AddGameResponse{},
		},
		{
			name:             "Some error",
			err:              errors.New("some error"),
			gameID:           0,
			expectedError:    status.Error(codes.Internal, outerror.InternalMessage),
			expectedResponse: &gamev2.AddGameResponse{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotResp, gotErr := AddGame(tc.err, tc.gameID)
			assert.Equal(t, tc.expectedResponse, gotResp)
			assert.Equal(t, tc.expectedError, gotErr)
		})
	}
}
