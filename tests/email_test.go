package tests

import (
	"testing"

	"github.com/sariya23/game_service/internal/config"
	"github.com/sariya23/game_service/internal/lib/email"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	cfg := config.MustLoadByPath("../config/local.env")
	sender := email.NewDialer(cfg.Email)
	err := sender.SendMessage("Alert", "Hello")
	require.NoError(t, err)
}
