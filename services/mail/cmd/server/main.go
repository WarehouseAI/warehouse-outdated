package server

import (
	"warehouseai/mail/server/handlers"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func NewMailHandler(dialer *gomail.Dialer, logger *logrus.Logger, sender string) *handlers.Handler {
	return &handlers.Handler{
		Dialer: dialer,
		Logger: logger,
		Sender: sender,
	}
}
