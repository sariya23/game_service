package handlers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sariya23/api_game_service/gen/game"
	"github.com/sariya23/game_service/internal/outerror"
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
		expectedResponse *game.AddGameResponse
	}{
		{
			name:             "GameAlreadyExist",
			err:              fmt.Errorf("%s: %w", "qweo", outerror.ErrGameAlreadyExist),
			gameID:           0,
			expectedError:    status.Error(codes.AlreadyExists, outerror.GameAlreadyExistMessage),
			expectedResponse: &game.AddGameResponse{},
		},
		{
			name:             "ErrCannotSaveGameImage",
			err:              fmt.Errorf("%s: %w", "qweo", outerror.ErrCannotSaveGameImage),
			gameID:           23,
			expectedError:    nil,
			expectedResponse: &game.AddGameResponse{GameId: 23},
		},
		{
			name:             "GenreNotFound",
			err:              fmt.Errorf("%s: %w", "qweo", outerror.ErrGenreNotFound),
			gameID:           0,
			expectedError:    status.Error(codes.InvalidArgument, outerror.GenreNotFoundMessage),
			expectedResponse: &game.AddGameResponse{},
		},
		{
			name:             "TagNotFound",
			err:              fmt.Errorf("%s: %w", "qweo", outerror.ErrTagNotFound),
			gameID:           0,
			expectedError:    status.Error(codes.InvalidArgument, outerror.TagNotFoundMessage),
			expectedResponse: &game.AddGameResponse{},
		},
		{
			name:             "Some error",
			err:              errors.New("some error"),
			gameID:           0,
			expectedError:    status.Error(codes.Internal, outerror.InternalMessage),
			expectedResponse: &game.AddGameResponse{},
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
