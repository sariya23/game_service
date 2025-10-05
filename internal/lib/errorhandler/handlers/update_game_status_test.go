package handlers

import (
	"fmt"
	"testing"

	"github.com/sariya23/game_service/internal/outerror"
	"github.com/sariya23/proto_api_games/v5/gen/gamev2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateGameStatus_errorhandler(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name         string
		err          error
		expectedResp *gamev2.UpdateGameStatusResponse
		expectedErr  error
	}{
		{
			name:         "UnknownGameStatus",
			err:          fmt.Errorf("%s: %w", "qwe", outerror.ErrUnknownGameStatus),
			expectedResp: &gamev2.UpdateGameStatusResponse{},
			expectedErr:  status.Error(codes.InvalidArgument, outerror.UnknownGameStatusMessage),
		},
		{
			name:         "InvalidNewStatus",
			err:          fmt.Errorf("%s: %w", "qwe", outerror.ErrInvalidNewGameStatus),
			expectedResp: &gamev2.UpdateGameStatusResponse{},
			expectedErr:  status.Error(codes.InvalidArgument, outerror.InvalidNewGameStatusMessage),
		},
		{
			name:         "GameNotFound",
			err:          fmt.Errorf("%s: %w", "qwe", outerror.ErrGameNotFound),
			expectedResp: &gamev2.UpdateGameStatusResponse{},
			expectedErr:  status.Error(codes.NotFound, outerror.GameNotFoundMessage),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotResp, gotErr := UpdateGameStatus(tc.err)
			assert.Equal(t, tc.expectedResp, gotResp)
			assert.Equal(t, tc.expectedErr, gotErr)
		})
	}

}
