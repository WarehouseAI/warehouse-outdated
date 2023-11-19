package mail

import (
	"crypto/tls"
	"warehouseai/mail/config"

	"gopkg.in/gomail.v2"
)

func NewMailDialer(config config.MailConfig) *gomail.Dialer {
	dialer := gomail.NewDialer(config.Host, 25, config.User, config.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dialer
}
