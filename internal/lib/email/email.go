package email

import "gopkg.in/gomail.v2"

type EmailDialer struct {
	dialer gomail.Dialer
	To     string
}

func NewDialer(smtpHost string, smtpPort int, user string, Password string, to string) *EmailDialer {
	dialer := gomail.NewDialer(smtpHost, smtpPort, user, Password)
	return &EmailDialer{
		dialer: *dialer,
		To:     to,
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
