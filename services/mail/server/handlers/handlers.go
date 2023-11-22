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
	Consumer   *broker.MailConsumer
	MailDialer *gomail.Dialer
	Sender     string
}

func (h *Handler) SendMailHandler() {
	messages, err := h.Consumer.Channel.Consume(
		h.Consumer.Queue.Name,
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

	stop := make(chan bool)

	go func() {
		for message := range messages {
			var emailEvent model.EmailReceivedEvent
			json.Unmarshal(message.Body, &emailEvent)

			if err := service.SendEmail(h.Sender, emailEvent.Data, h.MailDialer, h.Logger); err != nil {
				continue
			}
		}
	}()

	<-stop
}
