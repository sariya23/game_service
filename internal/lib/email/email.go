package email

import (
	"github.com/sariya23/game_service/internal/config"
	"gopkg.in/gomail.v2"
)

type EmailDialer struct {
	dialer gomail.Dialer
	To     string
}

func NewDialer(dialerConfig *config.Email) *EmailDialer {
	dialer := gomail.NewDialer(dialerConfig.SmtpHost, dialerConfig.SmtpPort, dialerConfig.EmailUser, dialerConfig.EmailPassword)
	return &EmailDialer{
		dialer: *dialer,
		To:     dialerConfig.AdminEmail,
	}
}

func (dialer *EmailDialer) SendMessage(subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", dialer.dialer.Username)
	message.SetHeader("To", dialer.To)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	if err := dialer.dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
