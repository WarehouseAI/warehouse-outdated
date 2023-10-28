package mailutils

import (
	"crypto/tls"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(email Email) error {
	m := gomail.NewMessage()

	m.SetHeader("From", email.From)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", email.Message)

	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), 25, email.From, os.Getenv("MAIL_PASSWORD"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
