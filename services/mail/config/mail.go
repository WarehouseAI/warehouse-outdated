package config

import "os"

type MailCfg struct {
	Host     string
	Sender   string
	User     string
	Password string
}

func NewMailCfg() MailCfg {
	return MailCfg{
		Host:     os.Getenv("MAIL_HOST"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Sender:   os.Getenv("MAIL_SENDER"),
		User:     os.Getenv("MAIL_USER"),
	}
}
