package email

import (
	"testing"

	"github.com/sariya23/game_service/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	cfg := config.MustLoadByPath("../../../config/local.env")
	sender := NewDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.EmailUser, cfg.EmailPassword)
	err := sender.SendMessage(cfg.EmailUser, cfg.AdminEmail, "Alert", "Hello")
	require.NoError(t, err)
}
