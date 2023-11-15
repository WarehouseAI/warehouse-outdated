package handlers

import (
	"encoding/json"
	"warehouseai/mail/model"
	"warehouseai/mail/service"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Handler struct {
	Logger *logrus.Logger
	Dialer *gomail.Dialer
	Sender string
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		receivedEmailEvent := model.EmailReceivedEvent{}

		if err := json.Unmarshal(msg.Value, &receivedEmailEvent); err != nil {
			continue
		}

		if mailErr := service.SendEmail(h.Sender, receivedEmailEvent.Data, h.Dialer, h.Logger); mailErr != nil {
			continue
		}
	}

	return nil
}

func (*Handler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (*Handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
