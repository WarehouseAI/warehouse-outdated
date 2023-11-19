package handlers

import (
	"encoding/json"
	"warehouseai/mail/cmd/broker"
	"warehouseai/mail/model"
	"warehouseai/mail/service"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Handler struct {
	Logger     *logrus.Logger
	MailDialer *gomail.Dialer
	Sender     string
}

func (h *Handler) SendMailHandler() {
	channel, queue, connection := broker.NewMailConsumer()

	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		channel.Close()
		connection.Close()
	}()

	go func() {
		for message := range messages {
			var email model.Email
			json.Unmarshal(message.Body, &email)

			if err := service.SendEmail(h.Sender, email, h.MailDialer, h.Logger); err != nil {
				continue
			}
		}
	}()
}
