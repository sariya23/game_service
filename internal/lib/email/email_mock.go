package email

import (
	"log/slog"

	"github.com/sariya23/game_service/internal/config"
)

type EmailDialerMock struct {
	log *slog.Logger
	To  string
}

func NewDialerMock(dialerConfig *config.Email) *EmailDialer {
	return &EmailDialer{
		To: dialerConfig.AdminEmail,
	}
}

func (dialer *EmailDialerMock) SendMessage(subject string, body string) error {
	return nil
}
