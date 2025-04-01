package email

import "gopkg.in/gomail.v2"

type EmailDialer struct {
	dialer gomail.Dialer
}

func NewDialer(smtpHost string, smtpPort int, user string, Password string) *EmailDialer {
	dialer := gomail.NewDialer(smtpHost, smtpPort, user, Password)
	return &EmailDialer{
		dialer: *dialer,
	}
}

func (dialer *EmailDialer) SendMessage(to string, subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", dialer.dialer.Username)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	if err := dialer.dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
