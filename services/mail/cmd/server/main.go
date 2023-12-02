package server

import (
	"warehouseai/mail/cmd/broker"
	"warehouseai/mail/server/handlers"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func NewMailHandler(
	dialer *gomail.Dialer,
	consumer *broker.MailConsumer,
	logger *logrus.Logger,
	sender string,
) *handlers.Handler {
	return &handlers.Handler{
		MailDialer: dialer,
		Logger:     logger,
		Consumer:   consumer,
		Sender:     sender,
	}
}
