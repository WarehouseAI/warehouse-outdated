package mail

import (
	"encoding/json"
	"warehouseai/user/model"

	"github.com/IBM/sarama"
)

type MailProducer struct {
	Producer sarama.SyncProducer
	Topic    string
}

func (k MailProducer) SendEmail(email model.Email) error {
	message := model.EmailRequest{Data: email}

	messageStr, err := json.Marshal(message)

	if err != nil {
		return err
	}

	producerMsg := &sarama.ProducerMessage{Topic: k.Topic, Value: sarama.StringEncoder(messageStr)}

	if _, _, err := k.Producer.SendMessage(producerMsg); err != nil {
		return err
	}

	return nil
}
