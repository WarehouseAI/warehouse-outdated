package config

import "os"

type MailConfig struct {
	Host     string
	Sender   string
	User     string
	Password string
}

func NewMailCfg() MailConfig {
	return MailConfig{
		Host:     os.Getenv("MAIL_HOST"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Sender:   os.Getenv("MAIL_SENDER"),
		User:     os.Getenv("MAIL_USER"),
	}
}
