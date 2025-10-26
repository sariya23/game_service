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

func TestGetGame_errorhandler(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name             string
		err              error
		expectedResponse *game.GetGameResponse
		expectedError    error
	}{
		{
			name:             "GameNotFound",
			err:              fmt.Errorf("%s: %w", "qwe", outerror.ErrGameNotFound),
			expectedResponse: &game.GetGameResponse{},
			expectedError:    status.Error(codes.NotFound, outerror.GameNotFoundMessage),
		},
		{
			name:             "SomeErr",
			err:              fmt.Errorf("%s: %w", "qwe", errors.New("some error")),
			expectedResponse: &game.GetGameResponse{},
			expectedError:    status.Error(codes.Internal, outerror.InternalMessage),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotResp, gotErr := GetGame(tc.err)
			assert.Equal(t, tc.expectedResponse, gotResp)
			assert.Equal(t, tc.expectedError, gotErr)
		})
	}
}
